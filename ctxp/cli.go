package ctxp

import (
	"context"
	"fmt"
	"os"

	ucli "github.com/urfave/cli/v3"
)

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

// NewCommand builds the urfave/cli command for ctxp.
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
		Name:  "ctxp",
		Usage: "print files with headers, delimiters, and optional head/tail slicing",

		// Small, focused help for `ctxp --help`.
		CustomRootCommandHelpTemplate: `USAGE:
  ctxp [flags] FILE...

FLAGS:
  --head, -h N              print first N lines (0 = no limit)
  --tail, -t N              print last N lines (0 = no limit)
  -n N                      print first and last N lines (0 = no limit)
  --prefix-delimiter VALUE  string printed before file contents; placeholders: <filename>, <row_count>
  --postfix-delimiter VALUE string printed after file contents
  --help                    show this help
`,

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
				Usage:       "string printed before file contents; placeholders: <filename>, <row_count>",
				Destination: &prefix,
			},
			&ucli.StringFlag{
				Name:        "postfix-delimiter",
				Usage:       "string printed after file contents",
				Destination: &postfix,
			},
		},

		Action: func(ctx context.Context, cmd *ucli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("ctxp: provide one or more file paths")
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

			files := cmd.Args().Slice()
			if err := r.Run(files, os.Stdout); err != nil {
				return fmt.Errorf("ctxp: %w", err)
			}
			return nil
		},
	}
}

