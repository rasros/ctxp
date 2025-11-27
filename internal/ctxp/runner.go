package ctxp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

type Runner struct {
	// Delimiter is the line printed between header and contents, e.g. "---".
	Delimiter string
}

// countLines counts the number of newline-terminated rows.
// Files without a trailing newline still count their last row.
func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	count := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		count++
	}
	return count
}

// Run prints, for each file, exactly:
//
//   <filename> (<rows> rows)
//   ---
//   <file contents>
//
func (r Runner) Run(files []string, out io.Writer) error {
	if r.Delimiter == "" {
		r.Delimiter = "---"
	}

	for _, p := range files {
		data, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("read %q: %w", p, err)
		}

		rows := countLines(data)

		// Header line
		if _, err := fmt.Fprintf(out, "%s (%d rows)\n", p, rows); err != nil {
			return err
		}

		// Delimiter line
		if _, err := fmt.Fprintf(out, "%s\n", r.Delimiter); err != nil {
			return err
		}

		// File contents (as-is)
		if _, err := out.Write(data); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}

	return nil
}

