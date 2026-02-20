package main

import (
	"fmt"
	"os"

	"github.com/AbdullahTarakji/repokit/internal/analyzer"
	"github.com/AbdullahTarakji/repokit/internal/fixer"
	"github.com/AbdullahTarakji/repokit/internal/reporter"
	"github.com/AbdullahTarakji/repokit/internal/scorer"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

var (
	flagFix    bool
	flagYes    bool
	flagFormat string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "repokit [path]",
		Short: "GitHub repository health auditor and fixer",
		Long:  "RepoKit audits GitHub repositories for health and quality, assigns a score (0-100), and can auto-fix missing files.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAudit,
	}

	rootCmd.Flags().BoolVar(&flagFix, "fix", false, "Auto-fix missing files")
	rootCmd.Flags().BoolVar(&flagYes, "yes", false, "Skip confirmation when fixing")
	rootCmd.Flags().StringVar(&flagFormat, "format", "", "Output format: text, json (default: TUI)")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("repokit %s\n", version)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runAudit(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	result, err := analyzer.Analyze(path)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	report := scorer.Score(result)

	switch flagFormat {
	case "json":
		return reporter.RenderJSON(os.Stdout, report, result.RepoName)
	case "text":
		reporter.RenderText(os.Stdout, report, result.RepoName)
	default:
		if err := reporter.RunTUI(report, result.RepoName); err != nil {
			// Fall back to text if TUI fails
			reporter.RenderText(os.Stdout, report, result.RepoName)
		}
	}

	if flagFix {
		fmt.Println("\n🔧 Fixing missing files...")
		fixResult, err := fixer.Fix(result, result.RepoPath)
		if err != nil {
			return fmt.Errorf("fix failed: %w", err)
		}
		if len(fixResult.FilesCreated) == 0 {
			fmt.Println("Nothing to fix!")
		} else {
			for _, f := range fixResult.FilesCreated {
				fmt.Printf("  ✅ Created %s\n", f)
			}
		}
		for _, e := range fixResult.Errors {
			fmt.Fprintf(os.Stderr, "  ⚠️  %v\n", e)
		}
	}

	return nil
}
