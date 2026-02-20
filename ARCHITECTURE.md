# RepoKit Architecture

## Overview
RepoKit is a CLI/TUI tool that audits GitHub repositories for health and quality, assigns a score (0-100), and can auto-fix missing files.

## Package Structure

```
cmd/repokit/          — CLI entry point (cobra)
internal/
  analyzer/           — Scans repo, collects facts (what files exist, languages detected, etc.)
  scorer/             — Takes analysis results, computes scores per category and overall
  fixer/              — Generates missing files from templates
  reporter/           — TUI display (Bubble Tea) + plain text output
  github/             — GitHub API client for remote repo analysis
  config/             — Configuration file support (.repokit.yaml)
templates/            — Embedded Go templates for generated files
```

## Core Flow

```
1. User runs: repokit [path|github-url] [--fix] [--format text|json]
2. Analyzer scans the repo → produces AnalysisResult
3. Scorer evaluates AnalysisResult → produces ScoreReport (0-100 + per-category)
4. Reporter displays ScoreReport (TUI or text)
5. If --fix: Fixer generates missing files from templates
```

## Analysis Categories (each scored 0-20, total 0-100)

### 1. Documentation (0-20)
- README.md exists (4pts) + has badges (1pt) + has install section (1pt) + has usage section (1pt) + >500 chars (1pt)
- LICENSE exists (3pts)
- CONTRIBUTING.md exists (2pts)
- CHANGELOG.md exists (2pts)
- CODE_OF_CONDUCT.md exists (1pt)
- docs/ directory exists (1pt)
- README quality signals: has headings (1pt), has code blocks (1pt), has links (1pt)

### 2. CI/CD (0-20)
- .github/workflows/ has at least one .yml file (6pts)
- Workflow includes test step (3pts)
- Workflow includes lint step (3pts)
- Workflow includes build step (2pts)
- Workflow uses matrix strategy (2pts)
- Workflow triggers on push+PR (2pts)
- Has Makefile or build script (2pts)

### 3. Security (0-20)
- .gitignore exists (3pts) + covers language-specific patterns (2pts)
- No secrets detected in tracked files (5pts) — check for API keys, tokens, passwords
- SECURITY.md exists (3pts)
- Dependabot config exists (2pts)
- No .env files committed (3pts)
- Pre-commit hooks configured (2pts)

### 4. Community (0-20)
- Issue templates exist (4pts)
- PR template exists (3pts)
- Has description set on GitHub (2pts) — only for remote repos
- Has topics/tags (2pts) — only for remote repos
- CODEOWNERS file exists (2pts)
- Has discussions enabled (2pts) — only for remote repos
- Funding/sponsors config (1pt)
- .github/FUNDING.yml exists (1pt)
- Star count > 0 as engagement signal (3pts) — only for remote repos

### 5. Maintenance (0-20)
- Last commit within 30 days (4pts), 90 days (2pts), 365 days (1pt)
- Has releases/tags (3pts)
- No stale branches (merged but not deleted) (3pts)
- Has .editorconfig (2pts)
- go.sum/package-lock.json/poetry.lock committed (2pts) — lock files present
- Reasonable repo size (2pts)
- Has version file or version in code (2pts)
- .gitattributes exists (2pts)

## Language Detection
Scan file extensions to detect languages. Map:
- .go → Go
- .js/.jsx/.ts/.tsx → JavaScript/TypeScript
- .py → Python
- .rs → Rust
- .rb → Ruby
- .java → Java
- .c/.cpp/.h → C/C++
- .php → PHP
- .swift → Swift
- Dockerfile → Docker

Also check for: go.mod, package.json, pyproject.toml, Cargo.toml, Gemfile, pom.xml, composer.json

## Secret Detection Patterns
Check tracked files for:
- `AKIA[0-9A-Z]{16}` — AWS access key
- `[a-zA-Z0-9/+=]{40}` near "secret" — AWS secret key
- `sk-[a-zA-Z0-9]{48}` — OpenAI API key
- `ghp_[a-zA-Z0-9]{36}` — GitHub personal access token
- `-----BEGIN (RSA |EC |DSA )?PRIVATE KEY-----`
- Common patterns: `password\s*=\s*["'][^"']+["']`, `api_key`, `token\s*=`
- .env files in tracked files

## Fixer Templates
Templates are embedded via Go embed. Each template is language-aware:

### LICENSE
- Input: license type (default MIT), author name, year
- Template: Full license text

### .gitignore
- Input: detected languages
- Template: Combined gitignore patterns per language

### CONTRIBUTING.md
- Input: project name, language, build commands
- Template: Standard contributing guide with language-specific setup

### CI Workflow (.github/workflows/ci.yml)
- Input: language, has tests
- Template: Language-specific GitHub Actions workflow

### Issue Templates
- Bug report template
- Feature request template

### SECURITY.md
- Input: project name
- Template: Standard security policy

### CODE_OF_CONDUCT.md
- Contributor Covenant v2.1

### .editorconfig
- Standard config based on detected languages

### PR Template
- Standard PR template

## CLI Commands (Cobra)
- `repokit [path]` — Scan local repo (default: current dir), show TUI report
- `repokit [github-url]` — Scan remote repo via GitHub API
- `repokit --fix` — Auto-fix missing files (interactive selection)
- `repokit --fix --yes` — Auto-fix all without prompting
- `repokit --format json` — Output as JSON (for CI integration)
- `repokit --format text` — Plain text output (no TUI)
- `repokit --category docs,ci` — Only check specific categories
- `repokit compare <path1> <path2>` — Compare two repos
- `repokit version` — Print version

## TUI Design (Bubble Tea + Lip Gloss)
Report card layout:
```
╭──────────────────────────────────────────╮
│  🏥 RepoKit Health Report               │
│  Repository: my-project                  │
│  Overall Score: 73/100  [████████░░] 🟡  │
╰──────────────────────────────────────────╯

 📝 Documentation    16/20  [████████░░] 🟢
    ✅ README.md (with badges, install, usage)
    ✅ LICENSE (MIT)
    ❌ CONTRIBUTING.md — run `repokit --fix`
    ❌ CHANGELOG.md — run `repokit --fix`

 ⚙️  CI/CD            14/20  [███████░░░] 🟡
    ✅ GitHub Actions workflow found
    ✅ Tests detected
    ❌ No lint step
    ❌ No matrix strategy

 🔒 Security          18/20  [█████████░] 🟢
    ✅ .gitignore (Go patterns)
    ✅ No secrets detected
    ❌ No SECURITY.md

 👥 Community         12/20  [██████░░░░] 🟡
    ✅ Issue templates
    ❌ No PR template
    ❌ No CODEOWNERS

 🔧 Maintenance       13/20  [██████░░░░] 🟡
    ✅ Recent commits (2 days ago)
    ✅ Has releases
    ❌ 3 stale branches
    ❌ No .editorconfig
```

Color coding:
- 🟢 Green: 16-20 (80%+)
- 🟡 Yellow: 10-15 (50-79%)
- 🔴 Red: 0-9 (below 50%)

## Dependencies
- cobra — CLI framework
- bubbletea — TUI framework
- lipgloss — TUI styling
- bubbles — TUI components (progress bars, tables)
- go-github — GitHub API client
- yaml.v3 — Config parsing
- embed — Template files (stdlib)

## Testing Strategy
- Unit tests for analyzer (mock file systems)
- Unit tests for scorer (known inputs → expected scores)
- Unit tests for fixer (template output validation)
- Unit tests for secret detection patterns
- Integration tests with real repo fixtures in testdata/
- Table-driven tests throughout
