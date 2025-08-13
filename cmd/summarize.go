package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize recent commits with advanced (AI-powered) summaries",
	Long: `Summarize recent or unique commits using LLMs (OpenAI, Claude, Gemini, etc).
Supports different styles (twitter, blog, linkedin, etc) and user-provided context.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Wire up provider selection, style, and context
		fmt.Println("Summarize functionality coming soon!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(summarizeCmd)
	summarizeCmd.Flags().String("provider", "", "LLM provider to use (openai, claude, gemini, etc)")
	summarizeCmd.Flags().String("style", "blog", "Summary style (blog, twitter, linkedin, note, etc)")
	summarizeCmd.Flags().String("context", "", "Extra context for the LLM to improve the summary")
	summarizeCmd.Flags().String("numbers", "", "Numbers of latest commits to summarize (e.g. 2 or 3)")
	summarizeCmd.Flags().String("commits", "", "Commit range to summarize (e.g. 123abc..456def)")
	summarizeCmd.Flags().Bool("unique", false, "Summarize only unique commits (PR-style)")
	summarizeCmd.Flags().String("base", "main", "Base branch for unique commit calculation")
	summarizeCmd.Flags().String("output", "", "Write summary to file (optional)")
}
