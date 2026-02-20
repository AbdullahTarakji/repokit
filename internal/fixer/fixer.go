package fixer

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/AbdullahTarakji/repokit/internal/analyzer"
)

//go:embed templates/*
var templateFS embed.FS

// FixResult describes what was fixed.
type FixResult struct {
	FilesCreated []string
	Errors       []error
}

// TemplateData holds data passed to templates.
type TemplateData struct {
	ProjectName string
	Year        int
	AuthorName  string
	Language    string
	Languages   []string
}

// Fix applies fixes based on analysis results.
func Fix(result *analyzer.AnalysisResult, repoPath string) (*FixResult, error) {
	fr := &FixResult{}
	lang := analyzer.PrimaryLanguage(result.Languages)
	data := TemplateData{
		ProjectName: result.RepoName,
		Year:        time.Now().Year(),
		AuthorName:  "Contributors",
		Language:    lang,
		Languages:   result.Languages,
	}

	if !result.HasLicense {
		fr.generate(repoPath, "LICENSE", "templates/LICENSE.tmpl", data)
	}
	if !result.HasGitignore {
		tmpl := gitignoreTemplate(lang)
		fr.generate(repoPath, ".gitignore", tmpl, data)
	}
	if !result.HasContributing {
		fr.generate(repoPath, "CONTRIBUTING.md", "templates/CONTRIBUTING.md.tmpl", data)
	}
	if !result.HasChangelog {
		fr.generate(repoPath, "CHANGELOG.md", "templates/CHANGELOG.md.tmpl", data)
	}
	if !result.HasCodeOfConduct {
		fr.generate(repoPath, "CODE_OF_CONDUCT.md", "templates/CODE_OF_CONDUCT.md.tmpl", data)
	}
	if !result.HasSecurity {
		fr.generate(repoPath, "SECURITY.md", "templates/SECURITY.md.tmpl", data)
	}
	if !result.HasEditorconfig {
		fr.generate(repoPath, ".editorconfig", "templates/editorconfig.tmpl", data)
	}
	if len(result.CIFiles) == 0 {
		tmpl := ciTemplate(lang)
		fr.generateWithDir(repoPath, filepath.Join(".github", "workflows", "ci.yml"), tmpl, data)
	}
	if len(result.IssueTemplates) == 0 {
		fr.generateWithDir(repoPath, filepath.Join(".github", "ISSUE_TEMPLATE", "bug_report.md"), "templates/issue-templates/bug_report.md.tmpl", data)
		fr.generateWithDir(repoPath, filepath.Join(".github", "ISSUE_TEMPLATE", "feature_request.md"), "templates/issue-templates/feature_request.md.tmpl", data)
	}
	if !result.HasPRTemplate {
		fr.generateWithDir(repoPath, filepath.Join(".github", "pull_request_template.md"), "templates/pull_request_template.md.tmpl", data)
	}

	return fr, nil
}

func (fr *FixResult) generate(repoPath, filename, tmplPath string, data TemplateData) {
	fr.generateWithDir(repoPath, filename, tmplPath, data)
}

func (fr *FixResult) generateWithDir(repoPath, filename, tmplPath string, data TemplateData) {
	content, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		fr.Errors = append(fr.Errors, fmt.Errorf("read template %s: %w", tmplPath, err))
		return
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(content))
	if err != nil {
		fr.Errors = append(fr.Errors, fmt.Errorf("parse template %s: %w", tmplPath, err))
		return
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		fr.Errors = append(fr.Errors, fmt.Errorf("execute template %s: %w", tmplPath, err))
		return
	}

	fullPath := filepath.Join(repoPath, filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		fr.Errors = append(fr.Errors, fmt.Errorf("create dir for %s: %w", filename, err))
		return
	}

	if err := os.WriteFile(fullPath, []byte(buf.String()), 0o644); err != nil {
		fr.Errors = append(fr.Errors, fmt.Errorf("write %s: %w", filename, err))
		return
	}

	fr.FilesCreated = append(fr.FilesCreated, filename)
}

func gitignoreTemplate(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		return "templates/gitignore/go.tmpl"
	case "javascript", "typescript":
		return "templates/gitignore/node.tmpl"
	case "python":
		return "templates/gitignore/python.tmpl"
	case "rust":
		return "templates/gitignore/rust.tmpl"
	default:
		return "templates/gitignore/default.tmpl"
	}
}

func ciTemplate(lang string) string {
	switch strings.ToLower(lang) {
	case "go":
		return "templates/ci/go.yml.tmpl"
	case "javascript", "typescript":
		return "templates/ci/node.yml.tmpl"
	case "python":
		return "templates/ci/python.yml.tmpl"
	default:
		return "templates/ci/default.yml.tmpl"
	}
}
