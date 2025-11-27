package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rasros/ctxp/internal/ctxp"
	"github.com/urfave/cli/v3"
)

// normalizeArgs rewrites "-n2" / "-t10" into ["-n","2"] / ["-t","10"] so that
// urfave/cli/v3 parses them as int flags.
func normalizeArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	out := make([]string, 0, len(args)+4)
	for _, a := range args {
		if len(a) > 2 && a[0] == '-' && (a[1] == 'n' || a[1] == 't') {
			// Check that the rest is all digits.
			digits := a[2:]
			isDigits := true
			for _, ch := range digits {
				if ch < '0' || ch > '9' {
					isDigits = false
					break
				}
			}
			if isDigits {
				// Split "-n2" → "-n", "2"
				out = append(out, a[:2], digits)
				continue
			}
		}
		out = append(out, a)
	}

	return out
}

func main() {
	var head int
	var tail int
	var prefix string
	var postfix string

	app := &cli.Command{
		Name:  "ctxp",
		Usage: "print files with headers, delimiters, and optional head/tail slicing",

		Flags: []cli.Flag{
			// long + short (-n / -n2)
			&cli.IntFlag{
				Name:        "head",
				Usage:       "print first N lines (0 = no limit)",
				Destination: &head,
			},
			&cli.IntFlag{
				Name:        "n",
				Usage:       "alias for --head",
				Destination: &head,
				Hidden:      true,
			},

			// long + short (-t / -t2)
			&cli.IntFlag{
				Name:        "tail",
				Usage:       "print last N lines (0 = no limit)",
				Destination: &tail,
			},
			&cli.IntFlag{
				Name:        "t",
				Usage:       "alias for --tail",
				Destination: &tail,
				Hidden:      true,
			},

			&cli.StringFlag{
				Name:        "prefix-delimiter",
				Usage:       "string printed before file contents; placeholders: <filename>, <row_count>",
				Destination: &prefix,
			},
			&cli.StringFlag{
				Name:        "postfix-delimiter",
				Usage:       "string printed after file contents",
				Destination: &postfix,
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				return fmt.Errorf("ctxp: provide one or more file paths")
			}

			r := ctxp.Runner{
				Head:             head,
				Tail:             tail,
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

	// Normalize args so "-n2" / "-t3" work.
	args := normalizeArgs(os.Args)

	if err := app.Run(context.Background(), args); err != nil {
		log.Fatal(err)
	}
}

