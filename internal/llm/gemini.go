package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	config ClientConfig
}

func (c *GeminiClient) GetProvider() Provider {
	return Gemini
}

func NewGeminiClient(config ClientConfig) (*GeminiClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required for Gemini client")
	}
	if config.Model == "" {
		config.Model = "gemini-2.5-flash-lite"
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		config: config,
	}, nil
}

// ValidateCredentials checks if the API key is valid
func (c *GeminiClient) ValidateCredentials(ctx context.Context) error {

	// Simple test request to validate credentials
	response, err := c.client.Models.GenerateContent(ctx, c.config.Model, genai.Text("Hello"), nil)
	if err != nil {
		return fmt.Errorf("invalid Gemini credentials: %w", err)
	}

	if response == nil {
		return fmt.Errorf("empty response from Gemini API")
	}

	return nil
}

func (c *GeminiClient) Summarize(ctx context.Context, request *SummaryRequest) (*SummaryResponse, error) {
	systemPrompt := getSystemPrompt(request.Platform)
	userPrompt := buildPrompt(request)
	// For Gemini, we need to combine system and user prompts since it doesn't have separate system messages
	// We'll structure it as: "You are X. Here's the task: Y"
	fullPrompt := fmt.Sprintf("%s\n\n--- TASK ---\n%s", systemPrompt, userPrompt)

	config := &genai.GenerateContentConfig{
		Temperature:     genai.Ptr[float32](0.7),
		MaxOutputTokens: int32(getMaxTokensForPlatform(request.Platform)),
	}
	result, err := c.client.Models.GenerateContent(ctx, c.config.Model, genai.Text(fullPrompt), config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}
	if result == nil || len(result.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	return &SummaryResponse{
		Platform: request.Platform,
		Summary:  result.Text(),
	}, nil
}
