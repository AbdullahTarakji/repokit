package scorer

import (
	"testing"
	"time"

	"github.com/AbdullahTarakji/repokit/internal/analyzer"
)

func TestScore_Empty(t *testing.T) {
	result := &analyzer.AnalysisResult{}
	report := Score(result)

	if report.MaxScore != 100 {
		t.Errorf("max score should be 100, got %d", report.MaxScore)
	}
	if len(report.Categories) != 5 {
		t.Errorf("expected 5 categories, got %d", len(report.Categories))
	}
}

func TestScore_PerfectDocs(t *testing.T) {
	result := &analyzer.AnalysisResult{
		HasREADME:           true,
		READMELength:        1000,
		READMEHasBadges:     true,
		READMEHasInstall:    true,
		READMEHasUsage:      true,
		READMEHasHeadings:   true,
		READMEHasCodeBlocks: true,
		READMEHasLinks:      true,
		HasLicense:          true,
		HasContributing:     true,
		HasChangelog:        true,
		HasCodeOfConduct:    true,
		HasDocsDir:          true,
	}

	report := Score(result)
	docsCat := report.Categories[0]
	if docsCat.Score != 20 {
		t.Errorf("perfect docs should score 20, got %d", docsCat.Score)
	}
}

func TestScore_PerfectCI(t *testing.T) {
	result := &analyzer.AnalysisResult{
		CIFiles:         []string{"ci.yml"},
		CIHasTest:       true,
		CIHasLint:       true,
		CIHasBuild:      true,
		CIHasMatrix:     true,
		CITriggerPushPR: true,
		HasMakefile:     true,
	}

	report := Score(result)
	ciCat := report.Categories[1]
	if ciCat.Score != 20 {
		t.Errorf("perfect CI should score 20, got %d", ciCat.Score)
	}
}

func TestScore_PerfectSecurity(t *testing.T) {
	result := &analyzer.AnalysisResult{
		HasGitignore:        true,
		GitignoreCoversLang: true,
		SecretsFound:        nil,
		HasSecurityMD:       true,
		HasDependabot:       true,
		HasEnvFiles:         false,
		HasPreCommit:        true,
	}

	report := Score(result)
	secCat := report.Categories[2]
	if secCat.Score != 20 {
		t.Errorf("perfect security should score 20, got %d", secCat.Score)
	}
}

func TestScore_MaintenanceRecency(t *testing.T) {
	tests := []struct {
		name       string
		lastCommit time.Time
		wantMin    int
	}{
		{"recent", time.Now().Add(-24 * time.Hour), 4},
		{"30 days", time.Now().Add(-60 * 24 * time.Hour), 2},
		{"6 months", time.Now().Add(-200 * 24 * time.Hour), 1},
		{"stale", time.Now().Add(-400 * 24 * time.Hour), 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := &analyzer.AnalysisResult{
				LastCommitDate: tc.lastCommit,
			}
			report := Score(result)
			maintCat := report.Categories[4]
			commitCheck := maintCat.Checks[0]
			if commitCheck.Points < tc.wantMin {
				t.Errorf("expected at least %d commit points, got %d", tc.wantMin, commitCheck.Points)
			}
		})
	}
}

func TestScore_OverallSum(t *testing.T) {
	result := &analyzer.AnalysisResult{
		HasREADME:    true,
		HasLicense:   true,
		HasGitignore: true,
	}

	report := Score(result)
	sum := 0
	for _, cat := range report.Categories {
		sum += cat.Score
	}
	if report.Overall != sum {
		t.Errorf("overall %d != sum of categories %d", report.Overall, sum)
	}
}

func TestScore_CategoryMaxScores(t *testing.T) {
	report := Score(&analyzer.AnalysisResult{})
	for _, cat := range report.Categories {
		if cat.MaxScore != 20 {
			t.Errorf("category %s max should be 20, got %d", cat.Name, cat.MaxScore)
		}
	}
}

func TestScore_SecretsReduceScore(t *testing.T) {
	clean := Score(&analyzer.AnalysisResult{})
	dirty := Score(&analyzer.AnalysisResult{
		SecretsFound: []analyzer.SecretFinding{{File: "test", Type: "AWS"}},
	})

	cleanSec := clean.Categories[2].Score
	dirtySec := dirty.Categories[2].Score
	if dirtySec >= cleanSec {
		t.Error("secrets should reduce security score")
	}
}

func TestAddCheck(t *testing.T) {
	cat := &CategoryScore{Name: "Test", MaxScore: 20}
	cat.addCheck("check1", true, 5)
	cat.addCheck("check2", false, 3)

	if len(cat.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(cat.Checks))
	}
	if cat.Score != 5 {
		t.Errorf("expected score 5, got %d", cat.Score)
	}
	if cat.Checks[0].Points != 5 {
		t.Error("passed check should have full points")
	}
	if cat.Checks[1].Points != 0 {
		t.Error("failed check should have 0 points")
	}
}
