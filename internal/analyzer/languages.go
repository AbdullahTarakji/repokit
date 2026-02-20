package analyzer

import (
	"os"
	"path/filepath"
	"strings"
)

// Language extension mappings.
var extToLang = map[string]string{
	".go":    "Go",
	".js":    "JavaScript",
	".jsx":   "JavaScript",
	".ts":    "TypeScript",
	".tsx":   "TypeScript",
	".py":    "Python",
	".rs":    "Rust",
	".rb":    "Ruby",
	".java":  "Java",
	".c":     "C",
	".cpp":   "C++",
	".h":     "C",
	".hpp":   "C++",
	".php":   "PHP",
	".swift": "Swift",
}

// Config file to language mappings.
var configToLang = map[string]string{
	"go.mod":          "Go",
	"go.sum":          "Go",
	"package.json":    "JavaScript",
	"tsconfig.json":   "TypeScript",
	"pyproject.toml":  "Python",
	"setup.py":        "Python",
	"requirements.txt": "Python",
	"Cargo.toml":      "Rust",
	"Gemfile":         "Ruby",
	"pom.xml":         "Java",
	"build.gradle":    "Java",
	"composer.json":   "PHP",
	"Package.swift":   "Swift",
	"Dockerfile":      "Docker",
}

// DetectLanguages scans a repository and returns detected languages.
func DetectLanguages(repoPath string) []string {
	seen := make(map[string]bool)

	// Check config files first
	for file, lang := range configToLang {
		if _, err := os.Stat(filepath.Join(repoPath, file)); err == nil {
			seen[lang] = true
		}
	}

	// Walk files for extensions (limit depth to avoid huge repos)
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(repoPath, path)

		// Skip hidden dirs, vendor, node_modules
		if info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules" || base == "__pycache__" || base == "target" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip deep paths
		if strings.Count(rel, string(os.PathSeparator)) > 5 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if lang, ok := extToLang[ext]; ok {
			seen[lang] = true
		}

		// Check for Dockerfile
		if filepath.Base(path) == "Dockerfile" {
			seen["Docker"] = true
		}

		return nil
	})

	langs := make([]string, 0, len(seen))
	for lang := range seen {
		langs = append(langs, lang)
	}
	return langs
}

// PrimaryLanguage returns the most likely primary language.
func PrimaryLanguage(langs []string) string {
	if len(langs) == 0 {
		return ""
	}
	// Prefer non-Docker, non-config languages
	for _, l := range langs {
		if l != "Docker" {
			return l
		}
	}
	return langs[0]
}
