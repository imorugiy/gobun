package pgdriver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net"
	"time"
)

type Connector struct {
	cfg *Config
}

func NewConnector() *Connector {
	c := &Connector{cfg: newDefaultConfig()}
	return c
}

var _ driver.Connector = (*Connector)(nil)

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if err := c.cfg.verify(); err != nil {
		return nil, err
	}
	return newConn(ctx, c.cfg)
}

func (c *Connector) Driver() driver.Driver {
	fmt.Println("fix here")
	return &Driver{}
}

type Driver struct{}

var _ driver.Driver = (*Driver)(nil)

func (d *Driver) Open(name string) (driver.Conn, error) {
	fmt.Println("fix here")
	return &Conn{}, nil
}

type Conn struct {
	cfg *Config

	netConn net.Conn
	rd      *reader

	processID int32
	secretKey int32
}

func newConn(ctx context.Context, cfg *Config) (*Conn, error) {
	netConn, err := cfg.Dialer(ctx, cfg.Network, cfg.Addr)
	if err != nil {
		return nil, err
	}

	cn := &Conn{
		cfg:     cfg,
		netConn: netConn,
		rd:      newReader(netConn),
	}

	if cfg.TLSConfig != nil {
		if err := enableSSL(ctx, cn, cfg.TLSConfig); err != nil {
			return nil, err
		}
	}

	if err := startup(ctx, cn); err != nil {
		return nil, err
	}

	return cn, nil
}

func (cn *Conn) reader(ctx context.Context, timeout time.Duration) *reader {
	cn.setReadDeadline(ctx, timeout)
	return cn.rd
}

func (cn *Conn) write(ctx context.Context, wb *writeBuffer) error {
	cn.setWriteDeadline(ctx, -1)

	n, err := cn.netConn.Write(wb.Bytes)
	wb.Reset()

	if err != nil {
		if n == 0 {
			return driver.ErrBadConn
		}
		return err
	}
	return nil
}

var _ driver.Conn = (*Conn)(nil)

var _ driver.ExecerContext = (*Conn)(nil)

func (cn *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	// query, err := formatQuery(query, args)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

var _ driver.Pinger = (*Conn)(nil)

func (cn *Conn) Ping(ctx context.Context) error {
	_, err := cn.ExecContext(ctx, "SELECT 1", nil)
	return err
}

func (cn *Conn) setReadDeadline(ctx context.Context, timeout time.Duration) {
	if timeout == -1 {
		timeout = cn.cfg.ReadTimeout
	}
	_ = cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout))
}

func (cn *Conn) setWriteDeadline(ctx context.Context, timeout time.Duration) {
	if timeout == -1 {
		timeout = cn.cfg.WriteTimeout
	}
	_ = cn.netConn.SetWriteDeadline(cn.deadline(ctx, timeout))
}

func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	deadline, ok := ctx.Deadline()
	if !ok {
		if timeout == 0 {
			return time.Time{}
		}
		return time.Now().Add(timeout)
	}

	if timeout == 0 {
		return deadline
	}
	if tm := time.Now().Add(timeout); tm.Before(deadline) {
		return tm
	}
	return deadline
}

func (cn *Conn) Begin() (driver.Tx, error) {
	return nil, nil
}

func (cn *Conn) Close() error {
	return nil
}

func (cn *Conn) Prepare(query string) (driver.Stmt, error) {
	return newStmt(cn, "", nil), nil
}

type rows struct {
	cn       *Conn
	rowDesc  *rowDescription
	reusable bool
	closed   bool
}

var _ driver.Rows = (*rows)(nil)

func newRows(cn *Conn, rowDesc *rowDescription, reusable bool) *rows {
	return &rows{
		cn:       cn,
		rowDesc:  rowDesc,
		reusable: reusable,
	}
}

func (r *rows) Columns() []string {
	return []string{}
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	return nil
}

type stmt struct {
	cn      *Conn
	name    string
	rowDesc *rowDescription
}

var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)
)

func newStmt(cn *Conn, name string, rowDesc *rowDescription) *stmt {
	return &stmt{cn: cn, name: name, rowDesc: rowDesc}
}

func (stmt *stmt) Close() error {
	return nil
}

func (stmt *stmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("not implemented")
}

func (stmt *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	panic("not implemented")
}

func (stmt *stmt) NumInput() int {
	return 0
}

func (stmt *stmt) Query(args []driver.Value) (driver.Rows, error) {
	panic("not implemented")
}

func (stmt *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if err := writeBindExecute(ctx, stmt.cn, stmt.name, args); err != nil {
		return nil, err
	}
	// fmt.Println(stmt.cn.)
	return readExtQueryData(ctx, stmt.cn, stmt.rowDesc)
}
