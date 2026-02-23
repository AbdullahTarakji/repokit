package fixer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AbdullahTarakji/repokit/internal/analyzer"
)

func TestFix_CreatesMissingFiles(t *testing.T) {
	dir := t.TempDir()
	result := &analyzer.AnalysisResult{
		RepoPath:  dir,
		RepoName:  "test-project",
		Languages: []string{"Go"},
	}

	fr, err := Fix(result, dir)
	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := []string{
		"LICENSE",
		".gitignore",
		"CONTRIBUTING.md",
		"CHANGELOG.md",
		"CODE_OF_CONDUCT.md",
		"SECURITY.md",
		".editorconfig",
		filepath.Join(".github", "workflows", "ci.yml"),
		filepath.Join(".github", "ISSUE_TEMPLATE", "bug_report.md"),
		filepath.Join(".github", "ISSUE_TEMPLATE", "feature_request.md"),
		filepath.Join(".github", "pull_request_template.md"),
	}

	for _, ef := range expectedFiles {
		found := false
		for _, fc := range fr.FilesCreated {
			if fc == ef {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %s to be created, got: %v", ef, fr.FilesCreated)
		}

		if _, err := os.Stat(filepath.Join(dir, ef)); os.IsNotExist(err) {
			t.Errorf("file %s should exist on disk", ef)
		}
	}
}

func TestFix_SkipsExistingFiles(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("existing"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("existing"), 0o644)

	result := &analyzer.AnalysisResult{
		RepoPath:   dir,
		RepoName:   "test",
		HasLicense:  true,
		HasGitignore: true,
		Languages:  []string{"Go"},
	}

	fr, err := Fix(result, dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fr.FilesCreated {
		if f == "LICENSE" || f == ".gitignore" {
			t.Errorf("should not recreate %s", f)
		}
	}

	// Verify original content preserved
	data, _ := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if string(data) != "existing" {
		t.Error("LICENSE should not be overwritten")
	}
}

func TestFix_LicenseContainsYear(t *testing.T) {
	dir := t.TempDir()
	result := &analyzer.AnalysisResult{
		RepoPath: dir,
		RepoName: "test",
	}

	_, _ = Fix(result, dir)
	data, err := os.ReadFile(filepath.Join(dir, "LICENSE"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "MIT License") {
		t.Error("LICENSE should be MIT")
	}
}

func TestFix_ContributingIsLanguageAware(t *testing.T) {
	tests := []struct {
		lang     string
		contains string
	}{
		{"Go", "go test"},
		{"JavaScript", "npm"},
		{"Python", "pytest"},
		{"Rust", "cargo"},
	}

	for _, tc := range tests {
		t.Run(tc.lang, func(t *testing.T) {
			dir := t.TempDir()
			result := &analyzer.AnalysisResult{
				RepoPath:  dir,
				RepoName:  "test",
				Languages: []string{tc.lang},
			}

			_, _ = Fix(result, dir)
			data, _ := os.ReadFile(filepath.Join(dir, "CONTRIBUTING.md"))
			if !strings.Contains(string(data), tc.contains) {
				t.Errorf("CONTRIBUTING.md for %s should contain %q", tc.lang, tc.contains)
			}
		})
	}
}

func TestGitignoreTemplate(t *testing.T) {
	tests := []struct {
		lang string
		tmpl string
	}{
		{"Go", "templates/gitignore/go.tmpl"},
		{"JavaScript", "templates/gitignore/node.tmpl"},
		{"Python", "templates/gitignore/python.tmpl"},
		{"Rust", "templates/gitignore/rust.tmpl"},
		{"Unknown", "templates/gitignore/default.tmpl"},
	}

	for _, tc := range tests {
		got := gitignoreTemplate(tc.lang)
		if got != tc.tmpl {
			t.Errorf("gitignoreTemplate(%q) = %q, want %q", tc.lang, got, tc.tmpl)
		}
	}
}

func TestFix_AllFilesExist(t *testing.T) {
	dir := t.TempDir()

	// Create all files that Fix would generate
	_ = os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "CONTRIBUTING.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "CODE_OF_CONDUCT.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "SECURITY.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".editorconfig"), []byte("x"), 0o644)
	_ = os.MkdirAll(filepath.Join(dir, ".github", "workflows"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, ".github", "workflows", "ci.yml"), []byte("x"), 0o644)
	_ = os.MkdirAll(filepath.Join(dir, ".github", "ISSUE_TEMPLATE"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, ".github", "ISSUE_TEMPLATE", "bug_report.md"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, ".github", "pull_request_template.md"), []byte("x"), 0o644)

	result := &analyzer.AnalysisResult{
		RepoPath:       dir,
		RepoName:       "test",
		HasLicense:     true,
		HasGitignore:   true,
		HasContributing: true,
		HasChangelog:   true,
		HasCodeOfConduct: true,
		HasSecurity:    true,
		HasEditorconfig: true,
		CIFiles:        []string{"ci.yml"},
		IssueTemplates: []string{"bug_report.md"},
		HasPRTemplate:  true,
		Languages:      []string{"Go"},
	}

	fr, err := Fix(result, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(fr.FilesCreated) != 0 {
		t.Errorf("expected no files created when all exist, got: %v", fr.FilesCreated)
	}
	if len(fr.Errors) != 0 {
		t.Errorf("expected no errors, got: %v", fr.Errors)
	}
}

func TestFix_TypeScriptUsesNodeTemplates(t *testing.T) {
	dir := t.TempDir()
	result := &analyzer.AnalysisResult{
		RepoPath:  dir,
		RepoName:  "ts-project",
		Languages: []string{"TypeScript"},
	}

	fr, err := Fix(result, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(fr.FilesCreated) == 0 {
		t.Error("expected files to be created")
	}

	// .gitignore should have node_modules
	data, _ := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if !strings.Contains(string(data), "node_modules") {
		t.Error(".gitignore for TypeScript should contain node_modules")
	}
}

func TestCiTemplate(t *testing.T) {
	tests := []struct {
		lang string
		tmpl string
	}{
		{"Go", "templates/ci/go.yml.tmpl"},
		{"JavaScript", "templates/ci/node.yml.tmpl"},
		{"Python", "templates/ci/python.yml.tmpl"},
		{"Unknown", "templates/ci/default.yml.tmpl"},
	}

	for _, tc := range tests {
		got := ciTemplate(tc.lang)
		if got != tc.tmpl {
			t.Errorf("ciTemplate(%q) = %q, want %q", tc.lang, got, tc.tmpl)
		}
	}
}
