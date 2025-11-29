package lx

import "testing"

func TestLanguageFromPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "go file",
			path: "main.go",
			want: "go",
		},
		{
			name: "python file",
			path: "script.py",
			want: "python",
		},
		{
			name: "javascript file",
			path: "app.js",
			want: "javascript",
		},
		{
			name: "jsx file",
			path: "component.jsx",
			want: "jsx",
		},
		{
			name: "typescript file",
			path: "types.ts",
			want: "typescript",
		},
		{
			name: "tsx file",
			path: "component.tsx",
			want: "tsx",
		},
		{
			name: "rust file",
			path: "lib.rs",
			want: "rust",
		},
		{
			name: "java file",
			path: "Main.java",
			want: "java",
		},
		{
			name: "c file",
			path: "prog.c",
			want: "c",
		},
		{
			name: "header file c",
			path: "prog.h",
			want: "c",
		},
		{
			name: "cpp file",
			path: "prog.cpp",
			want: "cpp",
		},
		{
			name: "cpp alt extension cc",
			path: "prog.cc",
			want: "cpp",
		},
		{
			name: "cpp header hpp",
			path: "prog.hpp",
			want: "cpp",
		},
		{
			name: "bash script",
			path: "script.sh",
			want: "bash",
		},
		{
			name: "zsh script",
			path: "script.zsh",
			want: "zsh",
		},
		{
			name: "ruby file",
			path: "app.rb",
			want: "ruby",
		},
		{
			name: "php file",
			path: "index.php",
			want: "php",
		},
		{
			name: "html file",
			path: "index.html",
			want: "html",
		},
		{
			name: "htm file",
			path: "index.htm",
			want: "html",
		},
		{
			name: "css file",
			path: "styles.css",
			want: "css",
		},
		{
			name: "json file",
			path: "data.json",
			want: "json",
		},
		{
			name: "yaml file yml",
			path: "config.yml",
			want: "yaml",
		},
		{
			name: "yaml file yaml",
			path: "config.yaml",
			want: "yaml",
		},
		{
			name: "toml file",
			path: "config.toml",
			want: "toml",
		},
		{
			name: "markdown file md",
			path: "README.md",
			want: "markdown",
		},
		{
			name: "markdown file markdown",
			path: "doc.markdown",
			want: "markdown",
		},
		{
			name: "text file",
			path: "notes.txt",
			want: "text",
		},
		{
			name: "uppercase extension",
			path: "COMPONENT.TSX",
			want: "tsx",
		},
		{
			name: "mixed case extension",
			path: "Component.JsX",
			want: "jsx",
		},
		{
			name: "multi dot path uses last ext",
			path: "archive.tar.gz",
			want: "",
		},
		{
			name: "no extension",
			path: "LICENSE",
			want: "",
		},
		{
			name: "hidden file no extension",
			path: ".gitignore",
			want: "",
		},
		{
			name: "unknown extension",
			path: "file.unknown",
			want: "",
		},
	}

	for _, tt := range tests {
		got := languageFromPath(tt.path)
		if got != tt.want {
			t.Errorf("%s: languageFromPath(%q) = %q, want %q",
				tt.name, tt.path, got, tt.want)
		}
	}
}

