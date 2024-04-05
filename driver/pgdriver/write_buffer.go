package pgdriver

import (
	"encoding/binary"
	"sync"
)

var wbPool = sync.Pool{
	New: func() interface{} {
		return newWriteBuffer()
	},
}

func getWriteBuffer() *writeBuffer {
	wb := wbPool.Get().(*writeBuffer)
	return wb
}

func putWriteBuffer(wb *writeBuffer) {
	wb.Reset()
	wbPool.Put(wb)
}

type writeBuffer struct {
	Bytes []byte

	msgStart   int
	paramStart int
}

func newWriteBuffer() *writeBuffer {
	return &writeBuffer{
		Bytes: make([]byte, 0, 1024),
	}
}

func (b *writeBuffer) Reset() {
	b.Bytes = b.Bytes[:0]
}

func (b *writeBuffer) StartMessage(c byte) {
	if c == 0 {
		b.msgStart = len(b.Bytes)
		b.Bytes = append(b.Bytes, 0, 0, 0, 0)
	} else {
		b.msgStart = len(b.Bytes) + 1
		b.Bytes = append(b.Bytes, c, 0, 0, 0, 0)
	}
}

func (b *writeBuffer) WriteInt32(num int32) {
	b.Bytes = append(b.Bytes, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(b.Bytes[len(b.Bytes)-4:], uint32(num))
}

func (b *writeBuffer) WriteString(s string) {
	b.Bytes = append(b.Bytes, s...)
	b.Bytes = append(b.Bytes, 0)
}

func (b *writeBuffer) FinishMessage() {
	binary.BigEndian.PutUint32(
		b.Bytes[b.msgStart:], uint32(len(b.Bytes)-b.msgStart),
	)
}
