package cmd

import (
	"fmt"
	"os"

	"github.com/frfahim/gitstory/internal/analyzer"
	"github.com/frfahim/gitstory/internal/git"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze the current Git repository",
	Long:  `Perform a detailed analysis of the Git repository, including commit statistics and author contributions.`,
	Run: func(cmd *cobra.Command, args []string) {
		var commits []*object.Commit
		currentDir, _ := os.Getwd()
		repo, err := git.OpenRepository(currentDir)
		if err != nil {
			fmt.Println("❌ Not a Git repository.")
			return
		}
		num, _ := cmd.Flags().GetInt("number")
		unique, _ := cmd.Flags().GetBool("unique")
		base, _ := cmd.Flags().GetString("base")
		if num < 1 {
			num = 5
		}
		// If unique mode, try to auto-detect base if not explicitly set
		if unique && (base == "" || base == "auto") {
			autoBase, err := repo.DetectDefaultBranch()
			if err != nil {
				fmt.Printf("⚠️  Could not auto-detect default branch: %v\n", err)
				return
			}
			base = autoBase
		}

		if unique {
			commits, err = repo.ListUniqueCommits(base, num)
			if err != nil {
				fmt.Printf("❌ Error listing unique commits (base=%s): %v\n", base, err)
				return
			}
			fmt.Printf("🔎 Showing last %d commits unique to branch '%s' (vs base '%s'):\n\n", len(commits), repo.CurrentBranchName(), base)
		} else {
			commits, err = repo.ListCommits(num)
			if err != nil {
				fmt.Printf("❌ Error listing commits: %v\n", err)
				return
			}
			fmt.Printf("🔎 Showing last %d commits on branch '%s':\n\n", len(commits), repo.CurrentBranchName())
		}

		summarize, err := repo.ListCommitSummarize(commits)
		if err != nil {
			fmt.Printf("❌ Error listing commit summarizes: %v\n", err)
			return
		}

		fmt.Printf("Commit summarizes (showing last %d):\n\n", len(summarize))
		summary_analyzer := analyzer.SummarizeCommits(summarize)
		fmt.Println(summary_analyzer)
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().IntP("number", "n", 5, "Number of commits to analyze")
	analyzeCmd.Flags().Bool("unique", false, "Show only commits unique to this branch (compared to main)")
	analyzeCmd.Flags().String("base", "main", "Base branch name for unique commit comparison")
}
