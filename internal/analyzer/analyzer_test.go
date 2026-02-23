package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	// Create .git dir to make it look like a repo
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestAnalyze_EmptyRepo(t *testing.T) {
	dir := setupTestRepo(t, nil)
	result, err := Analyze(dir)
	if err != nil {
		t.Fatal(err)
	}
	if result.HasREADME {
		t.Error("expected no README")
	}
	if result.HasLicense {
		t.Error("expected no LICENSE")
	}
}

func TestAnalyze_WithFiles(t *testing.T) {
	dir := setupTestRepo(t, map[string]string{
		"README.md":        "# Hello\n\n## Installation\n\ninstall stuff\n\n## Usage\n\nuse it\n\n[![badge](url)]\n\n[link](url)\n\n```go\ncode\n```\n" + string(make([]byte, 500)),
		"LICENSE":          "MIT",
		"CONTRIBUTING.md":  "contrib",
		"CHANGELOG.md":     "changes",
		"CODE_OF_CONDUCT.md": "conduct",
		"SECURITY.md":      "security",
		".gitignore":       "*.exe\n/vendor\nnode_modules\n",
		"Makefile":         "build:",
		".editorconfig":    "root=true",
		".gitattributes":   "* text=auto",
		"go.mod":           "module test\n\ngo 1.21",
		"main.go":          "package main",
	})

	result, err := Analyze(dir)
	if err != nil {
		t.Fatal(err)
	}

	checks := []struct {
		name string
		got  bool
	}{
		{"HasREADME", result.HasREADME},
		{"HasLicense", result.HasLicense},
		{"HasContributing", result.HasContributing},
		{"HasChangelog", result.HasChangelog},
		{"HasCodeOfConduct", result.HasCodeOfConduct},
		{"HasSecurity", result.HasSecurity},
		{"HasGitignore", result.HasGitignore},
		{"HasMakefile", result.HasMakefile},
		{"HasEditorconfig", result.HasEditorconfig},
		{"HasGitattributes", result.HasGitattributes},
		{"READMEHasBadges", result.READMEHasBadges},
		{"READMEHasInstall", result.READMEHasInstall},
		{"READMEHasUsage", result.READMEHasUsage},
		{"READMEHasHeadings", result.READMEHasHeadings},
		{"READMEHasCodeBlocks", result.READMEHasCodeBlocks},
		{"READMEHasLinks", result.READMEHasLinks},
		{"README > 500 chars", result.READMELength > 500},
	}

	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.got {
				t.Errorf("%s should be true", tc.name)
			}
		})
	}
}

func TestAnalyze_InvalidPath(t *testing.T) {
	_, err := Analyze("/nonexistent/path")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestAnalyze_NotADirectory(t *testing.T) {
	f, _ := os.CreateTemp("", "repokit-test")
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()
	_, err := Analyze(f.Name())
	if err == nil {
		t.Error("expected error for non-directory")
	}
}

func TestAnalyze_CIDetection(t *testing.T) {
	dir := setupTestRepo(t, map[string]string{
		".github/workflows/ci.yml": "name: CI\non:\n  push:\n  pull_request:\njobs:\n  test:\n    runs-on: ubuntu-latest\n    strategy:\n      matrix:\n        go: [1.21, 1.22]\n    steps:\n      - run: go test ./...\n      - run: golangci-lint run\n      - run: go build ./...\n",
	})

	result, err := Analyze(dir)
	if err != nil {
		t.Fatal(err)
	}

	ciChecks := []struct {
		name string
		got  bool
	}{
		{"CIFiles found", len(result.CIFiles) > 0},
		{"CIHasTest", result.CIHasTest},
		{"CIHasLint", result.CIHasLint},
		{"CIHasBuild", result.CIHasBuild},
		{"CIHasMatrix", result.CIHasMatrix},
		{"CITriggerPushPR", result.CITriggerPushPR},
	}

	for _, tc := range ciChecks {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.got {
				t.Errorf("%s should be true", tc.name)
			}
		})
	}
}

func TestAnalyze_CommunityFiles(t *testing.T) {
	dir := setupTestRepo(t, map[string]string{
		".github/ISSUE_TEMPLATE/bug.md":          "bug",
		".github/ISSUE_TEMPLATE/feature.md":      "feature",
		".github/pull_request_template.md":        "pr",
		"CODEOWNERS":                              "* @owner",
		".github/FUNDING.yml":                     "github: user",
	})

	result, err := Analyze(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.IssueTemplates) != 2 {
		t.Errorf("expected 2 issue templates, got %d", len(result.IssueTemplates))
	}
	if !result.HasPRTemplate {
		t.Error("expected PR template")
	}
	if !result.HasCodeowners {
		t.Error("expected CODEOWNERS")
	}
	if !result.HasFundingYML {
		t.Error("expected FUNDING.yml")
	}
}

func TestAnalyze_LockFiles(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"go.sum", "go.sum"},
		{"package-lock.json", "package-lock.json"},
		{"yarn.lock", "yarn.lock"},
		{"Cargo.lock", "Cargo.lock"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := setupTestRepo(t, map[string]string{tc.file: "content"})
			result, err := Analyze(dir)
			if err != nil {
				t.Fatal(err)
			}
			if !result.HasLockFile {
				t.Errorf("expected HasLockFile for %s", tc.file)
			}
		})
	}
}

func TestFileExistsAny(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "README.md"), []byte("hi"), 0o644)

	if !fileExistsAny(dir, "README.md") {
		t.Error("should find README.md")
	}
	if fileExistsAny(dir, "NOPE") {
		t.Error("should not find NOPE")
	}
}

func TestContainsAnyCI(t *testing.T) {
	tests := []struct {
		s      string
		subs   []string
		expect bool
	}{
		{"## Installation", []string{"## install"}, true},
		{"nothing here", []string{"## install"}, false},
		{"## USAGE", []string{"## usage"}, true},
	}
	for _, tc := range tests {
		got := containsAnyCI(tc.s, tc.subs...)
		if got != tc.expect {
			t.Errorf("containsAnyCI(%q, %v) = %v, want %v", tc.s, tc.subs, got, tc.expect)
		}
	}
}

func TestGitignoreCoversLanguages(t *testing.T) {
	tests := []struct {
		content string
		langs   []string
		expect  bool
	}{
		{"node_modules/\n", []string{"JavaScript"}, true},
		{"*.exe\n/vendor\n", []string{"Go"}, true},
		{"__pycache__/\n", []string{"Python"}, true},
		{"/target\n", []string{"Rust"}, true},
		{"", []string{"Go"}, false},
	}
	for _, tc := range tests {
		got := gitignoreCoversLanguages(tc.content, tc.langs)
		if got != tc.expect {
			t.Errorf("gitignoreCoversLanguages(%q, %v) = %v, want %v", tc.content, tc.langs, got, tc.expect)
		}
	}
}
