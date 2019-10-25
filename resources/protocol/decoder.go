package protocol

type Decoder struct {
}

func (d *Decoder) Array(fn func(decoder *Decoder, n int) error) error {
	return nil
}

func (d *Decoder) Bool(b *bool) error {
	return nil
}

func (d *Decoder) Bytes(v *[]byte) error {
	return nil
}

func (d *Decoder) Int8(v *int8) error {
	return nil
}

func (d *Decoder) Int16(v *int16) error {
	return nil
}

func (d *Decoder) Int32(v *int32) error {
	return nil
}

func (d *Decoder) Int64(v *int64) error {
	return nil
}

func (d *Decoder) String(s *string) error {
	return nil
}
