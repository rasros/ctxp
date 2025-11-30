package lx

import (
	"bufio"
	"os"
	"strings"
)

// readFilenamesFromStdin reads filenames (one per line) from stdin when stdin
// is a pipe or redirection. If stdin is a TTY, it returns (nil, nil).
func readFilenamesFromStdin() ([]string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	// If stdin is a character device, there is no piped input.
	if info.Mode()&os.ModeCharDevice != 0 {
		return nil, nil
	}

	sc := bufio.NewScanner(os.Stdin)
	var paths []string
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		paths = append(paths, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return paths, nil
}
