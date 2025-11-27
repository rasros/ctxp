package ctxp

import (
	"bufio"
	"bytes"
	"strconv"
)

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

// splitLines splits data into logical lines, trimming the last empty chunk
// when the file ends with a newline.
func splitLines(data []byte) [][]byte {
	if len(data) == 0 {
		return nil
	}
	lines := bytes.SplitAfter(data, []byte("\n"))
	if len(lines) == 0 {
		return lines
	}
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// sliceLines returns data restricted by head/tail settings.
// Adds an explicit "... (N rows skipped)\n" line when both are used
// and the slice omits middle rows.
func sliceLines(data []byte, head, tail int) []byte {
	if head <= 0 && tail <= 0 {
		return data
	}

	lines := splitLines(data)
	total := len(lines)
	if total == 0 {
		return data
	}

	// If head or tail alone covers file, or together cover it, return full.
	if head >= total || tail >= total || (head > 0 && tail > 0 && head+tail >= total) {
		return data
	}

	var out [][]byte

	switch {
	case head > 0 && tail > 0:
		// Both specified; include an explanatory ellipsis line.
		skipped := total - head - tail

		out = append(out, lines[:head]...)
		out = append(out, []byte("... ("+strconv.Itoa(skipped)+" rows skipped)\n"))
		out = append(out, lines[total-tail:]...)

	case head > 0:
		out = lines[:head]

	case tail > 0:
		out = lines[total-tail:]
	}

	return bytes.Join(out, nil)
}

