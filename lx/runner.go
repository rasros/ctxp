package lx

import (
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

// Run prints file contents with optional slicing and delimiters.
func (r Runner) Run(files []string, out io.Writer) error {
	if r.PrefixDelimiter == "" {
		r.PrefixDelimiter = "<filename> (<row_count> rows)\n---\n```\n"
	}
	if r.PostfixDelimiter == "" {
		r.PostfixDelimiter = "```\n\n"
	}

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}

		totalRows := countLines(data)
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
