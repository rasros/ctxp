package ctxp

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_WithHeadersAndDelimiter(t *testing.T) {
	dir := t.TempDir()

	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")

	if err := os.WriteFile(f1, []byte("line1\nline2\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte("x\ny\nz\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	r := Runner{Delimiter: "---"}

	err := r.Run([]string{f1, f2}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want :=
		f1 + " (2 rows)\n" +
			"---\n" +
			"line1\nline2\n" +
			f2 + " (3 rows)\n" +
			"---\n" +
			"x\ny\nz\n"

	if got != want {
		t.Fatalf("output mismatch\n got:\n%q\nwant:\n%q", got, want)
	}
}

