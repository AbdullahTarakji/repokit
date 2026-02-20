package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SecretFinding represents a potential secret found in a file.
type SecretFinding struct {
	File    string
	Line    int
	Type    string
	Pattern string
}

type secretPattern struct {
	Name    string
	Pattern *regexp.Regexp
}

var secretPatterns = []secretPattern{
	{"AWS Access Key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{"OpenAI API Key", regexp.MustCompile(`sk-[a-zA-Z0-9]{48}`)},
	{"GitHub Token", regexp.MustCompile(`ghp_[a-zA-Z0-9]{36}`)},
	{"GitHub OAuth", regexp.MustCompile(`gho_[a-zA-Z0-9]{36}`)},
	{"Private Key", regexp.MustCompile(`-----BEGIN (RSA |EC |DSA )?PRIVATE KEY-----`)},
	{"Password Assignment", regexp.MustCompile(`(?i)password\s*=\s*["'][^"']{8,}["']`)},
	{"API Key Assignment", regexp.MustCompile(`(?i)api[_-]?key\s*=\s*["'][^"']{8,}["']`)},
	{"Token Assignment", regexp.MustCompile(`(?i)token\s*=\s*["'][^"']{8,}["']`)},
	{"Secret Assignment", regexp.MustCompile(`(?i)secret\s*=\s*["'][^"']{8,}["']`)},
}

// ScanSecrets scans tracked files for potential secrets.
func ScanSecrets(repoPath string) []SecretFinding {
	var findings []SecretFinding

	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(repoPath, path)

		// Skip .git directory
		if info.IsDir() {
			if filepath.Base(path) == ".git" || filepath.Base(path) == "node_modules" || filepath.Base(path) == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary and large files
		if info.Size() > 1024*1024 { // 1MB
			return nil
		}
		if isBinaryExtension(filepath.Ext(path)) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			for _, sp := range secretPatterns {
				if sp.Pattern.MatchString(line) {
					findings = append(findings, SecretFinding{
						File:    rel,
						Line:    i + 1,
						Type:    sp.Name,
						Pattern: sp.Pattern.String(),
					})
				}
			}
		}

		return nil
	})

	return findings
}

func isBinaryExtension(ext string) bool {
	binExts := map[string]bool{
		".exe": true, ".dll": true, ".so": true, ".dylib": true,
		".bin": true, ".obj": true, ".o": true, ".a": true,
		".zip": true, ".tar": true, ".gz": true, ".bz2": true,
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
		".ico": true, ".svg": true, ".woff": true, ".woff2": true,
		".ttf": true, ".eot": true, ".pdf": true, ".mp3": true,
		".mp4": true, ".mov": true, ".avi": true,
	}
	return binExts[strings.ToLower(ext)]
}
