package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/AbdullahTarakji/repokit/internal/scorer"
	tea "github.com/charmbracelet/bubbletea"
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

func TestTUIInit(t *testing.T) {
	m := NewModel(sampleReport(), "test")
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestTUIUpdateQuit(t *testing.T) {
	m := NewModel(sampleReport(), "test")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Error("pressing q should return quit cmd")
	}
	model := updated.(Model)
	// After quitting, View should be empty
	if model.View() != "" {
		t.Error("quitting view should be empty")
	}
}

func TestTUIUpdateEsc(t *testing.T) {
	m := NewModel(sampleReport(), "test")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Error("pressing esc should return quit cmd")
	}
}

func TestTUIUpdateEnter(t *testing.T) {
	m := NewModel(sampleReport(), "test")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("pressing enter should return quit cmd")
	}
}

func TestTUIUpdateCtrlC(t *testing.T) {
	m := NewModel(sampleReport(), "test")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("pressing ctrl+c should return quit cmd")
	}
}

func TestTUIUpdateOtherKey(t *testing.T) {
	m := NewModel(sampleReport(), "test")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if cmd != nil {
		t.Error("pressing other key should not quit")
	}
}

func TestGetColorStyleAllBranches(t *testing.T) {
	tests := []struct {
		score, max int
	}{
		{18, 20}, // green
		{12, 20}, // yellow
		{3, 20},  // red
		{0, 0},   // zero max
	}
	for _, tc := range tests {
		// Just ensure no panic
		_ = getColorStyle(tc.score, tc.max)
	}
}

func TestProgressBarOverflow(t *testing.T) {
	// score > max should cap at width
	bar := progressBar(25, 20, 10)
	filled := strings.Count(bar, "█")
	if filled != 10 {
		t.Errorf("overflow should cap at width, got %d filled", filled)
	}
}

func TestRenderTextWithDescription(t *testing.T) {
	report := &scorer.ScoreReport{
		Overall:  50,
		MaxScore: 100,
		Categories: []scorer.CategoryScore{
			{
				Name: "Test", MaxScore: 20, Emoji: "🔧", Score: 10,
				Checks: []scorer.CheckResult{
					{Name: "Check A", Passed: true, Points: 5, MaxPoints: 5, Description: "Custom desc"},
					{Name: "Check B", Passed: false, Points: 0, MaxPoints: 5},
				},
			},
		},
	}
	var buf bytes.Buffer
	RenderText(&buf, report, "desc-repo")
	output := buf.String()
	if !strings.Contains(output, "Custom desc") {
		t.Error("should show description when present")
	}
	if !strings.Contains(output, "Check B") {
		t.Error("should show name when no description")
	}
}

func TestTUIViewWithDescription(t *testing.T) {
	report := &scorer.ScoreReport{
		Overall:  50,
		MaxScore: 100,
		Categories: []scorer.CategoryScore{
			{
				Name: "Cat", MaxScore: 20, Emoji: "📝", Score: 10,
				Checks: []scorer.CheckResult{
					{Name: "A", Passed: true, Points: 5, MaxPoints: 5, Description: "Has desc"},
					{Name: "B", Passed: false, Points: 0, MaxPoints: 5},
				},
			},
		},
	}
	m := NewModel(report, "test")
	view := m.View()
	if !strings.Contains(view, "Has desc") {
		t.Error("TUI should show description")
	}
}
