package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/AbdullahTarakji/repokit/internal/scorer"
)

func sampleReport() *scorer.ScoreReport {
	return &scorer.ScoreReport{
		Overall:  42,
		MaxScore: 100,
		Categories: []scorer.CategoryScore{
			{
				Name: "Documentation", MaxScore: 20, Emoji: "📝", Score: 10,
				Checks: []scorer.CheckResult{
					{Name: "README.md exists", Passed: true, Points: 4, MaxPoints: 4},
					{Name: "LICENSE exists", Passed: false, Points: 0, MaxPoints: 3},
				},
			},
			{Name: "CI/CD", MaxScore: 20, Emoji: "⚙️", Score: 8},
			{Name: "Security", MaxScore: 20, Emoji: "🔒", Score: 12},
			{Name: "Community", MaxScore: 20, Emoji: "👥", Score: 5},
			{Name: "Maintenance", MaxScore: 20, Emoji: "🔧", Score: 7},
		},
	}
}

func TestRenderText(t *testing.T) {
	var buf bytes.Buffer
	RenderText(&buf, sampleReport(), "test-repo")
	output := buf.String()

	checks := []string{
		"test-repo",
		"42/100",
		"Documentation",
		"CI/CD",
		"Security",
		"Community",
		"Maintenance",
		"✅",
		"❌",
	}

	for _, c := range checks {
		if !strings.Contains(output, c) {
			t.Errorf("text output should contain %q", c)
		}
	}
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	err := RenderJSON(&buf, sampleReport(), "test-repo")
	if err != nil {
		t.Fatal(err)
	}

	var jr JSONReport
	if err := json.Unmarshal(buf.Bytes(), &jr); err != nil {
		t.Fatal("invalid JSON:", err)
	}

	if jr.Repository != "test-repo" {
		t.Errorf("expected repo name test-repo, got %s", jr.Repository)
	}
	if jr.Overall != 42 {
		t.Errorf("expected overall 42, got %d", jr.Overall)
	}
	if len(jr.Categories) != 5 {
		t.Errorf("expected 5 categories, got %d", len(jr.Categories))
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		score, max, width int
		wantFilled        int
	}{
		{10, 20, 10, 5},
		{0, 20, 10, 0},
		{20, 20, 10, 10},
		{0, 0, 10, 0},
	}

	for _, tc := range tests {
		bar := progressBar(tc.score, tc.max, tc.width)
		filled := strings.Count(bar, "█")
		if filled != tc.wantFilled {
			t.Errorf("progressBar(%d,%d,%d): filled=%d, want %d, bar=%q",
				tc.score, tc.max, tc.width, filled, tc.wantFilled, bar)
		}
		if len([]rune(bar)) != tc.width {
			t.Errorf("bar width should be %d, got %d", tc.width, len([]rune(bar)))
		}
	}
}

func TestColorIndicator(t *testing.T) {
	tests := []struct {
		score, max int
		want       string
	}{
		{18, 20, "🟢"},
		{12, 20, "🟡"},
		{3, 20, "🔴"},
		{0, 0, "🔴"},
	}

	for _, tc := range tests {
		got := colorIndicator(tc.score, tc.max)
		if got != tc.want {
			t.Errorf("colorIndicator(%d,%d) = %s, want %s", tc.score, tc.max, got, tc.want)
		}
	}
}

func TestScoreEmoji(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{85, "🟢"},
		{60, "🟡"},
		{30, "🔴"},
	}

	for _, tc := range tests {
		got := scoreEmoji(tc.score)
		if got != tc.want {
			t.Errorf("scoreEmoji(%d) = %s, want %s", tc.score, got, tc.want)
		}
	}
}

func TestTUIModel(t *testing.T) {
	report := sampleReport()
	m := NewModel(report, "test-repo")

	view := m.View()
	if !strings.Contains(view, "test-repo") {
		t.Error("TUI view should contain repo name")
	}
	if !strings.Contains(view, "42/100") {
		t.Error("TUI view should contain score")
	}
}
