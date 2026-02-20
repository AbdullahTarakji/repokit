# BACKLOG.md — RepoKit

## Phase 1: Core (P0) — v0.1.0
- [ ] Project scaffolding (go.mod, cobra CLI, directory structure)
- [ ] Analyzer: scan local repo for files, detect languages
- [ ] Scorer: compute per-category and overall scores (0-100)
- [ ] Reporter: TUI report card with Bubble Tea + Lip Gloss
- [ ] Reporter: plain text output (--format text)
- [ ] Reporter: JSON output (--format json)
- [ ] Fixer: auto-generate missing files from templates
- [ ] Templates: LICENSE, .gitignore, CONTRIBUTING.md, CI workflow, issue templates, SECURITY.md, CODE_OF_CONDUCT.md, .editorconfig, PR template, CHANGELOG.md
- [ ] Secret detection in tracked files
- [ ] Language detection from file extensions + config files
- [ ] Unit tests for all packages (40+ tests)
- [ ] CI: GitHub Actions (lint, test, build on macOS + Ubuntu)
- [ ] README.md with badges, installation, usage, screenshots
- [ ] CHANGELOG.md, CONTRIBUTING.md

## Phase 2: Enhanced (P1) — v0.2.0
- [ ] GitHub API: scan remote repos by URL
- [ ] Compare command: side-by-side repo comparison
- [ ] Config file: .repokit.yaml for custom rules
- [ ] Category filtering: --category docs,ci
- [ ] Fixer: interactive selection mode (checkboxes)

## Phase 3: Polish (P2) — v0.3.0
- [ ] GoReleaser for binary releases
- [ ] Homebrew formula
- [ ] More language templates (Rust, Ruby, Java, PHP)
- [ ] CI integration guide (run in GitHub Actions)
- [ ] Badge generation ("RepoKit Score: 85/100")
