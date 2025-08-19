# GitStory 

**Transform your Git commits into compelling stories with AI.**

GitStory is a powerful CLI tool that analyzes your Git history and uses AI to generate intelligent summaries optimized for different platforms - from technical documentation to social media posts.

## 🎯 Why GitStory?

As developers, we often struggle to:
- **Summarize sprint work** for standups and retrospectives
- **Create release notes** that actually make sense
- **Share achievements** on LinkedIn or Twitter
- **Document changes** for technical teams
- **Write blog posts** about interesting projects

GitStory solves this by analyzing your Git commits and generating platform-specific summaries using AI.

## ✨ Features

### 🤖 AI-Powered Summarization
- **Multiple AI Providers**: OpenAI, Google Gemini (Claude and more coming soon)
- **Smart Provider Detection**: Auto-detects available API keys
- **Context-Aware**: Uses your commit messages, file changes, and custom context

### 📱 Platform-Optimized Output
- **Twitter**: Engaging posts under 280 characters with hashtags
- **LinkedIn**: Professional summaries highlighting business impact
- **Blog**: Detailed technical narratives for your blog
- **Technical Docs**: Comprehensive documentation for teams
- **Personal Notes**: Organized summaries for your reference

### 🔧 Git Integration
- **Repository Analysis**: Extract commit history and metadata
- **Branch Comparison**: Summarize unique commits on feature branches
- **Flexible Filtering**: Analyze last N commits or specific date ranges

### ⚙️ Developer-Friendly
- **Smart Defaults**: Works with minimal configuration
- **Multiple Formats**: JSON, Markdown, plain text output
- **Environment Variables**: Easy CI/CD integration

## 🚀 Quick Start

### Installation

```bash
# Using Go (requires Go 1.21+)
go install github.com/frfahim/gitstory@latest

# Or clone and build
git clone https://github.com/frfahim/gitstory.git
cd gitstory
make install
```

### Setup API Keys

```bash
# Set your preferred AI provider
export OPENAI_API_KEY="your-openai-key"
# OR
export GEMINI_API_KEY="your-gemini-key"
```

## 💡 Examples

### Basic Usage

```bash
# Quick summary for Twitter
gitstory summarize --platform twitter

# Professional LinkedIn post
gitstory summarize --platform linkedin --commits 10

# Technical documentation
gitstory summarize --platform technical --include-diff

# Store summary to file
gitstory summarize --platform blog --output blog_summary.md

# Blog post with context
gitstory summarize --platform blog --context "Sprint 23: User Authentication Overhaul"
```

### 📱 Social Media Post
```bash
$ gitstory summarize --platform twitter --commits 3

🐦 Twitter Summary:
Just shipped user auth v2! 🚀 Added JWT tokens, OAuth integration, and fixed security vulnerabilities. Major improvement to user onboarding flow! #webdev #authentication #security
```

### 💼 Professional Update
```bash
$ gitstory summarize --platform linkedin --commits 5

💼 LinkedIn Summary:
Completed a significant authentication system overhaul this sprint...

Key technical achievements:
• Implemented JWT token authentication with bcrypt password hashing
• Added OAuth integration for Google and GitHub login
• Enhanced security with rate limiting and vulnerability fixes

This work improves our platform's security posture while reducing user friction. The OAuth integration alone should increase signup conversion by 40%.

What security improvements have made the biggest impact in your projects?
```

### 🔧 Technical Documentation
```bash
$ gitstory summarize --platform technical --commits 8

## Summary
Implemented comprehensive authentication system with OAuth integration and security enhancements.

## Technical Changes
- auth/jwt.go: JWT token implementation with refresh rotation
- auth/oauth.go: Google and GitHub OAuth providers  
- middleware/security.go: Rate limiting and password policies

## Impact
- Reduced login time by 60% with optimized token validation
- Enhanced security with bcrypt cost factor 12
- Improved user experience with social login options
```

## 📖 Command Reference

### Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `summarize` | Generate AI summaries | `gitstory summarize --platform blog` |
| `list` | Show repository info and commits | `gitstory list --commits 10` |
| `status` | Display repository status | `gitstory status` |

### Summarize Options

```bash
# Platform options
--platform twitter|linkedin|blog|technical|notes

# Provider options  
--provider openai|gemini|claude

# Commit selection
--commits N              # Last N commits (default: 5)
--since "1 week ago"     # Commits since date
--unique --base main     # Only commits unique to current branch

# Content options
--context "description"  # Add context for better summaries
--include-diff          # Include code changes in analysis

# Output options
--output file.md        # Save to file
--format json|markdown  # Output format
```

### Configuration

```bash
# Auto-detect available providers
gitstory summarize --platform blog
# Uses first available: OpenAI → Gemini → Claude

# Override provider
gitstory summarize --provider gemini --platform twitter

# Set custom model
export OPENAI_MODEL="gpt-5"
export GEMINI_MODEL="gemini-2.5-pro"
```

## 🎯 Use Cases

### For Individual Developers
- **Daily Standups**: Quickly summarize yesterday's work
- **Weekly Updates**: Create professional progress reports
- **Social Presence**: Share achievements on Twitter/LinkedIn
- **Portfolio**: Generate project descriptions for your website
- **Learning Log**: Document technical decisions and learnings

### For Teams
- **Sprint Reviews**: Summarize team accomplishments  
- **Release Notes**: Generate user-facing change descriptions
- **Technical Documentation**: Create clear change logs
- **Stakeholder Updates**: Translate technical work to business value

### For Content Creators
- **Blog Posts**: Start drafts from actual development work
- **Tutorials**: Document real implementation approaches
- **Case Studies**: Extract insights from project development


## 🛠 Current Status

### ✅ Implemented Features
- ✅ Git repository analysis and commit extraction
- ✅ OpenAI GPT integration 
- ✅ Google Gemini integration
- ✅ Platform-specific prompt optimization
- ✅ Flexible commit filtering (count, date range, unique commits)
- ✅ Multiple output formats (JSON, Markdown, plain text)
- ✅ Comprehensive test suite

### 🚧 Future Plan
- [ ] **Claude AI Integration**: Anthropic's Claude and more
- [ ] **Interactive TUI**: Beautiful terminal interface with Bubble Tea
- [ ] **Export System**: Direct export to Hugo, Jekyll, Obsidian
- [ ] **Diff Analysis**: Include actual code changes in summaries
- [ ] **Configuration Management**: Persistent settings and API key management
- [ ] **Template System**: Custom prompt templates for different use cases


## 🏗 Technical Architecture

```
gitstory/
├── cmd/                 # CLI commands (cobra-based)
│   ├── list.go         # Repository analysis
│   ├── summarize.go    # AI summarization  
│   └── status.go       # Status information
├── internal/
│   ├── git/            # Git operations (go-git)
│   │   ├── repository.go
│   │   ├── commits.go
│   │   └── info.go
│   ├── llm/            # AI/LLM integration
│   │   ├── types.go    # Common interfaces
│   │   ├── client.go   # Provider factory
│   │   ├── prompts.go  # Shared prompt system
│   │   ├── openai.go   # OpenAI implementation
│   │   └── gemini.go   # Google Gemini implementation
│   └── testutil/       # Testing utilities
└── Makefile           # Build and test automation
```

### Key Dependencies
- **Git Operations**: `go-git/go-git` for repository analysis
- **CLI Framework**: `spf13/cobra` for command structure
- **AI Providers**: `openai-go` and `google.golang.org/genai`
- **Testing**: `stretchr/testify` for comprehensive testing

## 🧪 Development

### Prerequisites
- Go 1.21+ 
- Git
- API keys for OpenAI and/or Google Gemini

### Development Setup

```bash
# Clone repository
git clone https://github.com/frfahim/gitstory.git
cd gitstory

# Install dependencies
go mod download

# Run tests
make test

# Run with coverage
make test-coverage

# Build local binary
make build

# Install to GOPATH
make install

# Format and lint
make fmt
make lint
```

### Testing
```bash
# Run all tests
make test

# Test with real API (set keys first)
export OPENAI_API_KEY="your-key"
export GEMINI_API_KEY="your-key"  
go run main.go summarize --platform twitter
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**GitStory** - Because every commit tells a story worth sharing! 🚀✨
