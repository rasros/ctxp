# lx

`lx` is a small CLI tool for printing files with clean, LLM-friendly delimiters and optional head/tail slicing. It works well with clipboard tools like `wl-copy` and `xclip` to streamline prompt preparation.

## Features

- Prints multiple files with markdown-style headers and fenced blocks.
- Optional slices for inspecting output:
  - `--head, -h` for the first N lines
  - `--tail, -t` for the last N lines
  - `-n` for both, inserting an ellipsis line
- Reads file paths from args or stdin for flexible usage.
- Customize prefix/postfix delimiters in shel aliases with `{filename}` and `{row_count}` placeholders.

## Installation

```bash
go install github.com/rasros/lx/cmd/lx@latest
```

## Usage

```bash
lx file.go
lx -h 20 file.go
lx -t 50 file.go
lx -n 10 file.go
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



