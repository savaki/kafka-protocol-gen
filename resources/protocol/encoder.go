package protocol

import (
	"encoding/binary"
	"io"
)

type flusher interface {
	Flush() error
}

// Encoder provides protocol primitive encoders
type Encoder struct {
	buf    [16]byte
	target io.Writer
	err    error
}

// NewEncoder returns a new Encoder
func NewEncoder(target io.Writer) *Encoder {
	return &Encoder{
		target: target,
	}
}

// Flush encoding buffer
func (e *Encoder) Flush() error {
	if e.err != nil {
		return e.err
	}

	if f, ok := e.target.(flusher); ok {
		return f.Flush()
	}

	return nil
}

// PutArrayLength encodes the array length (int32)
func (e *Encoder) PutArrayLength(n int) {
	e.PutInt32(int32(n))
}

// PutBool encodes a bool
func (e *Encoder) PutBool(b bool) {
	if b {
		e.PutInt8(1)
	} else {
		e.PutInt8(0)
	}
}

// PutBool encodes a byte array
func (e *Encoder) PutBytes(data []byte) {
	e.PutInt32(int32(len(data)))
	if e.err == nil {
		_, e.err = e.target.Write(data)
	}
}

// PutInt8 encodes an int8
func (e *Encoder) PutInt8(i int8) {
	if e.err != nil {
		return
	}

	e.buf[0] = byte(i)
	_, e.err = e.target.Write(e.buf[:1])
}

// PutInt16 encodes an int16
func (e *Encoder) PutInt16(i int16) {
	if e.err != nil {
		return
	}

	binary.BigEndian.PutUint16(e.buf[:2], uint16(i))
	_, e.err = e.target.Write(e.buf[:2])
}

// PutInt32 encodes an int32
func (e *Encoder) PutInt32(i int32) {
	if e.err != nil {
		return
	}

	binary.BigEndian.PutUint32(e.buf[:4], uint32(i))
	_, e.err = e.target.Write(e.buf[:4])
}

// PutInt32Array encodes an []int32
func (e *Encoder) PutInt32Array(ii []int32) {
	if e.err != nil {
		return
	}

	if ii == nil {
		e.PutInt32(-1)
		return
	}

	length := len(ii)
	e.PutArrayLength(length)
	for _, i := range ii {
		e.PutInt32(i)
	}
}

// PutInt64 encodes an int64
func (e *Encoder) PutInt64(i int64) {
	if e.err != nil {
		return
	}

	binary.BigEndian.PutUint64(e.buf[:8], uint64(i))
	_, e.err = e.target.Write(e.buf[:8])
}

// PutInt64Array encodes an []int64
func (e *Encoder) PutInt64Array(ii []int64) {
	if e.err != nil {
		return
	}

	if ii == nil {
		e.PutInt32(-1)
		return
	}

	length := len(ii)
	e.PutArrayLength(length)
	for _, i := range ii {
		e.PutInt64(i)
	}
}

// PutNullableString encodes a *string
func (e *Encoder) PutNullableString(s *string) {
	if s == nil {
		e.PutInt16(-1)
		return
	}
	e.PutString(*s)
}

// PutString encodes a string
func (e *Encoder) PutString(s string) {
	e.PutInt16(int16(len(s)))
	if e.err == nil {
		_, e.err = io.WriteString(e.target, s)
	}
}

// PutStringArray encodes a []string
func (e *Encoder) PutStringArray(ss []string) {
	if e.err != nil {
		return
	}

	if ss == nil {
		e.PutInt32(-1)
		return
	}

	length := len(ss)
	e.PutArrayLength(length)
	for _, s := range ss {
		e.PutString(s)
	}
}

// PutVarInt encodes a var int
func (e *Encoder) PutVarInt(i int64) {
	if e.err != nil {
		return
	}

	length := binary.PutVarint(e.buf[:], i)
	_, e.err = e.target.Write(e.buf[0:length])
}
