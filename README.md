# lx

[![Go Reference](https://pkg.go.dev/badge/github.com/rasros/lx.svg)](https://pkg.go.dev/github.com/rasros/lx)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasros/lx)](https://goreportcard.com/report/github.com/rasros/lx)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`lx` is a small CLI for turning one or more files into clean, Markdown-fenced blocks that are easy to paste into LLM chat windows.

The purpose is not just formatting — it's workflow control. `lx` lets you **capture prompt context in your shell**, so you can recreate it instantly. Instead of repairing a misaligned AI conversation, you restart with a refined prompt in seconds. Your context becomes a repeatable artifact, not a manual copy-paste ritual.

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

````text
cmd/lx/main.go (18 rows)
---
```go
package main
... (rest omitted)
````

````

Copy directly to clipboard (Wayland):

```bash
lx cmd/lx/main.go | wl-copy
````

Use cases:

* Create a “context recipe” you can re-run for any new chat session.
* Quickly modify the recipe (more files, different slicing) and re-paste into a clean conversation.
* Treat prompts as reproducible commands instead of fragile in-chat constructions.

---

## More examples

### Multiple files and globs

```bash
lx **/*.py
```

Exclude tests (zsh):

```bash
lx **/*.py~*_test.py
```

Collect files with ripgrep, then format:

```bash
rg "def save" -l **/database/*.py | lx
```

This produces a stable, shell-based definition of “what I want the model to see,” easily re-run anytime.

---

### Slicing only what matters

Keep prompts lean and focused:

```bash
# First 40 lines
lx -h40 server.log

# Last 80 lines
lx -t80 server.log

# Both ends (split 60 lines between head & tail)
lx -n60 server.log
```

Short forms like `-h5`, `-t10`, `-n2` are automatically normalized.

Slicing lives in your command, not in manual editing. Adjust and re-run to refine context.

---

### Line numbers for precise references

```bash
lx -l file.go
```

Or combined with slicing:

```bash
lx --line-numbers -n80 -t40 file.go
```

Line numbers make follow-up instructions concrete (“check line 37–45”) and consistent across restarts.

Behavior summary:

* Full file: numbered `1..N`.
* Head only: numbers start at `1`.
* Tail only: numbers reflect original line positions.
* Head + tail: middle replaced by an unnumbered ellipsis.

---

## Features

* Generates Markdown headers and fenced blocks for one or many files.
* Automatically detects fenced-code language from file extension.
* Supports lightweight, ergonomic slicing (`-h`, `-t`, `-n`).
* Optional line numbers for precise AI instructions.
* Reads filenames from CLI args or stdin (great with `rg`, `find`, etc.).
* Customizable delimiters with placeholders.
* Designed for reproducibility: encode your prompt-context workflow in a command or alias instead of building it manually in an AI UI.

---

## Custom delimiters and placeholders

Default delimiters:

````text
{filename} ({row_count} rows)
---
```{language}
...file contents...
````

````

Override them:

```bash
lx \
  --prefix-delimiter="### {filename}\n```{language}\n" \
  --postfix-delimiter="```\n\n" \
  file.go
````

Placeholders:

* `{filename}`
* `{row_count}`
* `{byte_size}`
* `{last_modified}`
* `{language}`

Use these to enforce consistent prompt structure and regenerate identical context anytime.

