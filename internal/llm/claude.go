package llm

import (
	"context"
	"fmt"
)

type ClaudeClient struct {
	client interface{} // Placeholder for actual Claude client type
	config ClientConfig
}

func (c *ClaudeClient) GetProvider() Provider {
	return Claude
}

func NewClaudeClient(config ClientConfig) (*ClaudeClient, error) {
	return nil, fmt.Errorf("claude isn't yet implemented")
}

func (c *ClaudeClient) Summarize(ctx context.Context, request *SummaryRequest) (*SummaryResponse, error) {
	return nil, fmt.Errorf("claude isn't yet implemented")
}
