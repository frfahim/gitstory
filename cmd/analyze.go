package cmd

import (
	"fmt"
	"os"

	"github.com/frfahim/gitstory/internal/analyzer"
	"github.com/frfahim/gitstory/internal/git"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze the current Git repository",
	Long:  `Perform a detailed analysis of the Git repository, including commit statistics and author contributions.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, _ := os.Getwd()
		repo, err := git.OpenRepository(currentDir)
		if err != nil {
			fmt.Println("❌ Not a Git repository.")
			return
		}
		num, _ := cmd.Flags().GetInt("number")
		if num < 1 {
			num = 5
		}
		summaries, err := repo.ListCommitSummaries(num)
		if err != nil {
			fmt.Printf("❌ Error listing commit summaries: %v\n", err)
			return
		}

		fmt.Printf("Commit summaries (showing last %d):\n\n", len(summaries))
		summary_analyzer := analyzer.SummarizeCommits(summaries)
		fmt.Println(summary_analyzer)
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().IntP("number", "n", 5, "Number of commits to analyze")
}
