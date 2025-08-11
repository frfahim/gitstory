package cmd

import (
	"fmt"
	"os"

	"github.com/frfahim/gitstory/internal/git"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Git repository status and information",
	Long: `Display information about the current Git repository including:
- Repository path
- Current branch
- Remote URL
- Basic commit count
- Git repository validity`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("❌ Error getting current directory: %v\n", err)
			return
		}

		fmt.Printf("🔍 Checking Git repository status...\n\n")

		// Check if we're in a Git repository
		repo, err := git.OpenRepository(currentDir)
		if err != nil {
			fmt.Printf("❌ Not a Git repository\n")
			fmt.Printf("   Path: %s\n", currentDir)
			fmt.Printf("   Error: %v\n", err)
			fmt.Println("\n💡 Navigate to a Git repository directory or initialize one with 'git init'")
			return
		}

		// Get repository information
		info, err := repo.GetInfo()
		if err != nil {
			fmt.Printf("⚠️  Git repository found but error getting info: %v\n", err)
			return
		}

		// Display the information
		fmt.Printf("✅ Git Repository Found\n")
		fmt.Printf("   📁 Path: %s\n", info.Path)
		fmt.Printf("   🌿 Branch: %s\n", info.CurrentBranch)

		if info.RemoteURL != "" {
			fmt.Printf("   🔗 Remote: %s\n", info.RemoteURL)
		} else {
			fmt.Printf("   🔗 Remote: (no remote configured)\n")
		}

		fmt.Printf("   📊 Commits: %d\n", info.CommitCount)

		fmt.Println("\n🎉 Ready to analyze commits with GitStory!")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
