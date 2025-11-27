package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rasros/ctxp/internal/ctxp"
	"github.com/urfave/cli/v3"
)

func main() {
	var head int
	var tail int
	var prefix string
	var postfix string

	app := &cli.Command{
		Name:  "ctxp",
		Usage: "print files with headers, delimiters, and optional head/tail slicing",

		Flags: []cli.Flag{
			// long + short (-n)
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

			// long + short (-t)
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

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

