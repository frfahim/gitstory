package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/frfahim/gitstory/internal/git"
	"github.com/frfahim/gitstory/internal/llm"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Generate AI-powered summarize of your commits details",
	Long: `Generate intelligent summarize of your git commits details using AI providers like OpenAI and Gemini.
Supports different platforms (twitter/X, blog, linkedin, technical, notes) with optimized prompts.

Examples:
  gitstory summarize                                    # Interactive mode (coming soon)
  gitstory summarize --platform blog                    # Auto-detect provider
  gitstory summarize --provider gemini --platform twitter/X
  gitstory summarize --platform technical --commits 10 --context "Sprint 23"`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return runSummarize(cmd, args)
	},
}

func runSummarize(cmd *cobra.Command, args []string) error {
	// Get flags
	provider, _ := cmd.Flags().GetString("provider")
	platform, _ := cmd.Flags().GetString("platform")
	userContext, _ := cmd.Flags().GetString("context")
	numbers, _ := cmd.Flags().GetString("numbers")
	unique, _ := cmd.Flags().GetBool("unique")
	base, _ := cmd.Flags().GetString("base")
	output, _ := cmd.Flags().GetString("output")

	// Validate platform
	if platform == "" {
		platform = "technical" // default
	}
	if err := llm.ValidatePlatform(platform); err != nil {
		return fmt.Errorf("invalid platform: %w", err)
	}

	// Normalize platform (e.g., convert "X" to "twitter")
	normalizedPlatform := llm.NormalizePlatform(platform)

	// Get current directory and open repository
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := git.OpenRepository(currentDir)
	if err != nil {
		return fmt.Errorf("❌ Not a Git repository: %w", err)
	}

	// Determine number of commits
	numCommits := 5 // default
	if numbers != "" {
		if n, err := strconv.Atoi(numbers); err == nil && n > 0 {
			numCommits = n
		}
	}

	// Get commits based on options
	var commits []*object.Commit
	if unique {
		fmt.Printf("🔍 Getting unique commits from current branch compared to %s...\n", base)
		commits, err = repo.ListUniqueCommits(base, numCommits)
	} else {
		fmt.Printf("🔍 Getting last %d commits...\n", numCommits)
		commits, err = repo.ListCommits(numCommits)
	}

	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println("ℹ️ No commits found to summarize.")
		return nil
	}

	// Convert to commit summarize then to LLM format
	summarizeCommitList, err := repo.ListCommitSummarize(commits)
	if err != nil {
		return fmt.Errorf("failed to get commit summarizes: %w", err)
	}

	fmt.Printf("📝 Found %d commit(s) to summarize\n", len(summarizeCommitList))

	// Smart provider selection
	if provider == "" {
		available := llm.DetectAvailableProviders()
		if len(available) == 0 {
			return fmt.Errorf("❌ No LLM providers configured. Please set OPENAI_API_KEY or GEMINI_API_KEY")
		}
		if len(available) == 1 {
			provider = string(available[0])
			fmt.Printf("🤖 Using %s (auto-detected)\n", provider)
		} else {
			fmt.Printf("Multiple providers available: %v\n", available)
			fmt.Println("Use --provider to specify or run with --interactive (coming soon)")
			provider = string(available[0]) // Use first available
			fmt.Printf("🤖 Using %s (first available)\n", provider)
		}
	}

	// Validate provider
	if err := llm.ValidateProvider(provider); err != nil {
		return fmt.Errorf("invalid provider: %w", err)
	}

	// Create LLM client
	fmt.Printf("🔧 Creating %s client...\n", provider)
	client, err := llm.NewClient(llm.ClientConfig{
		Provider: llm.Provider(provider),
	})
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Validate credentials (skip for now since it's commented out in interface)
	fmt.Printf("🔑 Using %s provider...\n", provider)
	// TODO: Add credential validation when interface is updated
	// if err := client.ValidateCredentials(ctx); err != nil {
	//     return fmt.Errorf("credential validation failed: %w", err)
	// }

	// Create summary request
	request := &llm.SummaryRequest{
		Commits:     summarizeCommitList,
		Platform:    normalizedPlatform,
		UserContext: userContext,
	}

	// Generate summary
	ctx := context.Background()
	fmt.Printf("🧠 Generating %s summary using %s...\n", normalizedPlatform, provider)
	response, err := client.Summarize(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	// Display result
	displaySummary(response, normalizedPlatform)

	// Save to file if requested
	if output != "" {
		if err := saveSummaryToFile(response, output); err != nil {
			fmt.Printf("⚠️ Failed to save to file: %v\n", err)
		} else {
			fmt.Printf("💾 Summary saved to %s\n", output)
		}
	}

	return nil
}

func displaySummary(response *llm.SummaryResponse, platform llm.Platform) {
	icons := map[llm.Platform]string{
		"blog":      "📝",
		"twitter/X": "🐦",
		"linkedin":  "💼",
		"technical": "🔧",
		"notes":     "📋",
	}

	icon := icons[platform]
	if icon == "" {
		icon = "✨"
	}

	// Use cases.Title instead of deprecated strings.Title
	caser := cases.Title(language.English)
	platformTitle := caser.String(string(platform))

	fmt.Printf("\n%s %s Summary:\n", icon, platformTitle)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(response.Summary)

	// Show statistics using helper methods
	fmt.Printf("\n📊 %s\n", response.GetStats())

	if !response.MeetsRequirements() {
		fmt.Printf("⚠️ Warning: Summary may not meet %s platform requirements\n", platform)
	}
}

func saveSummaryToFile(response *llm.SummaryResponse, filename string) error {
	// Use cases.Title for consistent title casing
	caser := cases.Title(language.English)
	platformTitle := caser.String(string(response.Platform))

	content := fmt.Sprintf("# %s Summary\n\n%s\n\n---\nGenerated by [GitStory](https://github.com/frfahim/gitstory)\nPlatform: %s\nStats: %s\n",
		platformTitle,
		response.Summary,
		response.Platform,
		response.GetStats())

	return os.WriteFile(filename, []byte(content), 0644)
}

func init() {
	rootCmd.AddCommand(summarizeCmd)

	// Provider and platform options
	summarizeCmd.Flags().String("provider", "", "LLM provider (openai, gemini, claude)")
	summarizeCmd.Flags().String("platform", "", "Target platform (twitter/X, linkedin, blog, technical, notes)")

	// Commit selection options
	summarizeCmd.Flags().String("numbers", "", "Number of latest commits to summarize (e.g. 5)")
	summarizeCmd.Flags().Bool("unique", false, "Summarize only commits unique to current branch")
	summarizeCmd.Flags().String("base", "main", "Base branch for unique commit comparison")

	// Content options
	summarizeCmd.Flags().String("context", "", "Additional context to improve the summary")

	// Output options
	summarizeCmd.Flags().String("output", "", "Save summary to file (optional)")

	// Shell completion
	summarizeCmd.RegisterFlagCompletionFunc("platform", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"twitter", "X", "linkedin", "blog", "technical", "notes"}, cobra.ShellCompDirectiveNoFileComp
	})

	summarizeCmd.RegisterFlagCompletionFunc("provider", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"openai", "gemini", "claude"}, cobra.ShellCompDirectiveNoFileComp
	})
}
