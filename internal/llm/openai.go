package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	client *openai.Client
	config ClientConfig
}

// GetProvider returns the provider type
func (c *OpenAIClient) GetProvider() Provider {
	return OpenAI
}

// NewOpenAIClient creates a new OpenAI client using official package
func NewOpenAIClient(config ClientConfig) (Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	// Set default model
	if config.Model == "" {
		config.Model = "gpt-4o-mini"
	}

	// Create client using official OpenAI package
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	return &OpenAIClient{
		client: &client,
		config: config,
	}, nil
}

func (c *OpenAIClient) Summarize(ctx context.Context, request *SummaryRequest) (*SummaryResponse, error) {
	systemPrompt := getSystemPrompt(request.Platform)
	userContext := buildPrompt(request)
	fullPrompt := fmt.Sprintf("%s\n\n%s", systemPrompt, userContext)

	// Prepare messages for OpenAI chat completion
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(fullPrompt),
	}

	// Call OpenAI API
	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:               c.config.Model,
		Messages:            messages,
		Temperature:         param.Opt[float64]{Value: 0.7},
		MaxCompletionTokens: param.Opt[int64]{Value: int64(getMaxTokensForPlatform(request.Platform))},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI API")
	}
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)

	return &SummaryResponse{
		Summary:  summary,
		Platform: request.Platform,
	}, nil
}
