package reporter

import (
	"fmt"
	"strings"

	"github.com/AbdullahTarakji/repokit/internal/scorer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	categoryStyle = lipgloss.NewStyle().Bold(true)

	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			BorderForeground(lipgloss.Color("62"))
)

// Model is the Bubble Tea model for the TUI report.
type Model struct {
	Report   *scorer.ScoreReport
	RepoName string
	quitting bool
}

// NewModel creates a new TUI model.
func NewModel(report *scorer.ScoreReport, repoName string) Model {
	return Model{Report: report, RepoName: repoName}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model and handles key events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc", "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model and renders the health report.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Header
	header := fmt.Sprintf("🏥 RepoKit Health Report\nRepository: %s\nOverall Score: %d/%d  [%s] %s",
		m.RepoName, m.Report.Overall, m.Report.MaxScore,
		progressBar(m.Report.Overall, m.Report.MaxScore, 10),
		scoreEmoji(m.Report.Overall))

	b.WriteString(boxStyle.Render(header))
	b.WriteString("\n\n")

	// Categories
	for _, cat := range m.Report.Categories {
		bar := progressBar(cat.Score, cat.MaxScore, 10)
		color := colorIndicator(cat.Score, cat.MaxScore)
		style := getColorStyle(cat.Score, cat.MaxScore)

		line := fmt.Sprintf("%s %-16s %s  [%s] %s",
			cat.Emoji, cat.Name,
			style.Render(fmt.Sprintf("%2d/%d", cat.Score, cat.MaxScore)),
			bar, color)
		b.WriteString(categoryStyle.Render(line))
		b.WriteString("\n")

		for _, check := range cat.Checks {
			icon := "❌"
			if check.Passed {
				icon = "✅"
			}
			desc := check.Name
			if check.Description != "" {
				desc = check.Description
			}
			fmt.Fprintf(&b, "   %s %s\n", icon, desc)
		}
		b.WriteString("\n")
	}

	b.WriteString("Press q to quit\n")
	return b.String()
}

func getColorStyle(score, max int) lipgloss.Style {
	pct := 0
	if max > 0 {
		pct = (score * 100) / max
	}
	if pct >= 80 {
		return greenStyle
	}
	if pct >= 50 {
		return yellowStyle
	}
	return redStyle
}

// RunTUI runs the Bubble Tea TUI.
func RunTUI(report *scorer.ScoreReport, repoName string) error {
	m := NewModel(report, repoName)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}
