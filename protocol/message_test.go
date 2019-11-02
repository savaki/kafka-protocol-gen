package protocol

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestValidVersions_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    ValidVersions
		data    string
		wantErr bool
	}{
		{
			name: "scalar",
			data: `"0"`,
			want: ValidVersions{},
		},
		{
			name: "range",
			data: `"1-3"`,
			want: ValidVersions{
				From: 1,
				To:   3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ValidVersions
			if err := json.Unmarshal([]byte(tt.data), &got); (err != nil) != tt.wantErr {
				t.Errorf("got %v; wantErr is %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v; want %#v\n", got, tt.want)
			}
		})
	}
}

func TestVersions_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    Versions
		data    string
		wantErr bool
	}{
		{
			name: "scalar",
			data: `"0"`,
			want: Versions{},
		},
		{
			name: "range",
			data: `"1-3"`,
			want: Versions{
				From: 1,
				To:   3,
			},
		},
		{
			name: "up o current",
			data: `"1+"`,
			want: Versions{
				From:        1,
				UpToCurrent: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Versions
			if err := json.Unmarshal([]byte(tt.data), &got); (err != nil) != tt.wantErr {
				t.Errorf("got %v; wantErr is %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v; want %#v\n", got, tt.want)
			}
		})
	}
}

func TestVersions_IsValid(t *testing.T) {
	testCases := map[string]struct {
		Versions Versions
		Version  int16
		Want     bool
	}{
		"scalar": {
			Versions: Versions{},
			Version:  0,
			Want:     true,
		},
		"0+": {
			Versions: Versions{
				UpToCurrent: true,
			},
			Version: 1,
			Want:    true,
		},
		"above range": {
			Versions: Versions{
				From: 1,
				To:   3,
			},
			Version: 4,
			Want:    false,
		},
		"below range": {
			Versions: Versions{
				From: 1,
				To:   3,
			},
			Version: 0,
			Want:    false,
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			got := tc.Versions.IsValid(tc.Version)
			if got != tc.Want {
				t.Fatalf("got %v; want %v", got, tc.Want)
			}
		})
	}
}
