package scorer

import (
	"fmt"
	"time"

	"github.com/AbdullahTarakji/repokit/internal/analyzer"
)

// CheckResult represents a single check's pass/fail status.
type CheckResult struct {
	Name        string `json:"name"`
	Passed      bool   `json:"passed"`
	Points      int    `json:"points"`
	MaxPoints   int    `json:"max_points"`
	Description string `json:"description"`
}

// CategoryScore represents a category's score.
type CategoryScore struct {
	Name      string        `json:"name"`
	Score     int           `json:"score"`
	MaxScore  int           `json:"max_score"`
	Emoji     string        `json:"emoji"`
	Checks    []CheckResult `json:"checks"`
}

// ScoreReport is the full scoring output.
type ScoreReport struct {
	Overall      int             `json:"overall"`
	MaxScore     int             `json:"max_score"`
	Categories   []CategoryScore `json:"categories"`
}

// Score computes the score report from analysis results.
func Score(result *analyzer.AnalysisResult) *ScoreReport {
	report := &ScoreReport{MaxScore: 100}

	cats := []CategoryScore{
		scoreDocumentation(result),
		scoreCI(result),
		scoreSecurity(result),
		scoreCommunity(result),
		scoreMaintenance(result),
	}

	for _, c := range cats {
		report.Overall += c.Score
	}
	report.Categories = cats

	return report
}

func scoreDocumentation(r *analyzer.AnalysisResult) CategoryScore {
	cat := CategoryScore{Name: "Documentation", MaxScore: 20, Emoji: "📝"}

	cat.addCheck("README.md exists", r.HasREADME, 4)
	cat.addCheck("README has badges", r.READMEHasBadges, 1)
	cat.addCheck("README has install section", r.READMEHasInstall, 1)
	cat.addCheck("README has usage section", r.READMEHasUsage, 1)
	cat.addCheck("README > 500 chars", r.READMELength > 500, 1)
	cat.addCheck("LICENSE exists", r.HasLicense, 3)
	cat.addCheck("CONTRIBUTING.md exists", r.HasContributing, 2)
	cat.addCheck("CHANGELOG.md exists", r.HasChangelog, 2)
	cat.addCheck("CODE_OF_CONDUCT.md exists", r.HasCodeOfConduct, 1)
	cat.addCheck("docs/ directory exists", r.HasDocsDir, 1)
	cat.addCheck("README has headings", r.READMEHasHeadings, 1)
	cat.addCheck("README has code blocks", r.READMEHasCodeBlocks, 1)
	cat.addCheck("README has links", r.READMEHasLinks, 1)

	return cat
}

func scoreCI(r *analyzer.AnalysisResult) CategoryScore {
	cat := CategoryScore{Name: "CI/CD", MaxScore: 20, Emoji: "⚙️"}

	cat.addCheck("CI workflow exists", len(r.CIFiles) > 0, 6)
	cat.addCheck("CI includes test step", r.CIHasTest, 3)
	cat.addCheck("CI includes lint step", r.CIHasLint, 3)
	cat.addCheck("CI includes build step", r.CIHasBuild, 2)
	cat.addCheck("CI uses matrix strategy", r.CIHasMatrix, 2)
	cat.addCheck("CI triggers on push + PR", r.CITriggerPushPR, 2)
	cat.addCheck("Makefile or build script exists", r.HasMakefile, 2)

	return cat
}

func scoreSecurity(r *analyzer.AnalysisResult) CategoryScore {
	cat := CategoryScore{Name: "Security", MaxScore: 20, Emoji: "🔒"}

	cat.addCheck(".gitignore exists", r.HasGitignore, 3)
	cat.addCheck(".gitignore covers project languages", r.GitignoreCoversLang, 2)
	cat.addCheck("No secrets detected", len(r.SecretsFound) == 0, 5)
	cat.addCheck("SECURITY.md exists", r.HasSecurityMD, 3)
	cat.addCheck("Dependabot configured", r.HasDependabot, 2)
	cat.addCheck("No .env files committed", !r.HasEnvFiles, 3)
	cat.addCheck("Pre-commit hooks configured", r.HasPreCommit, 2)

	return cat
}

func scoreCommunity(r *analyzer.AnalysisResult) CategoryScore {
	cat := CategoryScore{Name: "Community", MaxScore: 20, Emoji: "👥"}

	cat.addCheck("Issue templates exist", len(r.IssueTemplates) > 0, 4)
	cat.addCheck("PR template exists", r.HasPRTemplate, 3)
	// Remote-only checks scored 0 for local analysis
	cat.addCheck("Has repo description", false, 2)
	cat.addCheck("Has topics/tags", false, 2)
	cat.addCheck("CODEOWNERS file exists", r.HasCodeowners, 2)
	cat.addCheck("Discussions enabled", false, 2)
	cat.addCheck("FUNDING.yml exists", r.HasFundingYML, 1)
	cat.addCheck("Has sponsors config", false, 1)
	cat.addCheck("Star engagement", false, 3)

	return cat
}

func scoreMaintenance(r *analyzer.AnalysisResult) CategoryScore {
	cat := CategoryScore{Name: "Maintenance", MaxScore: 20, Emoji: "🔧"}

	// Last commit recency
	commitPts := 0
	commitDesc := "No commit info"
	if !r.LastCommitDate.IsZero() {
		days := int(time.Since(r.LastCommitDate).Hours() / 24)
		if days <= 30 {
			commitPts = 4
			commitDesc = fmt.Sprintf("Last commit %d days ago", days)
		} else if days <= 90 {
			commitPts = 2
			commitDesc = fmt.Sprintf("Last commit %d days ago", days)
		} else if days <= 365 {
			commitPts = 1
			commitDesc = fmt.Sprintf("Last commit %d days ago", days)
		} else {
			commitDesc = fmt.Sprintf("Last commit %d days ago (stale)", days)
		}
	}
	cat.Checks = append(cat.Checks, CheckResult{
		Name: "Recent commits", Passed: commitPts > 0,
		Points: commitPts, MaxPoints: 4, Description: commitDesc,
	})
	cat.Score += commitPts

	cat.addCheck("Has releases/tags", r.HasReleases, 3)
	cat.addCheck("No stale branches", r.StaleBranches == 0, 3)
	cat.addCheck("Has .editorconfig", r.HasEditorconfig, 2)
	cat.addCheck("Lock files committed", r.HasLockFile, 2)
	cat.addCheck("Reasonable repo size", r.RepoSizeKB < 100*1024, 2)
	cat.addCheck("Has version file", r.HasVersionFile, 2)
	cat.addCheck("Has .gitattributes", r.HasGitattributes, 2)

	return cat
}

func (c *CategoryScore) addCheck(name string, passed bool, maxPoints int) {
	pts := 0
	if passed {
		pts = maxPoints
	}
	c.Checks = append(c.Checks, CheckResult{
		Name:      name,
		Passed:    passed,
		Points:    pts,
		MaxPoints: maxPoints,
	})
	c.Score += pts
}
