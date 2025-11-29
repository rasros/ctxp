package lx

import (
	"reflect"
	"testing"
)

func TestNormalizeArgs(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "n flag short",
			in:   []string{"lx", "-n2", "file"},
			want: []string{"lx", "-n", "2", "file"},
		},
		{
			name: "t flag short",
			in:   []string{"lx", "-t10"},
			want: []string{"lx", "-t", "10"},
		},
		{
			name: "h flag short",
			in:   []string{"lx", "-h5"},
			want: []string{"lx", "-h", "5"},
		},
		{
			name: "mixed",
			in:   []string{"lx", "-n2", "-t3", "file"},
			want: []string{"lx", "-n", "2", "-t", "3", "file"},
		},
		{
			name: "non digit suffix",
			in:   []string{"lx", "-nfoo"},
			want: []string{"lx", "-nfoo"},
		},
		{
			name: "looks like n but not digits",
			in:   []string{"lx", "-notflag"},
			want: []string{"lx", "-notflag"},
		},
		{
			name: "negative",
			in:   []string{"lx", "-h-5"},
			want: []string{"lx", "-h-5"},
		},
		{
			name: "other flags",
			in:   []string{"lx", "-x2"},
			want: []string{"lx", "-x2"},
		},
	}

	for _, tt := range tests {
		got := NormalizeArgs(tt.in)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: NormalizeArgs(%v) = %v, want %v",
				tt.name, tt.in, got, tt.want)
		}
	}
}

