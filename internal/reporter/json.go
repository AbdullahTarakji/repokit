package reporter

import (
	"encoding/json"
	"io"

	"github.com/AbdullahTarakji/repokit/internal/scorer"
)

// JSONReport is the JSON output structure.
type JSONReport struct {
	Repository string              `json:"repository"`
	Overall    int                 `json:"overall"`
	MaxScore   int                 `json:"max_score"`
	Categories []scorer.CategoryScore `json:"categories"`
}

// RenderJSON writes a JSON report to the writer.
func RenderJSON(w io.Writer, report *scorer.ScoreReport, repoName string) error {
	jr := JSONReport{
		Repository: repoName,
		Overall:    report.Overall,
		MaxScore:   report.MaxScore,
		Categories: report.Categories,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jr)
}
