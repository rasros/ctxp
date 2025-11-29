# lx

`lx` is a small CLI tool for printing files with clean, LLM-friendly delimiters and optional head/tail slicing. It works well with clipboard tools like `wl-copy` and `xclip` to streamline prompt preparation.

## Features

* Prints multiple files with markdown-style headers and fenced code blocks.
* Automatically detects the language for fenced code blocks based on file extension.
* Optional head/tail slicing:

  * `--head`, `-h` prints the first N lines
  * `--tail`, `-t` prints the last N lines
  * `-n` prints both the first and last N lines with an ellipsis line between them
* Reads file paths from CLI args or from stdin (one per line).
* Delimiters are customizable using placeholders such as `{filename}`, `{row_count}`, `{language}`, etc.

## Installation

```bash
go install github.com/rasros/lx/cmd/lx@latest
```

## Usage

Example:

~~~bash
> lx cmd/lx/main.go
lx cmd/lx/main.go
cmd/lx/main.go (18 rows)
---
```go
package main
...
~~~

```

Glob example:

```

lx **/*.py~test_*.py

```

From stdin:

```

rg "Init" -l | lx -n 5

```

Custom delimiters:

```

lx --prefix-delimiter="### {filename}\n```{language}\n" file.go

```

## LLM Workflows

```

lx main.go | wl-copy
rg -tpy -l "def handler(" | lx | wl-copy

```

## Placeholders

```

{filename}
{row_count}
{byte_size}
{last_modified}
{language}

```


