package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	BuildTime = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GitStory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GitStory %s\n", Version)
		fmt.Printf("Built at: %s\n", BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
