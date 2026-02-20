package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanSecrets(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantHits int
	}{
		{"AWS key", "AKIAIOSFODNN7EXAMPLE", 1},
		{"GitHub token", "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghij", 1},
		{"Private key", "-----BEGIN RSA PRIVATE KEY-----", 1},
		{"Password", `password = "supersecret123"`, 1},
		{"API key", `api_key = "abcdefgh12345678"`, 1},
		{"Token assign", `token = "my-secret-token-value"`, 1},
		{"Clean file", "package main\n\nfunc main() {}\n", 0},
		{"Short password", `password = "short"`, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			os.WriteFile(filepath.Join(dir, "test.go"), []byte(tc.content), 0o644)
			findings := ScanSecrets(dir)
			if len(findings) != tc.wantHits {
				t.Errorf("expected %d findings, got %d: %+v", tc.wantHits, len(findings), findings)
			}
		})
	}
}

func TestScanSecrets_SkipsBinary(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "image.png"), []byte("AKIAIOSFODNN7EXAMPLE"), 0o644)
	findings := ScanSecrets(dir)
	if len(findings) != 0 {
		t.Error("should skip binary files")
	}
}

func TestScanSecrets_SkipsGitDir(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0o755)
	os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("AKIAIOSFODNN7EXAMPLE"), 0o644)
	findings := ScanSecrets(dir)
	if len(findings) != 0 {
		t.Error("should skip .git directory")
	}
}

func TestScanSecrets_SkipsLargeFiles(t *testing.T) {
	dir := t.TempDir()
	large := make([]byte, 2*1024*1024)
	copy(large, []byte("AKIAIOSFODNN7EXAMPLE"))
	os.WriteFile(filepath.Join(dir, "big.txt"), large, 0o644)
	findings := ScanSecrets(dir)
	if len(findings) != 0 {
		t.Error("should skip large files")
	}
}

func TestScanSecrets_SkipsTestFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "main_test.go"), []byte(`password = "supersecret123"`), 0o644)
	findings := ScanSecrets(dir)
	if len(findings) != 0 {
		t.Error("should skip test files")
	}
}

func TestIsBinaryExtension(t *testing.T) {
	tests := []struct {
		ext    string
		binary bool
	}{
		{".png", true},
		{".go", false},
		{".exe", true},
		{".zip", true},
		{".txt", false},
	}
	for _, tc := range tests {
		if got := isBinaryExtension(tc.ext); got != tc.binary {
			t.Errorf("isBinaryExtension(%q) = %v, want %v", tc.ext, got, tc.binary)
		}
	}
}
