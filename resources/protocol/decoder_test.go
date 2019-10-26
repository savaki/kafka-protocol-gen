package protocol

import (
	"bytes"
	"math"
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
			got, err = decoder.Bool()
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
			got, err = decoder.Bytes()
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
			got, err = decoder.Int8()
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
			got, err = decoder.Int16()
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
			got, err = decoder.Int32()
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
		"nil": {
			want: nil,
		},
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
			got, err = decoder.Int32Array()
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
			got, err = decoder.Int64()
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
		"nil": {
			want: nil,
		},
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
			got, err = decoder.Int64Array()
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
			got, err = decoder.String()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutNullableString(t *testing.T) {
	var (
		blank = ""
		some  = "some"
	)

	testCases := map[string]struct {
		want *string
	}{
		"blank": {
			want: &blank,
		},
		"nil": {
			want: nil,
		},
		"some": {
			want: &some,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got *string
			)

			encoder := &Encoder{target: buf}
			encoder.PutNullableString(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			got, err = decoder.NullableString()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got == nil && tc.want != nil {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
			if got != nil && tc.want == nil {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
			if got != nil && tc.want != nil && *got != *tc.want {
				t.Fatalf("got %v; want %v", *got, *tc.want)
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
			got, err = decoder.StringArray()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestDecoder_PutVarInt(t *testing.T) {
	testCases := map[string]struct {
		want int64
	}{
		"0": {
			want: 0,
		},
		"1": {
			want: 1,
		},
		"-1": {
			want: 1,
		},
		"-int64": {
			want: math.MinInt64,
		},
		"int64": {
			want: math.MaxInt64,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			var (
				buf = bytes.NewBuffer(nil)
				got int64
			)

			encoder := &Encoder{target: buf}
			encoder.PutVarInt(tc.want)
			err := encoder.Flush()
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			decoder := &Decoder{raw: buf.Bytes()}
			err = decoder.VarInt(&got)
			if err != nil {
				t.Fatalf("got %v; want nil", err)
			}

			if got != tc.want {
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

	_, err = decoder.Bool()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Bytes()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Int8()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Int16()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Int32()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Int32Array()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Int64()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.Int64Array()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.NullableString()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.String()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}

	_, err = decoder.StringArray()
	if !IsInsufficientDataError(err) {
		t.Fatalf("got %v; want %v", err, errInsufficientData)
	}
}
