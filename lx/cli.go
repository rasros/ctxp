package lx

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	ucli "github.com/urfave/cli/v3"
)

// Version is the application version. It can be overridden at build time via ldflags.
var Version = "dev"

// NormalizeArgs rewrites "-n2" / "-t10" / "-h5" into ["-n","2"] / ["-t","10"] / ["-h","5"]
// so that urfave/cli/v3 parses them as int flags.
func NormalizeArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	out := make([]string, 0, len(args)+4)
	for _, a := range args {
		if len(a) > 2 && a[0] == '-' && (a[1] == 'n' || a[1] == 't' || a[1] == 'h') {
			digits := a[2:]
			isDigits := true
			for _, ch := range digits {
				if ch < '0' || ch > '9' {
					isDigits = false
					break
				}
			}
			if isDigits {
				out = append(out, a[:2], digits)
				continue
			}
		}
		out = append(out, a)
	}

	return out
}

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

// NewCommand builds the urfave/cli command for lx.
func NewCommand() *ucli.Command {
	var head int
	var tail int
	var nBoth int
	var prefix string
	var postfix string

	// Make --help the only help flag (freeing -h for --head).
	ucli.HelpFlag = &ucli.BoolFlag{
		Name:        "help",
		Usage:       "show help",
		HideDefault: true,
		Local:       true,
	}

	return &ucli.Command{
		Name:    "lx",
		Usage:   "print files with headers, delimiters, and optional head/tail slicing",
		Version: Version,

		Flags: []ucli.Flag{
			&ucli.IntFlag{
				Name:        "head",
				Aliases:     []string{"h"},
				Usage:       "print first N lines (0 = no limit)",
				Destination: &head,
			},

			&ucli.IntFlag{
				Name:        "tail",
				Aliases:     []string{"t"},
				Usage:       "print last N lines (0 = no limit)",
				Destination: &tail,
			},

			// -n applies both head and tail.
			&ucli.IntFlag{
				Name:        "n",
				Usage:       "print first and last N lines (0 = no limit)",
				Destination: &nBoth,
			},

			&ucli.StringFlag{
				Name:        "prefix-delimiter",
				Usage:       "string printed before file contents; placeholders: {filename}, {row_count}",
				Destination: &prefix,
			},
			&ucli.StringFlag{
				Name:        "postfix-delimiter",
				Usage:       "string printed after file contents",
				Destination: &postfix,
			},
		},

		Action: func(ctx context.Context, cmd *ucli.Command) error {
			// Start with filenames from CLI args.
			files := cmd.Args().Slice()

			// Add filenames from piped stdin (one per line), if any.
			stdinFiles, err := readFilenamesFromStdin()
			if err != nil {
				return fmt.Errorf("lx: read stdin: %w", err)
			}
			if len(stdinFiles) > 0 {
				files = append(files, stdinFiles...)
			}

			if len(files) == 0 {
				return fmt.Errorf("lx: provide one or more file paths via args or stdin")
			}

			// Derive effective head/tail with override rules:
			// -n sets both head and tail, unless overridden by explicit --head/-h or --tail/-t.
			effHead := head
			effTail := tail

			if cmd.IsSet("n") {
				if !cmd.IsSet("head") { // -h/--head present => override -n for head
					effHead = nBoth
				}
				if !cmd.IsSet("tail") { // -t/--tail present => override -n for tail
					effTail = nBoth
				}
			}

			r := Runner{
				Head:             effHead,
				Tail:             effTail,
				PrefixDelimiter:  prefix,
				PostfixDelimiter: postfix,
			}

			if err := r.Run(files, os.Stdout); err != nil {
				return fmt.Errorf("lx: %w", err)
			}
			return nil
		},
	}
}

