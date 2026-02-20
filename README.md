# 🏥 RepoKit

[![CI](https://github.com/AbdullahTarakji/repokit/actions/workflows/ci.yml/badge.svg)](https://github.com/AbdullahTarakji/repokit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/AbdullahTarakji/repokit)](https://goreportcard.com/report/github.com/AbdullahTarakji/repokit)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A CLI tool that audits GitHub repositories for health and quality, assigns a score (0-100), and can auto-fix missing files.

![Demo](docs/demo.gif)

## Installation

```bash
go install github.com/AbdullahTarakji/repokit/cmd/repokit@latest
```

Or build from source:

```bash
git clone https://github.com/AbdullahTarakji/repokit.git
cd repokit
make build
```

## Usage

```bash
# Scan current directory (TUI mode)
repokit scan .

# Scan with plain text output
repokit scan --format text .

# JSON output (great for CI)
repokit scan --format json .

# Auto-fix missing files
repokit fix .

# Preview what would be fixed (dry run)
repokit fix --dry-run .

# Legacy: scan and fix in one command
repokit --fix --yes /path/to/repo

# Print version
repokit version
```

## Scoring

RepoKit evaluates repositories across 5 categories, each scored 0-20 for a total of 0-100:

### 📝 Documentation (0-20)
- README.md exists and has quality signals (badges, install/usage sections, headings, code blocks, links)
- LICENSE, CONTRIBUTING.md, CHANGELOG.md, CODE_OF_CONDUCT.md, docs/ directory

### ⚙️ CI/CD (0-20)
- GitHub Actions workflows with test, lint, and build steps
- Matrix strategy, push+PR triggers, Makefile

### 🔒 Security (0-20)
- .gitignore with language-specific patterns
- No secrets detected in tracked files (AWS keys, API tokens, private keys, passwords)
- SECURITY.md, Dependabot config, no .env files, pre-commit hooks

### 👥 Community (0-20)
- Issue templates, PR template, CODEOWNERS, FUNDING.yml
- Remote-only: description, topics, discussions, stars

### 🔧 Maintenance (0-20)
- Recent commits, releases/tags, no stale branches
- .editorconfig, lock files, reasonable repo size, .gitattributes

### Score Colors
- 🟢 **Green**: 80%+ (16-20 per category)
- 🟡 **Yellow**: 50-79% (10-15 per category)
- 🔴 **Red**: Below 50% (0-9 per category)

## Auto-Fix

When you run `repokit --fix`, it generates missing files from templates:

- **LICENSE** (MIT)
- **.gitignore** (language-specific)
- **CONTRIBUTING.md** (language-aware setup instructions)
- **CHANGELOG.md**
- **CODE_OF_CONDUCT.md** (Contributor Covenant v2.1)
- **SECURITY.md**
- **.editorconfig**
- **CI workflow** (language-specific GitHub Actions)
- **Issue templates** (bug report + feature request)
- **PR template**

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
