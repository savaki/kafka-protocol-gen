package protocol

import (
	"encoding/binary"
	"errors"
)

var errInsufficientData = errors.New("insufficient data to decode packet")

func IsInsufficientDataError(err error) bool {
	return errors.Is(err, errInsufficientData)
}

type Decoder struct {
	raw    []byte
	offset int
}

func (d *Decoder) remains(n int) error {
	if remain := len(d.raw) - d.offset; remain < n {
		return errInsufficientData
	}
	return nil
}

func (d *Decoder) Array(fn func(n int, decoder *Decoder) error) error {
	if err := d.remains(4); err != nil {
		return err
	}

	var n int32
	if err := d.Int32(&n); err != nil {
		return err
	}

	if n < 0 {
		return nil
	}

	return fn(int(n), d)
}

func (d *Decoder) Bool(b *bool) error {
	if err := d.remains(1); err != nil {
		return err
	}
	*b = d.raw[d.offset] == 1
	d.offset += 1
	return nil
}

func (d *Decoder) Bytes(v *[]byte) error {
	var n int32
	if err := d.Int32(&n); err != nil {
		return err
	}
	length := int(n)

	if err := d.remains(length); err != nil {
		return err
	}
	*v = d.raw[d.offset : d.offset+length]
	d.offset += length
	return nil
}

func (d *Decoder) Int8(v *int8) error {
	if err := d.remains(1); err != nil {
		return err
	}
	*v = int8(d.raw[d.offset])
	d.offset += 1
	return nil
}

func (d *Decoder) Int16(v *int16) error {
	if err := d.remains(2); err != nil {
		return err
	}
	*v = int16(binary.BigEndian.Uint16(d.raw[d.offset:]))
	d.offset += 2
	return nil
}

func (d *Decoder) Int32(v *int32) error {
	if err := d.remains(4); err != nil {
		return err
	}
	*v = int32(binary.BigEndian.Uint32(d.raw[d.offset:]))
	d.offset += 4
	return nil
}

func (d *Decoder) Int32Array(vv *[]int32) error {
	fn := func(n int, decoder *Decoder) error {
		var items []int32
		for i := 0; i < n; i++ {
			var item int32
			if err := d.Int32(&item); err != nil {
				return err
			}
			items = append(items, item)
		}
		*vv = items
		return nil
	}
	return d.Array(fn)
}

func (d *Decoder) Int64(v *int64) error {
	if err := d.remains(8); err != nil {
		return err
	}
	*v = int64(binary.BigEndian.Uint64(d.raw[d.offset:]))
	d.offset += 8
	return nil
}

func (d *Decoder) Int64Array(vv *[]int64) error {
	fn := func(n int, decoder *Decoder) error {
		var items []int64
		for i := 0; i < n; i++ {
			var item int64
			if err := d.Int64(&item); err != nil {
				return err
			}
			items = append(items, item)
		}
		*vv = items
		return nil
	}
	return d.Array(fn)
}

func (d *Decoder) String(s *string) error {
	var n int16
	if err := d.Int16(&n); err != nil {
		return err
	}

	length := int(n)
	if err := d.remains(length); err != nil {
		return err
	}

	*s = string(d.raw[d.offset : d.offset+length])
	d.offset += length
	return nil
}

func (d *Decoder) StringArray(ss *[]string) error {
	fn := func(n int, decoder *Decoder) error {
		items := make([]string, 0, n)
		for i := 0; i < n; i++ {
			var item string
			if err := d.String(&item); err != nil {
				return err
			}
			items = append(items, item)
		}
		*ss = items
		return nil
	}
	return d.Array(fn)
}
