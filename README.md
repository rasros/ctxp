# lx

[![Go Reference](https://pkg.go.dev/badge/github.com/rasros/lx.svg)](https://pkg.go.dev/github.com/rasros/lx)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasros/lx)](https://goreportcard.com/report/github.com/rasros/lx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**`lx` formats one or more files into clean, Markdown-fenced blocks ready to paste into LLM chat windows.**

The goal is to make prompt setup **repeatable**. Instead of manually letting an agent guess context or clicking through files in a UI, you define the exact context you want in one shell command and rerun it whenever you need a fresh session. It works smoothly with tools like `rg -l`, `fd`, and recursive shell globs.

This gives you a stable, controllable workflow:

- You decide exactly what context the model sees.
- If a conversation drifts, you restart instantly with the same context.
- Adjust the command, rerun it, and paste the updated output.

---

## Installation

```bash
go install github.com/rasros/lx/cmd/lx@latest
````

---

## Basic usage

Format a single file as an LLM-ready snippet:

```bash
lx cmd/lx/main.go
```

Example output (trimmed):

~~~text
cmd/lx/main.go (18 rows)
---
```go
package main
... (rest omitted)
```
~~~

Copy directly to clipboard (Wayland):

```bash
lx cmd/lx/main.go | wl-copy
```

---

## Features

* Generates Markdown headers and fenced blocks for one or many files.
* Automatically detects fenced-code language from file extension.
* Supports lightweight, ergonomic slicing (`-h`, `-t`, `-n`).
* Optional line numbers for precise AI instructions.
* Reads filenames from CLI args or stdin (great with `rg`, `find`, etc.).
* Customizable delimiters with placeholders.

---

## More examples

### Multiple files with glob
We can easily filter multiple files in shells that allow recursive glob:
```bash
lx **/*.py
```

`lx` relies on standard tools for file selection. This example includes all python files except those with name ending in `_test.py`:

Either do it through `fd` or `find`:
```bash
# fd using stdin-mode
fd -e py -E "*_test.py" | lx

# find using stdin-mode
find . -name '*.py' ! -name '*_test.py' | lx
```

Or through glob:
```bash
# zsh glob
lx **/*.py~*_test.py 

# bash glob
shopt -s globstar extglob
lx **/*.py~*_test.py

# fish glob
lx **/*.py ^**/*_test.py

# msys2 glob
shopt -s globstar extglob
lx **/!(*_test).py
```

### Pattern search
Searching for patterns is easily done through grep and pipe matching files to `lx`.

This example searches for files with a function starting with `save` under the src folder:
```bash
# grep
grep -rl "def save" src | lx

# ripgrep
rg -l "def save" src | lx
```

### Slicing

While iterating on the command it's convenient to slice files so you can more easily see what's included:

```bash
# First 40 lines
lx -h40 server.log

# Last 80 lines
lx -t80 server.log

# Both ends (split 60 lines between head & tail)
lx -n60 server.log
```

Short forms like `-h5`, `-t10`, `-n2` are supported.

### Line numbers for precise references

TOON supports line numbers so we do too 🤷.

Useful for non-coding or if you include prompting about line number references.

### Custom delimiters and placeholders

Default delimiters:

~~~text
{filename} ({row_count} rows)
---
```{language}
...file contents...
```
~~~

Override them:

~~~bash
lx \
  --prefix-delimiter="### {filename}\n```{language}\n" \
  --postfix-delimiter="```\n\n" \
  file.go
~~~

Placeholders:

* `{filename}`
* `{row_count}`
* `{byte_size}`
* `{last_modified}`
* `{language}`

Use these to enforce consistent prompt structure and regenerate identical context anytime.

