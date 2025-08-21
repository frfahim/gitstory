package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitstory",
	Short: "Turn your commits into stories worth sharing",
	Long: `GitStory analyzes your Git commits and generates intelligent summarize 
that you can share on social media, blogs, or documentation.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
