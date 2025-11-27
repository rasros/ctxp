package ctxp

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Defaults(t *testing.T) {
	dir := t.TempDir()

	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")

	os.WriteFile(f1, []byte("line1\nline2\n"), 0o644)
	os.WriteFile(f2, []byte("x\ny\nz\n"), 0o644)

	var buf bytes.Buffer
	r := Runner{}

	if err := r.Run([]string{f1, f2}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	want :=
		f1 + " (2 rows)\n" +
			"---\n" +
			"```\n" +
			"line1\nline2\n" +
			"```\n" +
			f2 + " (3 rows)\n" +
			"---\n" +
			"```\n" +
			"x\ny\nz\n" +
			"```\n"

	if got != want {
		t.Fatalf("mismatch\nGot:\n%q\nWant:\n%q", got, want)
	}
}

func TestRun_HeadAndTail(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")

	content := []byte("one\ntwo\nthree\nfour\nfive\nsix\n")
	os.WriteFile(f, content, 0o644)

	var buf bytes.Buffer
	r := Runner{Head: 2, Tail: 2}

	if err := r.Run([]string{f}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want :=
		f + " (4 rows)\n" +
			"---\n" +
			"```\n" +
			"one\ntwo\nfive\nsix\n" +
			"```\n"

	if buf.String() != want {
		t.Fatalf("mismatch\nGot:\n%q\nWant:\n%q", buf.String(), want)
	}
}

func TestRun_HeadPlusTailLargerThanFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")

	os.WriteFile(f, []byte("l1\nl2\nl3\nl4\nl5\n"), 0o644)

	var buf bytes.Buffer
	r := Runner{Head: 3, Tail: 3} // head+tail > rows => full file

	if err := r.Run([]string{f}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want :=
		f + " (5 rows)\n" +
			"---\n" +
			"```\n" +
			"l1\nl2\nl3\nl4\nl5\n" +
			"```\n"

	if buf.String() != want {
		t.Fatalf("mismatch\nGot:\n%q\nWant:\n%q", buf.String(), want)
	}
}

