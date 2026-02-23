package analyzer

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestDetectLanguages(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected []string
	}{
		{
			"Go project",
			map[string]string{"go.mod": "module test", "main.go": "package main"},
			[]string{"Go"},
		},
		{
			"Node project",
			map[string]string{"package.json": "{}", "index.js": ""},
			[]string{"JavaScript"},
		},
		{
			"Python project",
			map[string]string{"setup.py": "", "app.py": ""},
			[]string{"Python"},
		},
		{
			"Rust project",
			map[string]string{"Cargo.toml": "", "src/main.rs": "fn main(){}"},
			[]string{"Rust"},
		},
		{
			"Multi-language",
			map[string]string{"main.go": "", "Dockerfile": "FROM alpine"},
			[]string{"Docker", "Go"},
		},
		{
			"TypeScript",
			map[string]string{"tsconfig.json": "{}", "app.tsx": ""},
			[]string{"TypeScript"},
		},
		{
			"Empty repo",
			map[string]string{},
			[]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, content := range tc.files {
				path := filepath.Join(dir, name)
				_ = os.MkdirAll(filepath.Dir(path), 0o755)
				_ = os.WriteFile(path, []byte(content), 0o644)
			}

			got := DetectLanguages(dir)
			sort.Strings(got)
			sort.Strings(tc.expected)

			if len(got) != len(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
				return
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("expected %v, got %v", tc.expected, got)
					break
				}
			}
		})
	}
}

func TestPrimaryLanguage(t *testing.T) {
	tests := []struct {
		langs    []string
		expected string
	}{
		{[]string{"Go", "Docker"}, "Go"},
		{[]string{"Docker"}, "Docker"},
		{[]string{}, ""},
		{[]string{"Python", "JavaScript"}, "Python"},
	}

	for _, tc := range tests {
		got := PrimaryLanguage(tc.langs)
		if got != tc.expected {
			t.Errorf("PrimaryLanguage(%v) = %q, want %q", tc.langs, got, tc.expected)
		}
	}
}
