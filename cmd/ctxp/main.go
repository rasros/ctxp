package main

import (
	"context"
	"log"
	"os"

	"github.com/rasros/ctxp/internal/ctxp"
)

func main() {
	app := ctxp.NewCommand()
	args := ctxp.NormalizeArgs(os.Args)

	if err := app.Run(context.Background(), args); err != nil {
		log.Fatal(err)
	}
}
