package protocol

import (
	"bytes"
	"io"
	"testing"
)

func TestEncoder_err(t *testing.T) {
	var (
		buf     = bytes.NewBuffer(nil)
		wantErr = io.EOF
		encoder = Encoder{
			target: buf,
			err:    wantErr,
		}
	)

	var b bool
	encoder.PutBool(b)

	var data []byte
	encoder.PutBytes(data)

	var i8 int8
	encoder.PutInt8(i8)

	var i16 int16
	encoder.PutInt16(i16)

	var i32 int32
	encoder.PutInt32(i32)

	var ii32 []int32
	encoder.PutInt32Array(ii32)

	var i64 int64
	encoder.PutInt64(i64)

	var ii64 []int64
	encoder.PutInt64Array(ii64)

	var sp *string
	encoder.PutNullableString(sp)

	var s string
	encoder.PutString(s)

	var ss []string
	encoder.PutStringArray(ss)

	encoder.PutVarInt(123)

	if got, want := buf.Len(), 0; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}

	if got := encoder.Flush(); got != wantErr {
		t.Fatalf("got %v; want %v", got, wantErr)
	}
}
