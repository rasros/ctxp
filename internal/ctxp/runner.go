package ctxp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Runner struct {
	Head             int
	Tail             int
	PrefixDelimiter  string
	PostfixDelimiter string
}

// countLines counts the number of newline-terminated rows.
// Files without a trailing newline still count their last row.
func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	sc := bufio.NewScanner(bytes.NewReader(data))
	n := 0
	for sc.Scan() {
		n++
	}
	return n
}

// splitLines splits into logical lines, trimming a trailing empty chunk
// when the file ends with a newline, so the count matches countLines.
func splitLines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	lines := bytes.SplitAfter(data, []byte("\n"))
	if len(lines) == 0 {
		return lines
	}
	// If the last chunk is empty (which happens when the file ends with "\n"),
	// drop it so that len(lines) == countLines(data).
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// sliceLines returns a view of data restricted by head/tail settings.
// head or tail of 0 means "no restriction from that side".
// If head+tail >= total number of lines, the full file is returned.
func sliceLines(data []byte, head, tail int) []byte {
	if head <= 0 && tail <= 0 {
		return data
	}

	lines := splitLines(data)
	total := len(lines)
	if total == 0 {
		return data
	}

	// If head or tail alone reaches/exceeds total, or the sum does, keep full file.
	if head >= total || tail >= total || (head > 0 && tail > 0 && head+tail >= total) {
		return data
	}

	var out [][]byte

	switch {
	case head > 0 && tail > 0:
		out = append(out, lines[:head]...)
		out = append(out, lines[total-tail:]...)
	case head > 0:
		out = lines[:head]
	case tail > 0:
		out = lines[total-tail:]
	}

	return bytes.Join(out, nil)
}

// Run prints, for each file:
//
//   <filename> (<row_count> rows)
//   ---
//   ```
//   <possibly-sliced file contents>
//   ```
//
func (r Runner) Run(files []string, out io.Writer) error {
	if r.PrefixDelimiter == "" {
		r.PrefixDelimiter = "<filename> (<row_count> rows)\n---\n```\n"
	}
	if r.PostfixDelimiter == "" {
		r.PostfixDelimiter = "```\n"
	}

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}

		// Real row count from the full file (what you asked for).
		totalRows := countLines(data)

		// Apply head/tail slicing for the printed view.
		view := sliceLines(data, r.Head, r.Tail)

		prefix := r.PrefixDelimiter
		prefix = strings.ReplaceAll(prefix, "<filename>", path)
		prefix = strings.ReplaceAll(prefix, "<row_count>", strconv.Itoa(totalRows))

		if _, err := out.Write([]byte(prefix)); err != nil {
			return fmt.Errorf("write prefix: %w", err)
		}

		if _, err := out.Write(view); err != nil {
			return fmt.Errorf("write data: %w", err)
		}

		if _, err := out.Write([]byte(r.PostfixDelimiter)); err != nil {
			return fmt.Errorf("write postfix: %w", err)
		}
	}

	return nil
}

