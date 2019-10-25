package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecoder_PutBool(t *testing.T) {
	testCases := map[string]struct {
		want bool
	}{
		"true": {
			want: true,
		},
		"false": {
			want: false,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got bool
			)

			encoder := &Encoder{target: buf}
			encoder.PutBool(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Bool(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutBytes(t *testing.T) {
	testCases := map[string]struct {
		want []byte
	}{
		"nil": {
			want: nil,
		},
		"none": {
			want: []byte{},
		},
		"some": {
			want: []byte("hello"),
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got []byte
			)

			encoder := &Encoder{target: buf}
			encoder.PutBytes(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Bytes(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if !bytes.Equal(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutInt8(t *testing.T) {
	testCases := map[string]struct {
		want int8
	}{
		"pos": {
			want: 123,
		},
		"neg": {
			want: -123,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got int8
			)

			encoder := &Encoder{target: buf}
			encoder.PutInt8(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Int8(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutInt16(t *testing.T) {
	testCases := map[string]struct {
		want int16
	}{
		"pos": {
			want: 123,
		},
		"neg": {
			want: -123,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got int16
			)

			encoder := &Encoder{target: buf}
			encoder.PutInt16(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Int16(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutInt32(t *testing.T) {
	testCases := map[string]struct {
		want int32
	}{
		"pos": {
			want: 123,
		},
		"neg": {
			want: -123,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got int32
			)

			encoder := &Encoder{target: buf}
			encoder.PutInt32(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Int32(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutInt32Array(t *testing.T) {
	testCases := map[string]struct {
		want []int32
	}{
		"pos": {
			want: []int32{1, 2, 3},
		},
		"neg": {
			want: []int32{-1, -2, -3},
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got []int32
			)

			encoder := &Encoder{target: buf}
			encoder.PutInt32Array(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Int32Array(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutInt64(t *testing.T) {
	testCases := map[string]struct {
		want int64
	}{
		"pos": {
			want: 123,
		},
		"neg": {
			want: -123,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got int64
			)

			encoder := &Encoder{target: buf}
			encoder.PutInt64(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Int64(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutInt64Array(t *testing.T) {
	testCases := map[string]struct {
		want []int64
	}{
		"pos": {
			want: []int64{1, 2, 3},
		},
		"neg": {
			want: []int64{-1, -2, -3},
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got []int64
			)

			encoder := &Encoder{target: buf}
			encoder.PutInt64Array(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.Int64Array(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutString(t *testing.T) {
	testCases := map[string]struct {
		want string
	}{
		"none": {
			want: "",
		},
		"some": {
			want: "abc",
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got string
			)

			encoder := &Encoder{target: buf}
			encoder.PutString(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.String(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutStringArray(t *testing.T) {
	testCases := map[string]struct {
		want []string
	}{
		"nil": {
			want: nil,
		},
		"empty": {
			want: []string{},
		},
		"some": {
			want: []string{"a", "b", "c"},
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got []string
			)

			encoder := &Encoder{target: buf}
			encoder.PutStringArray(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.StringArray(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_remain(t *testing.T) {
	var (
		decoder = &Decoder{}
		err     error
	)

	var b bool
	err = decoder.Bool(&b)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var data []byte
	err = decoder.Bytes(&data)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var i8 int8
	err = decoder.Int8(&i8)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var i16 int16
	err = decoder.Int16(&i16)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var i32 int32
	err = decoder.Int32(&i32)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var ii32 []int32
	err = decoder.Int32Array(&ii32)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var i64 int64
	err = decoder.Int64(&i64)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var ii64 []int64
	err = decoder.Int64Array(&ii64)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var s string
	err = decoder.String(&s)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	var ss []string
	err = decoder.StringArray(&ss)
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}
}
