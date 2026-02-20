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
	flagDryRun bool
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

	scanCmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a repository and display health report",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAudit,
	}
	scanCmd.Flags().StringVar(&flagFormat, "format", "", "Output format: text, json (default: TUI)")

	fixCmd := &cobra.Command{
		Use:   "fix [path]",
		Short: "Auto-fix missing repository files",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runFix,
	}
	fixCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Show what would be fixed without making changes")
	fixCmd.Flags().BoolVar(&flagYes, "yes", false, "Skip confirmation")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("repokit %s\n", version)
		},
	}

	rootCmd.AddCommand(scanCmd, fixCmd, versionCmd)

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
		return doFix(result, false)
	}

	return nil
}

func runFix(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	result, err := analyzer.Analyze(path)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Show report first
	reporter.RenderText(os.Stdout, scorer.Score(result), result.RepoName)

	return doFix(result, flagDryRun)
}

func doFix(result *analyzer.AnalysisResult, dryRun bool) error {
	if dryRun {
		fmt.Println("\n🔧 Dry run — would fix the following:")
		fr, err := fixer.DryRun(result)
		if err != nil {
			return fmt.Errorf("dry run failed: %w", err)
		}
		if len(fr) == 0 {
			fmt.Println("  Nothing to fix!")
		} else {
			for _, f := range fr {
				fmt.Printf("  📄 Would create %s\n", f)
			}
		}
		return nil
	}

	fmt.Println("\n🔧 Fixing missing files...")
	fixResult, err := fixer.Fix(result, result.RepoPath)
	if err != nil {
		return fmt.Errorf("fix failed: %w", err)
	}
	if len(fixResult.FilesCreated) == 0 {
		fmt.Println("  Nothing to fix!")
	} else {
		for _, f := range fixResult.FilesCreated {
			fmt.Printf("  ✅ Created %s\n", f)
		}
	}
	for _, e := range fixResult.Errors {
		fmt.Fprintf(os.Stderr, "  ⚠️  %v\n", e)
	}
	return nil
}
