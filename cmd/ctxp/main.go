package main

import (
	"log"
	"os"

	"github.com/rasros/ctxp/internal/ctxp"
)

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		log.Fatal("ctxp: provide one or more file paths")
	}

	r := ctxp.Runner{
		Delimiter: "---",
	}

	if err := r.Run(files, os.Stdout); err != nil {
		log.Fatalf("ctxp: %v", err)
	}
}

