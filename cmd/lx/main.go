package main

import (
	"context"
	"log"
	"os"

	"github.com/rasros/lx/lx"
)

func main() {
	app := lx.NewCommand()
	args := lx.NormalizeArgs(os.Args)

	if err := app.Run(context.Background(), args); err != nil {
		log.Fatal(err)
	}
}
