package sizeof

import "testing"

func TestBytes(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want int32
	}{
		{
			name: "nil",
			data: nil,
			want: ArrayLength,
		},
		{
			name: "some",
			data: []byte("hello world"),
			want: ArrayLength + int32(len("hello world")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes(tt.data); got != tt.want {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray(t *testing.T) {
	tests := []struct {
		name string
		ss   []string
		want int32
	}{
		{
			name: "nil",
			ss:   nil,
			want: 4,
		},
		{
			name: "simple",
			ss:   []string{"hello", "world"},
			want: 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringArray(tt.ss); got != tt.want {
				t.Errorf("StringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32Array(t *testing.T) {
	tests := []struct {
		name string
		ii   []int32
		want int32
	}{
		{
			name: "nil",
			ii:   nil,
			want: 4,
		},
		{
			name: "some",
			ii:   []int32{1, 2, 3},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int32Array(tt.ii); got != tt.want {
				t.Errorf("Int32Array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Array(t *testing.T) {
	tests := []struct {
		name string
		ii   []int64
		want int32
	}{
		{
			name: "nil",
			ii:   nil,
			want: 4,
		},
		{
			name: "some",
			ii:   []int64{1, 2, 3},
			want: 28,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64Array(tt.ii); got != tt.want {
				t.Errorf("Int64Array() = %v, want %v", got, tt.want)
			}
		})
	}
}
