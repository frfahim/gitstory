package analyzer

import (
	"fmt"
	"strings"

	"github.com/frfahim/gitstory/internal/git"
)

// SummarizeCommits creates a naive summary from commit messages
func SummarizeCommits(commits []git.CommitSummary) string {
	if len(commits) == 0 {
		return "No commits to summarize."
	}
	var lines []string
	for _, c := range commits {
		firstLine := strings.SplitN(c.Message, "\n", 2)[0]
		lines = append(lines, fmt.Sprintf("- %s: %s", c.Hash, firstLine))
	}
	return strings.Join(lines, "\n")
}
