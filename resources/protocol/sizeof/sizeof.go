package sizeof

const (
	Bool        int32 = 1
	Int8        int32 = 1
	Int16       int32 = 2
	Int32       int32 = 4
	Int64       int32 = 8
	ArrayLength       = Int32
)

func Bytes(data []byte) int32 {
	return ArrayLength + int32(len(data)) // int32 length + length of bytes
}

func Int32Array(ii []int32) int32 {
	return ArrayLength + int32(len(ii))*Int32 // int32 length + length of array * int32 length
}

func Int64Array(ii []int64) int32 {
	return ArrayLength + int32(len(ii))*Int64 // int32 length + length of array * int64 length
}

func String(s string) int32 {
	return Int16 + int32(len(s))
}

func StringArray(ss []string) int32 {
	var sz int32
	sz += ArrayLength // int32 length of array
	for _, s := range ss {
		sz += String(s)
	}
	return sz
}
