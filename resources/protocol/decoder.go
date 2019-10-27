package protocol

import (
	"encoding/binary"
	"errors"
)

var (
	errInsufficientData = errors.New("insufficient data to decode packet")
	errNullString       = errors.New("null string")
	errVarIntOverflow   = errors.New("var int overflow")
)

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

func (d *Decoder) ArrayLength() (int, error) {
	n, err := d.Int32()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func (d *Decoder) Bool() (bool, error) {
	if err := d.remains(1); err != nil {
		return false, err
	}
	b := d.raw[d.offset] == 1
	d.offset += 1
	return b, nil
}

func (d *Decoder) Bytes() ([]byte, error) {
	n, err := d.Int32()
	if err != nil {
		return nil, err
	}
	length := int(n)

	if err := d.remains(length); err != nil {
		return nil, err
	}
	v := d.raw[d.offset : d.offset+length]
	d.offset += length
	return v, nil
}

func (d *Decoder) Int8() (int8, error) {
	if err := d.remains(1); err != nil {
		return 0, err
	}
	v := int8(d.raw[d.offset])
	d.offset += 1
	return v, nil
}

func (d *Decoder) Int16() (int16, error) {
	if err := d.remains(2); err != nil {
		return 0, err
	}
	v := int16(binary.BigEndian.Uint16(d.raw[d.offset:]))
	d.offset += 2
	return v, nil
}

func (d *Decoder) Int32() (int32, error) {
	if err := d.remains(4); err != nil {
		return 0, err
	}
	v := int32(binary.BigEndian.Uint32(d.raw[d.offset:]))
	d.offset += 4
	return v, nil
}

func (d *Decoder) Int32Array() ([]int32, error) {
	n, err := d.ArrayLength()
	if err != nil {
		return nil, err
	}

	if n == -1 {
		return nil, nil
	}

	items := make([]int32, n)
	for i := 0; i < n; i++ {
		item, err := d.Int32()
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	return items, nil
}

func (d *Decoder) Int64() (int64, error) {
	if err := d.remains(8); err != nil {
		return 0, err
	}
	v := int64(binary.BigEndian.Uint64(d.raw[d.offset:]))
	d.offset += 8
	return v, nil
}

func (d *Decoder) Int64Array() ([]int64, error) {
	n, err := d.ArrayLength()
	if err != nil {
		return nil, err
	}

	if n == -1 {
		return nil, nil
	}

	items := make([]int64, n)
	for i := 0; i < n; i++ {
		item, err := d.Int64()
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	return items, nil
}

func (d *Decoder) NullableString() (*string, error) {
	s, err := d.String()
	if err != nil {
		if err == errNullString {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

func (d *Decoder) String() (string, error) {
	n, err := d.Int16()
	if err != nil {
		return "", err
	}

	if n == -1 {
		return "", errNullString
	}

	length := int(n)
	if err := d.remains(length); err != nil {
		return "", err
	}

	s := string(d.raw[d.offset : d.offset+length])
	d.offset += length
	return s, nil
}

func (d *Decoder) StringArray() ([]string, error) {
	n, err := d.ArrayLength()
	if err != nil {
		return nil, err
	}

	if n == -1 {
		return nil, nil
	}

	items := make([]string, n)
	for i := 0; i < n; i++ {
		item, err := d.String()
		if err != nil {
			return nil, err
		}
		items[i] = item
	}

	return items, nil
}

func (d *Decoder) VarInt(v *int64) error {
	tmp, n := binary.Varint(d.raw[d.offset:])
	switch n {
	case 0:
		d.offset = len(d.raw) // no further requests can be made
		return errInsufficientData

	case -1:
		d.offset = len(d.raw) // no further requests can be made
		return errVarIntOverflow

	default:
		d.offset += n
		*v = tmp
		return nil
	}
}
