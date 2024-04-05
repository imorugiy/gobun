package pgdriver

import (
	"bufio"
	"context"
	"crypto/tls"
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	bindMsg         = 'B'
	bindCompleteMsg = '2'

	authenticationOKMsg = 'R'
	errorResponseMsg    = 'E'
	parameterStatusMsg  = 'S'
	backendKeyDataMsg   = 'K'
	readyForQueryMsg    = 'Z'

	authenticationOK = 0
)

func writeBindExecute(ctx context.Context, cn *Conn, name string, args []driver.NamedValue) error {
	wb := getWriteBuffer()
	defer putWriteBuffer(wb)

	wb.StartMessage(bindMsg)

	return cn.write(ctx, wb)
}

type rowDescription struct {
	buf      []byte
	names    []string
	types    []int32
	numInput int16
}

type reader struct {
	*bufio.Reader
	buf []byte
}

func newReader(r io.Reader) *reader {
	return &reader{
		Reader: bufio.NewReader(r),
		buf:    make([]byte, 128),
	}
}

func (r *reader) ReadTemp(n int) ([]byte, error) {
	if n <= len(r.buf) {
		b := r.buf[:n]
		_, err := io.ReadFull(r.Reader, b)
		return b, err
	}

	b := make([]byte, n)
	_, err := io.ReadFull(r.Reader, b)
	return b, err
}

func (r *reader) Discard(n int) error {
	_, err := r.ReadTemp(n)
	return err
}

func readExtQueryData(ctx context.Context, cn *Conn, rowDesc *rowDescription) (*rows, error) {
	rd := cn.reader(ctx, -1)
	// var firstErr error
	for {
		c, msgLen, err := readMessageType(rd)
		if err != nil {
			return nil, err
		}

		switch c {
		case bindCompleteMsg:
			if err := rd.Discard(msgLen); err != nil {
				return nil, err
			}
			return newRows(cn, rowDesc, false), nil
		}
	}
}

func readMessageType(rd *reader) (byte, int, error) {
	c, err := rd.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	l, err := readInt32(rd)
	if err != nil {
		return 0, 0, err
	}
	return c, int(l) - 4, nil
}

func readInt32(rd *reader) (int32, error) {
	b, err := rd.ReadTemp(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

func startup(ctx context.Context, cn *Conn) error {
	if err := writeStartup(ctx, cn); err != nil {
		return err
	}

	rd := cn.reader(ctx, -1)

	for {
		c, msgLen, err := readMessageType(rd)
		fmt.Println(string(c), msgLen)
		if err != nil {
			return err
		}
		switch c {
		case backendKeyDataMsg:
			processID, err := readInt32(rd)
			if err != nil {
				return err
			}
			secretKey, err := readInt32(rd)
			if err != nil {
				return err
			}
			cn.processID = processID
			cn.secretKey = secretKey
		case authenticationOKMsg:
			if err := auth(ctx, cn, rd); err != nil {
				return err
			}
		case readyForQueryMsg:
			return rd.Discard(msgLen)
		case parameterStatusMsg:
			if err := rd.Discard(msgLen); err != nil {
				return err
			}
		case errorResponseMsg:
			e, err := readError(rd)
			if err != nil {
				return err
			}
			return e
		default:
			return fmt.Errorf("pgdriver: unexpected startup message: %q", c)
		}
	}
}

func writeStartup(ctx context.Context, cn *Conn) error {
	wb := getWriteBuffer()
	defer putWriteBuffer(wb)

	wb.StartMessage(0)
	wb.WriteInt32(196608)
	wb.WriteString("user")
	wb.WriteString(cn.cfg.User)
	wb.WriteString("database")
	wb.WriteString(cn.cfg.Database)
	if cn.cfg.AppName != "" {
		wb.WriteString("application_name")
		wb.WriteString(cn.cfg.AppName)
	}
	wb.WriteString("")
	wb.FinishMessage()

	return cn.write(ctx, wb)
}

func auth(ctx context.Context, cn *Conn, rd *reader) error {
	num, err := readInt32(rd)
	if err != nil {
		return err
	}

	switch num {
	case authenticationOK:
		return nil
	default:
		return fmt.Errorf("pgdriver: unknown authentication message: %q", num)
	}
}

func enableSSL(ctx context.Context, cn *Conn, tlsConf *tls.Config) error {
	if err := writeSSLMsg(ctx, cn); err != nil {
		return err
	}

	rd := cn.reader(ctx, -1)

	c, err := rd.ReadByte()
	if err != nil {
		return err
	}

	if c != 'S' {
		return errors.New("pgdriver: SSL is not enabled on the server")
	}

	tlsCN := tls.Client(cn.netConn, tlsConf)
	if err := tlsCN.HandshakeContext(ctx); err != nil {
		return fmt.Errorf("pgdriver: TLS handshake failed: %w", err)
	}
	cn.netConn = tlsCN
	rd.Reset(cn.netConn)

	return nil
}

func writeSSLMsg(ctx context.Context, cn *Conn) error {
	wb := getWriteBuffer()
	defer putWriteBuffer(wb)

	wb.StartMessage(0)
	wb.WriteInt32(80877103)
	wb.FinishMessage()

	return cn.write(ctx, wb)
}

// func readExtQuery(ctx context.Context, cn *Conn) (driver.Result, error) {
// }

func readString(rd *reader) (string, error) {
	b, err := rd.ReadSlice(0)
	if err != nil {
		return "", err
	}
	return string(b[:len(b)-1]), nil
}

func readError(rd *reader) (error, error) {
	m := make(map[byte]string)
	for {
		c, err := rd.ReadByte()
		if err != nil {
			return nil, err
		}
		if c == 0 {
			break
		}
		s, err := readString(rd)
		if err != nil {
			return nil, err
		}
		m[c] = s
	}
	switch err := (Error{m: m}); err.Field('V') {
	case "FATAL", "PANIC":
		return nil, err
	default:
		return err, nil
	}
}
