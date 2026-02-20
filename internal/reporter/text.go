package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/AbdullahTarakji/repokit/internal/scorer"
)

// RenderText writes a plain text report to the writer.
func RenderText(w io.Writer, report *scorer.ScoreReport, repoName string) {
	fmt.Fprintf(w, "\n🏥 RepoKit Health Report\n")
	fmt.Fprintf(w, "Repository: %s\n", repoName)
	fmt.Fprintf(w, "Overall Score: %d/%d %s\n\n", report.Overall, report.MaxScore, scoreEmoji(report.Overall))

	for _, cat := range report.Categories {
		bar := progressBar(cat.Score, cat.MaxScore, 10)
		color := colorIndicator(cat.Score, cat.MaxScore)
		fmt.Fprintf(w, "%s %-16s %2d/%d  [%s] %s\n", cat.Emoji, cat.Name, cat.Score, cat.MaxScore, bar, color)

		for _, check := range cat.Checks {
			icon := "❌"
			if check.Passed {
				icon = "✅"
			}
			desc := check.Name
			if check.Description != "" {
				desc = check.Description
			}
			fmt.Fprintf(w, "   %s %s\n", icon, desc)
		}
		fmt.Fprintln(w)
	}
}

func progressBar(score, max, width int) string {
	if max == 0 {
		return strings.Repeat("░", width)
	}
	filled := (score * width) / max
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func colorIndicator(score, max int) string {
	pct := 0
	if max > 0 {
		pct = (score * 100) / max
	}
	if pct >= 80 {
		return "🟢"
	}
	if pct >= 50 {
		return "🟡"
	}
	return "🔴"
}

func scoreEmoji(score int) string {
	if score >= 80 {
		return "🟢"
	}
	if score >= 50 {
		return "🟡"
	}
	return "🔴"
}
