# lx

**`lx` is a small CLI for turning one or more files into clean, Markdown-fenced blocks that are easy to paste into LLM chat windows.**

[![Go Reference](https://pkg.go.dev/badge/github.com/rasros/lx.svg)](https://pkg.go.dev/github.com/rasros/lx)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasros/lx)](https://goreportcard.com/report/github.com/rasros/lx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Its purpose is to make prompt setup *repeatable*. Instead of manually selecting, trimming, and pasting code into an AI chat window, you define the context you want in a single shell command. Re-run that command and you get the exact same prompt every time. It's designed to work with popular tools like `rg -l` and `fd` and shell glob.

This creates a more stable workflow. You keep control of the context, not the chat UI. If a conversation drifts or misinterprets something, you restart with the same context immediately.  

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

Use cases:

* Create a "context recipe" you can re-run for any new chat session.
* Quickly modify the recipe (more files, different slicing) and re-paste into a clean conversation.
* Treat prompts as reproducible commands instead of fragile in-chat constructions.

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
```bash

# fd using stdin-mode
fd -e py -E "*_test.py" | lx

# find using stdin-mode
find . -name '*.py' ! -name '*_test.py' | lx

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
Searching for patterns is easily done with `ripgrep`. This example search for a function beginning with "save" in the `database` folder.

Collect files with ripgrep, then format:
```bash
rg "def save" -l **/database/*.py | lx
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

