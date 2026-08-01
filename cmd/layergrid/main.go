package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layergrid/layergrid-cli/internal/baseline"
	"github.com/layergrid/layergrid-cli/internal/report"
	"github.com/layergrid/layergrid-cli/internal/scan"
	"github.com/layergrid/layergrid-cli/internal/trifecta"
	"github.com/layergrid/layergrid-cli/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func rootCmd() *cobra.Command {
	var opts scan.Options
	cmd := &cobra.Command{
		Use:           "layergrid",
		Short:         "Static risk scanner for AI agent stacks",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printBanner(cmd, opts.NoColor)
		},
	}
	cmd.SetHelpFunc(func(helpCmd *cobra.Command, args []string) {
		if helpCmd == cmd {
			_ = printBanner(helpCmd, opts.NoColor)
			return
		}
		_, _ = fmt.Fprint(helpCmd.OutOrStdout(), helpCmd.UsageString())
	})
	scanCmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan an agent stack",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			result, err := scan.Run(path, opts)
			if err != nil {
				return err
			}
			out, err := report.New(opts.Format, opts.NoColor).Format(result)
			if err != nil {
				return err
			}
			if opts.Output != "" {
				if err := os.WriteFile(opts.Output, out, 0o644); err != nil {
					return err
				}
			} else if !opts.Quiet {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(out)); err != nil {
					return err
				}
			}
			if scan.ShouldFail(result.Score, opts.FailOn) {
				os.Exit(1)
			}
			return nil
		},
	}
	scanCmd.Flags().StringVarP(&opts.Format, "format", "f", "human", "human, json, sarif, html, or markdown")
	scanCmd.Flags().StringVarP(&opts.Output, "output", "o", "", "write output to path")
	scanCmd.Flags().StringVar(&opts.FailOn, "fail-on", "never", "critical, high, medium, low, or never")
	scanCmd.Flags().StringVar(&opts.ConfigPath, "config", "", "path to .layergrid.yaml")
	scanCmd.Flags().StringSliceVar(&opts.Include, "include", nil, "include glob patterns")
	scanCmd.Flags().StringSliceVar(&opts.Exclude, "exclude", nil, "exclude glob patterns")
	scanCmd.Flags().StringSliceVar(&opts.Frameworks, "frameworks", nil, "only run selected detectors")
	scanCmd.Flags().StringSliceVar(&opts.Rules, "rules", nil, "only evaluate selected rule IDs or categories")
	scanCmd.Flags().BoolVar(&opts.NoColor, "no-color", false, "disable color")
	scanCmd.Flags().BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose output")
	scanCmd.Flags().BoolVarP(&opts.Quiet, "quiet", "q", false, "suppress stdout")
	cmd.AddCommand(scanCmd)
	cmd.AddCommand(baselineCmd())

	cmd.AddCommand(&cobra.Command{
		Use:   "list-rules",
		Short: "List built-in rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			rules, err := trifecta.LoadBuiltinRules()
			if err != nil {
				return err
			}
			for _, rule := range rules {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%-34s %-8s %s\n", rule.ID, strings.ToUpper(string(rule.Severity)), rule.Name); err != nil {
					return err
				}
			}
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "explain <rule-id>",
		Short: "Explain a built-in rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rules, err := trifecta.LoadBuiltinRules()
			if err != nil {
				return err
			}
			for _, rule := range rules {
				if rule.ID == args[0] {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n\nSeverity: %s\nCategory: %s\n\n%s\n\nFix: %s\n", rule.Name, strings.ToUpper(string(rule.Severity)), rule.Category, rule.Description, rule.Fix); err != nil {
						return err
					}
					if len(rule.References) > 0 {
						if _, err := fmt.Fprintf(cmd.OutOrStdout(), "\nReferences:\n"); err != nil {
							return err
						}
						for _, ref := range rule.References {
							if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", ref); err != nil {
								return err
							}
						}
					}
					return nil
				}
			}
			return fmt.Errorf("unknown rule %s", args[0])
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "layergrid %s\ncommit %s\nbuilt %s\nschema %s\nrubric %s\n", version.Version, version.Commit, version.Date, version.SchemaVersion, version.RubricVersion)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Write a starter .layergrid.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			const body = "version: 1\nexclude:\n  - node_modules/**\n  - .venv/**\n  - testdata/**\nfail_on: high\n"
			if err := os.WriteFile(".layergrid.yaml", []byte(body), 0o644); err != nil {
				return err
			}
			if err := os.MkdirAll(".layergrid", 0o755); err != nil {
				return err
			}
			const note = "# LayerGrid Baseline\n\nCommit `.layergrid/baseline.json` after running `layergrid baseline save` so later scans can compare trust-boundary drift.\n"
			return os.WriteFile(filepath.Join(".layergrid", "README.md"), []byte(note), 0o644)
		},
	})
	return cmd
}

func baselineCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "baseline",
		Short: "Save and compare LayerGrid trust baselines",
	}
	var saveOpts scan.Options
	var saveOutput string
	save := &cobra.Command{
		Use:   "save [path]",
		Short: "Save a LayerGrid baseline",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			result, err := scan.Run(path, saveOpts)
			if err != nil {
				return err
			}
			if saveOutput == "" {
				saveOutput = filepath.Join(".layergrid", "baseline.json")
			}
			if err := os.MkdirAll(filepath.Dir(saveOutput), 0o755); err != nil {
				return err
			}
			b := baseline.FromStack(result.Stack, result.StartedAt)
			if err := baseline.Save(saveOutput, b); err != nil {
				return err
			}
			if !saveOpts.Quiet {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Saved baseline to %s\n", saveOutput)
				return err
			}
			return nil
		},
	}
	save.Flags().StringVarP(&saveOutput, "output", "o", "", "write baseline to path")
	save.Flags().StringVar(&saveOpts.ConfigPath, "config", "", "path to .layergrid.yaml")
	save.Flags().BoolVarP(&saveOpts.Quiet, "quiet", "q", false, "suppress stdout")
	root.AddCommand(save)

	var compareOpts scan.Options
	var baselinePath string
	var compareFormat string
	var failOn string
	compare := &cobra.Command{
		Use:   "compare [path]",
		Short: "Compare current scan to a saved LayerGrid baseline",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			if baselinePath == "" {
				baselinePath = filepath.Join(".layergrid", "baseline.json")
			}
			before, err := baseline.Load(baselinePath)
			if err != nil {
				return err
			}
			result, err := scan.Run(path, compareOpts)
			if err != nil {
				return err
			}
			after := baseline.FromStack(result.Stack, result.StartedAt)
			diff := baseline.Compare(before, after)
			if !compareOpts.Quiet {
				switch compareFormat {
				case "json":
					data, err := json.MarshalIndent(diff, "", "  ")
					if err != nil {
						return err
					}
					if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(data)); err != nil {
						return err
					}
				case "human", "":
					if _, err := fmt.Fprint(cmd.OutOrStdout(), baseline.FormatHuman(diff)); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unsupported baseline format %q", compareFormat)
				}
			}
			if baseline.ShouldFail(diff, failOn) {
				os.Exit(1)
			}
			return nil
		},
	}
	compare.Flags().StringVar(&baselinePath, "baseline", "", "path to baseline JSON")
	compare.Flags().StringVar(&failOn, "fail-on", "any", "scope-widening, tool-added, descriptor-drift, any, or never")
	compare.Flags().StringVarP(&compareFormat, "format", "f", "human", "human or json")
	compare.Flags().StringVar(&compareOpts.ConfigPath, "config", "", "path to .layergrid.yaml")
	compare.Flags().BoolVarP(&compareOpts.Quiet, "quiet", "q", false, "suppress stdout")
	root.AddCommand(compare)

	return root
}

func printBanner(cmd *cobra.Command, noColor bool) error {
	rules, err := trifecta.LoadBuiltinRules()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(cmd.OutOrStdout(), report.Banner(report.BannerOptions{NoColor: noColor, RuleCount: len(rules)}))
	return err
}
