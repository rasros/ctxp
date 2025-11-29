package lx

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Runner struct {
	Head             int
	Tail             int
	PrefixDelimiter  string
	PostfixDelimiter string
	LineNumbers      bool
}

// addLineNumbers prefixes each logical line with "N: ", where N is the
// original file line number. Ellipsis lines ("... (N rows skipped)") are left
// unnumbered.
func addLineNumbers(data []byte, totalRows, head, tail int) []byte {
	if len(data) == 0 {
		return data
	}

	lines := splitLines(data)
	if len(lines) == 0 {
		return data
	}

	// Helper to number a slice of lines starting at a given file line number.
	numberLines := func(lines [][]byte, startLine int) []byte {
		buf := make([]byte, 0, len(data)+len(lines)*8)
		lineNum := startLine
		for _, ln := range lines {
			num := strconv.Itoa(lineNum)
			buf = append(buf, []byte(num)...)
			buf = append(buf, ':', ' ')
			buf = append(buf, ln...)
			lineNum++
		}
		return buf
	}

	// Decide if we effectively printed the full file (no slicing or slicing
	// that covers all rows). In that case, number sequentially 1..totalRows.
	fullFile := (head <= 0 && tail <= 0) ||
		head >= totalRows ||
		tail >= totalRows ||
		(head > 0 && tail > 0 && head+tail >= totalRows)

	if fullFile {
		return numberLines(lines, 1)
	}

	// Now we know that slicing actually removed some rows.
	switch {
	case head > 0 && tail > 0 && head+tail < totalRows:
		// Mixed head+tail with an ellipsis line in the middle.
		// Layout of `lines` (from sliceLines):
		//   lines[0:head]           => original 1..head
		//   lines[head]             => ellipsis ("... (N rows skipped)")
		//   lines[head+1:]          => original (totalRows-tail+1)..totalRows
		buf := make([]byte, 0, len(data)+len(lines)*8)

		// First head lines: line numbers 1..head.
		if head > len(lines) {
			head = len(lines)
		}
		if head > 0 {
			buf = append(buf, numberLines(lines[:head], 1)...)
		}

		// Ellipsis line (if present) unnumbered.
		if head < len(lines) {
			buf = append(buf, lines[head]...)
		}

		// Tail lines: their original line numbers start at totalRows-tail+1.
		if head+1 < len(lines) {
			startOrig := totalRows - tail + 1
			// lines[head+1:] corresponds exactly to tail lines.
			bufTail := numberLines(lines[head+1:], startOrig)
			buf = append(buf, bufTail...)
		}

		return buf

	case head > 0 && (tail <= 0):
		// Only head: first head lines are printed, numbered from 1.
		return numberLines(lines, 1)

	case tail > 0 && (head <= 0):
		// Only tail: last tail lines are printed, numbered from totalRows-tail+1.
		startOrig := totalRows - tail + 1
		return numberLines(lines, startOrig)
	}

	// Fallback: sequential numbering (should not normally reach here).
	return numberLines(lines, 1)
}

// Run prints file contents with optional slicing and delimiters.
func (r Runner) Run(files []string, out io.Writer) error {
	if r.PrefixDelimiter == "" {
		r.PrefixDelimiter = "{filename} ({row_count} rows)\n---\n```{language}\n"
	}
	if r.PostfixDelimiter == "" {
		r.PostfixDelimiter = "```\n\n"
	}

	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat %q: %w", path, err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}

		totalRows := countLines(data)
		byteSize := info.Size()
		lastMod := info.ModTime().Format(time.RFC3339)
		lang := languageFromPath(path)

		view := sliceLines(data, r.Head, r.Tail)

		prefix := r.PrefixDelimiter
		prefix = strings.ReplaceAll(prefix, "{filename}", path)
		prefix = strings.ReplaceAll(prefix, "{row_count}", strconv.Itoa(totalRows))
		prefix = strings.ReplaceAll(prefix, "{byte_size}", strconv.FormatInt(byteSize, 10))
		prefix = strings.ReplaceAll(prefix, "{last_modified}", lastMod)
		prefix = strings.ReplaceAll(prefix, "{language}", lang)

		if _, err := out.Write([]byte(prefix)); err != nil {
			return fmt.Errorf("write prefix: %w", err)
		}

		var toWrite []byte
		if r.LineNumbers {
			toWrite = addLineNumbers(view, totalRows, r.Head, r.Tail)
		} else {
			toWrite = view
		}

		if _, err := out.Write(toWrite); err != nil {
			return fmt.Errorf("write data: %w", err)
		}
		if _, err := out.Write([]byte(r.PostfixDelimiter)); err != nil {
			return fmt.Errorf("write postfix: %w", err)
		}
	}

	return nil
}

