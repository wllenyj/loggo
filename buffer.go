package loggo

import (
	"bytes"
	"io"
	"sync"
)

type Buffer interface {
	Bytes() []byte
	String() string
	Len() int
}

var (
	buffer_pool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

type ByteWarp []byte

func (b ByteWarp) Bytes() []byte {
	return b
}
func (b ByteWarp) String() string {
	return string(b)
}
func (b ByteWarp) Len() int {
	return len(b)
}

type StringWarp string
func (s StringWarp) Bytes() []byte {
	return []byte(s)
}
func (s StringWarp) String() string {
	return string(s)
}
func (s StringWarp) Len() int {
	return len(s)
}

type BufferWriter interface {
	WriteBuffer(Buffer) (int, error)
}

type BufferWriterWarp struct {
	w io.Writer
}

func (ww *BufferWriterWarp) WriteBuffer(b Buffer) (int, error) {
	n, err := ww.w.Write(b.Bytes())
	return n, err
}
