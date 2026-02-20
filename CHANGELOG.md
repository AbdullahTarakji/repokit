# Changelog

All notable changes to RepoKit will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-20

### Added
- Local repository analysis (files, README quality, CI, security, community, maintenance)
- 5-category scoring system (Documentation, CI/CD, Security, Community, Maintenance)
- Secret detection (AWS keys, GitHub tokens, private keys, passwords, API keys)
- Language detection for Go, JavaScript, TypeScript, Python, Rust, Ruby, Java, C/C++, PHP, Swift, Docker
- Auto-fix with embedded templates for LICENSE, .gitignore, CONTRIBUTING, CHANGELOG, CODE_OF_CONDUCT, SECURITY, .editorconfig, CI workflows, issue templates, PR template
- TUI report with Bubble Tea (color-coded, progress bars)
- Plain text and JSON output formats
- Cobra CLI with --fix, --format, version subcommand
