package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AnalysisResult holds all collected facts about a repository.
type AnalysisResult struct {
	RepoPath string
	RepoName string

	// Files found
	HasREADME        bool
	HasLicense       bool
	HasContributing  bool
	HasChangelog     bool
	HasCodeOfConduct bool
	HasSecurity      bool
	HasDocsDir       bool

	// README quality
	READMELength     int
	READMEHasBadges  bool
	READMEHasInstall bool
	READMEHasUsage   bool
	READMEHasHeadings  bool
	READMEHasCodeBlocks bool
	READMEHasLinks   bool

	// CI
	CIFiles         []string
	CIHasTest       bool
	CIHasLint       bool
	CIHasBuild      bool
	CIHasMatrix     bool
	CITriggerPushPR bool
	HasMakefile     bool

	// Security
	HasGitignore         bool
	GitignoreCoversLang  bool
	SecretsFound         []SecretFinding
	HasSecurityMD        bool
	HasDependabot        bool
	HasEnvFiles          bool
	HasPreCommit         bool

	// Community
	IssueTemplates    []string
	HasPRTemplate     bool
	HasCodeowners     bool
	HasFundingYML     bool

	// Maintenance
	LastCommitDate  time.Time
	HasReleases     bool
	StaleBranches   int
	HasEditorconfig bool
	HasLockFile     bool
	HasGitattributes bool
	HasVersionFile  bool
	RepoSizeKB     int64

	// Languages
	Languages []string
}

// Analyze scans a local repository path and returns analysis results.
func Analyze(repoPath string) (*AnalysisResult, error) {
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", absPath)
	}

	result := &AnalysisResult{
		RepoPath: absPath,
		RepoName: filepath.Base(absPath),
	}

	checkFiles(result)
	analyzeREADME(result)
	analyzeCI(result)
	analyzeSecurity(result)
	analyzeCommunity(result)
	analyzeMaintenance(result)
	result.Languages = DetectLanguages(absPath)

	return result, nil
}

func checkFiles(r *AnalysisResult) {
	r.HasREADME = fileExistsAny(r.RepoPath, "README.md", "README", "README.txt", "readme.md")
	r.HasLicense = fileExistsAny(r.RepoPath, "LICENSE", "LICENSE.md", "LICENSE.txt", "LICENCE", "LICENCE.md")
	r.HasContributing = fileExistsAny(r.RepoPath, "CONTRIBUTING.md", "CONTRIBUTING", "contributing.md")
	r.HasChangelog = fileExistsAny(r.RepoPath, "CHANGELOG.md", "CHANGELOG", "HISTORY.md", "CHANGES.md")
	r.HasCodeOfConduct = fileExistsAny(r.RepoPath, "CODE_OF_CONDUCT.md", "code_of_conduct.md")
	r.HasSecurity = fileExistsAny(r.RepoPath, "SECURITY.md", "security.md")
	r.HasDocsDir = dirExists(r.RepoPath, "docs")
	r.HasMakefile = fileExistsAny(r.RepoPath, "Makefile", "makefile", "GNUmakefile")
	r.HasEditorconfig = fileExistsAny(r.RepoPath, ".editorconfig")
	r.HasGitattributes = fileExistsAny(r.RepoPath, ".gitattributes")
}

func analyzeREADME(r *AnalysisResult) {
	if !r.HasREADME {
		return
	}

	content := readFileContent(r.RepoPath, "README.md")
	if content == "" {
		content = readFileContent(r.RepoPath, "readme.md")
	}
	if content == "" {
		content = readFileContent(r.RepoPath, "README")
	}

	r.READMELength = len(content)
	r.READMEHasBadges = strings.Contains(content, "[![") || strings.Contains(content, "![badge")
	r.READMEHasInstall = containsAnyCI(content, "## install", "## installation", "## getting started", "## setup", "## quick start")
	r.READMEHasUsage = containsAnyCI(content, "## usage", "## examples", "## how to use")
	r.READMEHasHeadings = strings.Contains(content, "# ")
	r.READMEHasCodeBlocks = strings.Contains(content, "```")
	r.READMEHasLinks = strings.Contains(content, "](")
}

func analyzeCI(r *AnalysisResult) {
	workflowDir := filepath.Join(r.RepoPath, ".github", "workflows")
	entries, err := os.ReadDir(workflowDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if !e.IsDir() && (strings.HasSuffix(e.Name(), ".yml") || strings.HasSuffix(e.Name(), ".yaml")) {
			r.CIFiles = append(r.CIFiles, e.Name())
			content := readFileContent(workflowDir, e.Name())
			lc := strings.ToLower(content)
			if strings.Contains(lc, "test") {
				r.CIHasTest = true
			}
			if strings.Contains(lc, "lint") || strings.Contains(lc, "golangci") || strings.Contains(lc, "eslint") || strings.Contains(lc, "flake8") || strings.Contains(lc, "clippy") {
				r.CIHasLint = true
			}
			if strings.Contains(lc, "build") || strings.Contains(lc, "go build") || strings.Contains(lc, "npm run build") {
				r.CIHasBuild = true
			}
			if strings.Contains(lc, "matrix") {
				r.CIHasMatrix = true
			}
			if strings.Contains(lc, "pull_request") && (strings.Contains(lc, "push") || strings.Contains(lc, "on:")) {
				r.CITriggerPushPR = true
			}
		}
	}
}

func analyzeSecurity(r *AnalysisResult) {
	r.HasGitignore = fileExistsAny(r.RepoPath, ".gitignore")
	if r.HasGitignore {
		content := readFileContent(r.RepoPath, ".gitignore")
		langs := DetectLanguages(r.RepoPath)
		r.GitignoreCoversLang = gitignoreCoversLanguages(content, langs)
	}

	r.SecretsFound = ScanSecrets(r.RepoPath)
	r.HasSecurityMD = r.HasSecurity
	r.HasDependabot = fileExistsAny(filepath.Join(r.RepoPath, ".github"), "dependabot.yml", "dependabot.yaml")
	r.HasEnvFiles = checkEnvFiles(r.RepoPath)
	r.HasPreCommit = fileExistsAny(r.RepoPath, ".pre-commit-config.yaml", ".pre-commit-config.yml")
}

func analyzeCommunity(r *AnalysisResult) {
	// Issue templates
	templateDir := filepath.Join(r.RepoPath, ".github", "ISSUE_TEMPLATE")
	entries, err := os.ReadDir(templateDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				r.IssueTemplates = append(r.IssueTemplates, e.Name())
			}
		}
	}

	r.HasPRTemplate = fileExistsAny(r.RepoPath, ".github/pull_request_template.md", ".github/PULL_REQUEST_TEMPLATE.md") ||
		fileExistsAny(filepath.Join(r.RepoPath, ".github"), "pull_request_template.md", "PULL_REQUEST_TEMPLATE.md")
	r.HasCodeowners = fileExistsAny(r.RepoPath, "CODEOWNERS", ".github/CODEOWNERS", "docs/CODEOWNERS")
	r.HasFundingYML = fileExistsAny(filepath.Join(r.RepoPath, ".github"), "FUNDING.yml", "FUNDING.yaml")
}

func analyzeMaintenance(r *AnalysisResult) {
	// Last commit date - check .git directory
	gitHead := filepath.Join(r.RepoPath, ".git", "HEAD")
	if info, err := os.Stat(gitHead); err == nil {
		r.LastCommitDate = info.ModTime()
	}

	// Releases/tags
	tagsDir := filepath.Join(r.RepoPath, ".git", "refs", "tags")
	if entries, err := os.ReadDir(tagsDir); err == nil && len(entries) > 0 {
		r.HasReleases = true
	}
	// Also check packed-refs
	if !r.HasReleases {
		packedRefs := readFileContent(filepath.Join(r.RepoPath, ".git"), "packed-refs")
		if strings.Contains(packedRefs, "refs/tags/") {
			r.HasReleases = true
		}
	}

	// Lock files
	r.HasLockFile = fileExistsAny(r.RepoPath, "go.sum", "package-lock.json", "yarn.lock", "poetry.lock", "Cargo.lock", "Gemfile.lock", "pnpm-lock.yaml", "composer.lock")

	// Version file
	r.HasVersionFile = fileExistsAny(r.RepoPath, "VERSION", "version.txt", ".version")

	// Repo size
	r.RepoSizeKB = calcDirSizeKB(r.RepoPath)
}

func gitignoreCoversLanguages(content string, langs []string) bool {
	lc := strings.ToLower(content)
	for _, lang := range langs {
		switch strings.ToLower(lang) {
		case "go":
			if strings.Contains(lc, "*.exe") || strings.Contains(lc, "/vendor") || strings.Contains(lc, "bin/") {
				return true
			}
		case "javascript", "typescript":
			if strings.Contains(lc, "node_modules") {
				return true
			}
		case "python":
			if strings.Contains(lc, "__pycache__") || strings.Contains(lc, "*.pyc") {
				return true
			}
		case "rust":
			if strings.Contains(lc, "/target") {
				return true
			}
		default:
			// If gitignore has any content, give partial credit
			if len(strings.TrimSpace(content)) > 10 {
				return true
			}
		}
	}
	return false
}

func checkEnvFiles(repoPath string) bool {
	patterns := []string{".env", ".env.local", ".env.production"}
	for _, p := range patterns {
		if fileExistsAny(repoPath, p) {
			return true
		}
	}
	return false
}

func fileExistsAny(dir string, names ...string) bool {
	for _, name := range names {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

func dirExists(base, name string) bool {
	info, err := os.Stat(filepath.Join(base, name))
	return err == nil && info.IsDir()
}

func readFileContent(dir, name string) string {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return ""
	}
	return string(data)
}

func containsAnyCI(s string, substrs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func calcDirSizeKB(root string) int64 {
	var size int64
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && filepath.Base(p) == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size / 1024
}
