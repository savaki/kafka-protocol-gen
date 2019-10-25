package protocol

import (
	"encoding/binary"
	"io"
)

type flusher interface {
	Flush() error
}

type Encoder struct {
	buf    [8]byte
	target io.Writer
	err    error
}

func (e *Encoder) Flush() error {
	if e.err != nil {
		return e.err
	}

	if f, ok := e.target.(flusher); ok {
		return f.Flush()
	}

	return nil
}

func (e *Encoder) PutArray(n int, fn func(int)) {
	if e.err != nil {
		return
	}

	e.PutInt32(int32(n))
	for i := 0; i < n; i++ {
		fn(i)
	}
}

func (e *Encoder) PutBool(b bool) {
	var v int8
	if b {
		v = 1
	}
	e.PutInt8(v)
}

func (e *Encoder) PutBytes(data []byte) {
	if e.err != nil {
		return
	}

	e.PutInt32(int32(len(data)))
	_, e.err = e.target.Write(data)
}

func (e *Encoder) PutInt8(i int8) {
	if e.err != nil {
		return
	}

	e.buf[0] = byte(i)
	_, e.err = e.target.Write(e.buf[:1])
}

func (e *Encoder) PutInt16(i int16) {
	if e.err != nil {
		return
	}

	binary.BigEndian.PutUint16(e.buf[:2], uint16(i))
	_, e.err = e.target.Write(e.buf[:2])
}

func (e *Encoder) PutInt32(i int32) {
	if e.err != nil {
		return
	}

	binary.BigEndian.PutUint32(e.buf[:4], uint32(i))
	_, e.err = e.target.Write(e.buf[:4])
}

func (e *Encoder) PutInt32Array(ii []int32) {
	if e.err != nil {
		return
	}

	if ii == nil {
		e.PutInt32(-1)
		return
	}

	e.PutArray(len(ii), func(i int) { e.PutInt32(ii[i]) })
}

func (e *Encoder) PutInt64(i int64) {
	if e.err != nil {
		return
	}

	binary.BigEndian.PutUint64(e.buf[:8], uint64(i))
	_, e.err = e.target.Write(e.buf[:8])
}

func (e *Encoder) PutInt64Array(ii []int64) {
	if e.err != nil {
		return
	}

	if ii == nil {
		e.PutInt32(-1)
		return
	}

	e.PutArray(len(ii), func(i int) { e.PutInt64(ii[i]) })
}

func (e *Encoder) PutNullableString(s *string) {

}

func (e *Encoder) PutString(s string) {
	if e.err != nil {
		return
	}

	e.PutInt16(int16(len(s)))
	_, e.err = io.WriteString(e.target, s)
}

func (e *Encoder) PutStringArray(ss []string) {
	if e.err != nil {
		return
	}

	if ss == nil {
		e.PutInt32(-1)
		return
	}

	e.PutArray(len(ss), func(i int) { e.PutString(ss[i]) })
}

func (e *Encoder) PutVarInt64(i int64) {}

func (e *Encoder) PutVarString(s string) {}
