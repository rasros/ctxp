# lx

`lx` is a small CLI tool for printing files with clean, LLM-friendly delimiters and optional head/tail slicing. It works well with clipboard tools like `wl-copy` and `xclip` to streamline prompt preparation.

## Features

- Prints multiple files with markdown-style headers and fenced blocks.
- Optional slices for inspecting output
- Reads file paths from args or stdin for flexible usage.
- Customize prefix/postfix delimiters in shell aliases with placeholders.

## Installation

```bash
go install github.com/rasros/lx/cmd/lx@latest
```

## Usage

```bash
> lx cmd/lx/main.go
lx cmd/lx/main.go
cmd/lx/main.go (18 rows)
---
` ```go
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
\```
```

```bash
lx **/*.py~test_*.py
```

From stdin:
```bash
rg "Init" -l | lx -n 5
```

Custom delimiters:
```bash
lx --prefix-delimiter="### <filename>\n" file.go
```

## LLM Workflows
TODO
lx main.go | wl-copy
rg -tpy -l "def handler\(" | lx | wl-copy


## Placeholders
`{filename}`
`{row_count}` 
`{byte_size}`
`{last_modified}`
