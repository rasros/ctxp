# lx

`lx` is a small CLI tool for printing files with clean, LLM-friendly delimiters 
and optional head/tail slicing. It works well with clipboard tools like 
`wl-copy` and `xclip` to streamline prompt preparation.

## Features

* Prints multiple files with markdown-style headers and fenced code blocks.
* Automatically detects the language for fenced code blocks based on file 
extension.
* Optional head/tail slicing:

  * `--head`, `-h` prints the first N lines
  * `--tail`, `-t` prints the last N lines
  * `-n` prints both the first and last N lines with an ellipsis line between them
* Reads file paths from CLI args or from stdin (one per line).
* Delimiters are customizable using placeholders such as `{filename}`, 
`{row_count}`, `{language}`, etc.

## Installation

```bash
go install github.com/rasros/lx/cmd/lx@latest
```

## Basic usage

Enter a filename to format it.

~~~bash
> lx cmd/lx/main.go
lx cmd/lx/main.go
cmd/lx/main.go (18 rows)
---
```go
package main
... (rest of the output omitted in README)
```
~~~

Instead of manually copy paste it use a copy util like `wl-copy` for wayland, 
like so:
```bash
lx cmd/lx/main.go | wl-copy
```
This will place the entire file in clipboard.


## More examples

Glob multiple files:
```bash
lx **/*.py
```
This will output each file after each other.

Glob multiple files except tests (zsh):
```bash
lx **/*.py~*_test.py
```

Copy from stdin is very convenient with `ripgrep -l` mode. It will search for 
files in recursive directories and output files with matches.

```bash
rg "def save" -l **/database/*.py | lx
```

Custom delimiters:

```
lx --prefix-delimiter="### {filename}\n```{language}\n" file.go
```
This example would make more sense as an alias.

## Placeholders

`{filename}`
`{row_count}`
`{byte_size}`
`{last_modified}`
`{language}`



