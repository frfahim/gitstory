package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/frfahim/gitstory/internal/git"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent commits in the current Git repository",
	Long:  `Show the recent Git commits with hash, author, date, and message.`,
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
		commits, err := repo.ListCommits(num)
		if err != nil {
			fmt.Printf("❌ Error listing commits: %v\n", err)
			return
		}
		fmt.Printf("Showing last %d commits:\n\n", len(commits))
		for _, c := range commits {
			fmt.Printf("• %s | %s | %s\n  %s\n\n",
				c.Hash.String()[:7],
				c.Author.Name,
				c.Author.When.Format(time.RFC822),
				c.Message)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntP("number", "n", 5, "Number of commits to show")
}
