package lx

import (
	"path/filepath"
	"strings"
)

var extToLang = map[string]string{
	".go":   "go",
	".py":   "python",
	".js":   "javascript",
	".jsx":  "jsx",
	".ts":   "typescript",
	".tsx":  "tsx",
	".rs":   "rust",
	".java": "java",

	".c": "c",
	".h": "c",

	".cc":  "cpp",
	".cpp": "cpp",
	".cxx": "cpp",
	".hpp": "cpp",
	".hh":  "cpp",
	".hxx": "cpp",

	".sh":   "bash",
	".bash": "bash",
	".zsh":  "zsh",

	".rb":  "ruby",
	".php": "php",

	".html": "html",
	".htm":  "html",
	".css":  "css",

	".json": "json",

	".yml":  "yaml",
	".yaml": "yaml",

	".toml": "toml",

	".md":       "markdown",
	".markdown": "markdown",

	".txt": "text",
}

// languageFromPath returns a markdown language identifier derived from the file extension.
func languageFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	return extToLang[ext]
}
