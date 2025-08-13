package llm

import (
	"context"
	"fmt"
	"strings"
)

// Provider represents different AI providers
type Provider string

const (
	OpenAI Provider = "openai"
	Claude Provider = "claude"
	Gemini Provider = "gemini"
)

// Platform represents target platforms for summaries
type Platform string

const (
	LinkedIn  Platform = "linkedin"  // LinkedIn posts
	Twitter   Platform = "twitter"   // Twitter/X posts
	Blog      Platform = "blog"      // Blog posts
	Note      Platform = "note"      // Personal notes
	Technical Platform = "technical" // Technical documentation
)

// Client interface that all AI providers must implement
type Client interface {
	// Summarize generates a summary based on the request
	Summarize(ctx context.Context, request *SummaryRequest) (*SummaryResponse, error)

	// ValidateCredentials checks if API credentials are valid
	// ValidateCredentials(ctx context.Context) error

	// GetProvider returns the provider type
	GetProvider() Provider
}

// ClientConfig contains configuration for AI clients
type ClientConfig struct {
	Provider Provider `json:"provider"`
	APIKey   string   `json:"api_key"`
	Model    string   `json:"model,omitempty"`
}

// CommitData represents basic commit information for AI
type CommitData struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// SummaryRequest contains all information needed for AI summarization
type SummaryRequest struct {
	// Core commit data
	Commits []CommitData `json:"commits"`

	// Summary configuration
	Platform  Platform `json:"platform"`
	MaxLength int64    `json:"max_length,omitempty"`

	// Additional context provided by the user
	UserContext string `json:"user_context,omitempty"`
}

// SummaryResponse contains the AI-generated summary
type SummaryResponse struct {
	Summary  string   `json:"summary"`
	Platform Platform `json:"platform"`
}

// CharCount returns the character count of the summary
func (s *SummaryResponse) CharCount() int {
	return len(s.Summary)
}

// WordCount returns the word count of the summary
func (s *SummaryResponse) WordCount() int {
	return len(strings.Fields(s.Summary))
}

// MeetsRequirements checks if the summary meets platform requirements
func (s *SummaryResponse) MeetsRequirements() bool {
	switch s.Platform {
	case Twitter:
		return s.CharCount() <= 280
	case LinkedIn:
		return s.CharCount() <= 3000
	case Blog:
		return s.WordCount() >= 100 && s.WordCount() <= 800
	case Technical:
		return s.WordCount() >= 50 && s.WordCount() <= 600
	case Note:
		return s.WordCount() <= 400
	default:
		return true
	}
}

// GetPlatformLimits returns the limits for the platform
func (s *SummaryResponse) GetPlatformLimits() (minWords, maxWords, maxChars int) {
	switch s.Platform {
	case Twitter:
		return 0, 40, 280
	case LinkedIn:
		return 20, 400, 3000
	case Blog:
		return 100, 800, 5000
	case Technical:
		return 50, 600, 4000
	case Note:
		return 0, 400, 2000
	default:
		return 0, 1000, 5000
	}
}

// GetStats returns formatted statistics
func (s *SummaryResponse) GetStats() string {
	chars := s.CharCount()
	words := s.WordCount()
	meetsReq := s.MeetsRequirements()

	status := "✅"
	if !meetsReq {
		status = "⚠️"
	}

	return fmt.Sprintf("%s %d characters, %d words", status, chars, words)
}

// Helper type for repository information
type RepoInfo struct {
	Name            string
	Description     string
	PrimaryLanguage string
}
