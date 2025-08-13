package llm

import (
	"fmt"
	"os"
	"strings"
)

// NewClient creates a new LLM client based on the provider
func NewClient(config ClientConfig) (Client, error) {
	// Auto-detect API key from environment if not provided
	if config.APIKey == "" {
		config.APIKey = getAPIKeyFromEnv(config.Provider)
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key not provided for %s. Set %s environment variable",
			config.Provider, getEnvKeyName(config.Provider))
	}

	if config.Model == "" {
		config.Model = getDefaultModel(config.Provider)
	}

	// Create provider-specific client
	switch config.Provider {
	case OpenAI:
		return NewOpenAIClient(config)
	case Gemini:
		return NewGeminiClient(config)
	case Claude:
		return nil, fmt.Errorf("Claude provider not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", config.Provider)
	}
}

// DetectAvailableProviders checks which providers have API keys configured
func DetectAvailableProviders() []Provider {
	var available []Provider

	providers := []Provider{OpenAI, Gemini, Claude}

	for _, provider := range providers {
		if getAPIKeyFromEnv(provider) != "" {
			available = append(available, provider)
		}
	}

	return available
}

// GetDefaultProvider returns the first available provider
func GetDefaultProvider() (Provider, error) {
	available := DetectAvailableProviders()

	if len(available) == 0 {
		return "", fmt.Errorf("no LLM providers configured. Please set an API key")
	}

	return available[0], nil
}

// getAPIKeyFromEnv retrieves API key from environment variables
func getAPIKeyFromEnv(provider Provider) string {
	envKeys := map[Provider]string{
		OpenAI: "OPENAI_API_KEY",
		Claude: "CLAUDE_API_KEY",
		Gemini: "GEMINI_API_KEY",
	}

	if envKey, exists := envKeys[provider]; exists {
		return strings.TrimSpace(os.Getenv(envKey))
	}
	return ""
}

// getEnvKeyName returns the environment variable name for a provider
func getEnvKeyName(provider Provider) string {
	envKeys := map[Provider]string{
		OpenAI: "OPENAI_API_KEY",
		Claude: "CLAUDE_API_KEY",
		Gemini: "GEMINI_API_KEY",
	}
	return envKeys[provider]
}

func getDefaultModel(provider Provider) string {
	defaultModels := map[Provider]string{
		OpenAI: "gpt-4o",
		Claude: "claude-3",
		Gemini: "gemini-2.5-flash-lite",
	}
	return defaultModels[provider]
}

// GetSupportedProviders returns list of all supported LLM providers
func GetSupportedProviders() []Provider {
	return []Provider{OpenAI, Gemini, Claude}
}

// GetSupportedPlatforms returns list of supported output platforms
func GetSupportedPlatforms() []Platform {
	return []Platform{Blog, Twitter, LinkedIn, Technical, Note}
}

// ValidateProvider checks if a provider string is valid
func ValidateProvider(provider string) error {
	for _, p := range GetSupportedProviders() {
		if string(p) == provider {
			return nil
		}
	}
	return fmt.Errorf("unsupported provider '%s'. Supported: %v", provider, GetSupportedProviders())
}

// ValidatePlatform checks if a platform string is valid
func ValidatePlatform(platform string) error {
	for _, p := range GetSupportedPlatforms() {
		if string(p) == platform {
			return nil
		}
	}
	return fmt.Errorf("unsupported platform '%s'. Supported: %v", platform, GetSupportedPlatforms())
}
