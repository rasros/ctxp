package ctxp

import (
	"os"
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
			in:   []string{"ctxp", "-n2", "file"},
			want: []string{"ctxp", "-n", "2", "file"},
		},
		{
			name: "t flag short",
			in:   []string{"ctxp", "-t10"},
			want: []string{"ctxp", "-t", "10"},
		},
		{
			name: "h flag short",
			in:   []string{"ctxp", "-h5"},
			want: []string{"ctxp", "-h", "5"},
		},
		{
			name: "mixed",
			in:   []string{"ctxp", "-n2", "-t3", "file"},
			want: []string{"ctxp", "-n", "2", "-t", "3", "file"},
		},
		{
			name: "non digit suffix",
			in:   []string{"ctxp", "-nfoo"},
			want: []string{"ctxp", "-nfoo"},
		},
		{
			name: "looks like n but not digits",
			in:   []string{"ctxp", "-notflag"},
			want: []string{"ctxp", "-notflag"},
		},
		{
			name: "negative",
			in:   []string{"ctxp", "-h-5"},
			want: []string{"ctxp", "-h-5"},
		},
		{
			name: "other flags",
			in:   []string{"ctxp", "-x2"},
			want: []string{"ctxp", "-x2"},
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

func TestReadFilenamesFromStdin_Piped(t *testing.T) {
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Stdin = origStdin
		r.Close()
		w.Close()
	}()

	os.Stdin = r

	input := "file1.txt\n\nfile2.txt \n  \nfile3\n"
	if _, err := w.Write([]byte(input)); err != nil {
		t.Fatal(err)
	}
	w.Close()

	got, err := readFilenamesFromStdin()
	if err != nil {
		t.Fatalf("readFilenamesFromStdin error: %v", err)
	}

	want := []string{"file1.txt", "file2.txt", "file3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readFilenamesFromStdin = %v, want %v", got, want)
	}
}

