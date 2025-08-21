package llm

import (
	"fmt"
	"strings"
)

// getSystemPrompt returns the system prompt for different platforms
func getSystemPrompt(platform Platform) string {
	prompts := map[Platform]string{
		Blog: `You are an experienced technical writer and software engineering blogger. You excel at:
- Translating technical work into engaging narratives
- Highlighting the "why" behind code changes, not just the "what"
- Creating content that both developers and technical managers find valuable
- Using clear, accessible language while maintaining technical accuracy
- Structuring content for easy scanning and comprehension`,

		Twitter: `You are a tech influencer who creates viral developer content on Twitter/X. You excel at:
- Condensing complex technical work into compelling 280-character stories
- Using relevant hashtags and emojis strategically
- Creating content that gets developers to engage and share
- Balancing technical accuracy with accessibility
- Highlighting achievements and learnings that resonate with the dev community`,

		LinkedIn: `You are a senior software engineer who shares professional insights on LinkedIn. You excel at:
- Highlighting business impact of technical work
- Demonstrating professional growth and technical leadership
- Creating content that showcases both technical skills and business acumen
- Writing posts that attract recruiters and technical peers
- Balancing technical details with broader professional relevance`,

		Technical: `You are a senior technical lead creating documentation for other developers. You excel at:
- Providing clear, actionable technical insights
- Explaining architectural decisions and their rationale
- Highlighting implementation details that matter for future development
- Creating documentation that helps team members understand changes quickly
- Focusing on technical impact, performance implications, and maintainability`,

		Note: `You are a thoughtful developer creating personal development notes. You excel at:
- Organizing information for easy future reference
- Highlighting key learning points and decisions made
- Creating concise but complete summarize
- Noting important context and follow-up actions
- Structuring information for personal productivity and growth tracking`,
	}

	if prompt, exists := prompts[platform]; exists {
		return prompt
	}
	return prompts[Technical] // default
}

// getPlatformInstructions returns platform-specific formatting instructions
func getPlatformInstructions(platform Platform) string {
	instructions := map[Platform]string{
		Blog: `
Create a blog post summary with this structure:
## What We Accomplished
- Lead with the main achievement or problem solved
- Use engaging, story-driven language

## Key Technical Highlights  
- 2-3 most significant technical changes
- Focus on interesting implementation details
- Mention technologies/frameworks used

## Impact & Why It Matters
- Business value or user benefit
- Technical improvements (performance, maintainability, etc.)
- What this enables for future development

Use markdown formatting. Aim for 300-500 words. Make it engaging but informative.`,

		Twitter: `
Create a Twitter/X thread or single post:
- Start with a hook that grabs attention
- Maximum 280 characters if single post, or 2-3 connected tweets
- Include 2-3 relevant hashtags (#coding #webdev #javascript etc.)
- Use 1-2 emojis strategically (🚀 ✨ 🔧 💡 🎯)
- Focus on the most impressive achievement or learning
- End with engagement (question, call to action, or relatable statement)

Examples:
"Just shipped user auth v2! 🚀 Reduced login time by 60% with smart caching and JWT optimization. Sometimes the smallest changes make the biggest impact 💡 #webdev #performance"`,

		LinkedIn: `
Create a professional LinkedIn post:
- Start with a professional hook about the business challenge or opportunity
- Highlight 2-3 key technical achievements and their business impact
- Mention specific technologies/skills used (great for keyword visibility)
- Include a learning or insight that adds professional value
- End with a question or call-to-action to encourage engagement
- Use professional tone but keep it conversational
- Aim for 150-300 words
- Consider using bullet points for readability

Structure: Challenge/Opportunity → Technical Solution → Business Impact → Personal Learning → Engagement Question`,

		Technical: `
Create comprehensive technical documentation:

## Summary
- Brief overview of what was accomplished

## Technical Changes
- List major code/architecture changes
- Include file/component names where relevant
- Mention new dependencies or libraries added

## Implementation Details
- Explain key technical decisions and their rationale
- Highlight any complex problem-solving approaches
- Note performance improvements or optimizations

## Breaking Changes & Migration Notes
- List any breaking changes
- Provide migration guidance if needed

## Testing & Quality
- Mention testing approach or coverage improvements
- Note any quality/security enhancements

## Next Steps
- List any follow-up work or technical debt created

Use clear, scannable formatting. Include code snippets or technical details where helpful.`,

		Note: `
Create organized personal notes:

## Summary
- What was accomplished in this work session

## Key Changes
- Most important modifications made
- Technologies/approaches used

## Decisions Made
- Important technical or architectural decisions
- Rationale behind choices made

## Learnings
- New things learned during implementation
- Challenges overcome and how

## Follow-up
- [ ] Tasks to complete later
- [ ] Technical debt created
- [ ] Ideas for future improvements

Use bullet points and checkboxes. Keep it concise but complete for future reference.`,
	}

	if instruction, exists := instructions[platform]; exists {
		return instruction
	}
	return instructions[Technical] // default
}

// getMaxTokensForPlatform returns appropriate token limits for each platform
func getMaxTokensForPlatform(platform Platform) int {
	limits := map[Platform]int{
		Twitter:   150,  // Short
		LinkedIn:  400,  // Professional detail
		Blog:      1000, // Rich content
		Technical: 800,  // Detailed but focused
		Note:      500,  // Personal note
	}

	if limit, exists := limits[platform]; exists {
		return limit
	}
	return 400 // default
}

// buildPrompt creates the prompt for any AI provider based on commits and context
func buildPrompt(request *SummaryRequest) string {
	var prompt strings.Builder

	// Add context if provided
	if request.UserContext != "" {
		prompt.WriteString(fmt.Sprintf("Project Context: %s\n\n", request.UserContext))
	}

	// Add commit summary stats
	prompt.WriteString(fmt.Sprintf("Analyzing %d git commit(s) with code changes:\n\n", len(request.Commits)))

	// Add each commit with enhanced formatting
	for i, commit := range request.Commits {
		prompt.WriteString(fmt.Sprintf("=== Commit %d ===\n", i+1))
		// prompt.WriteString(fmt.Sprintf("• Hash: %s\n", commit.Hash))
		prompt.WriteString(fmt.Sprintf("• Author: %s\n", commit.Author))
		prompt.WriteString(fmt.Sprintf("• Date: %s\n", commit.Date))
		prompt.WriteString(fmt.Sprintf("• Message: %s\n", commit.Message))

		// Add file statistics (always available from ListCommitSummarize)
		if commit.Stats.TotalFiles > 0 {
			prompt.WriteString(fmt.Sprintf("• Files changed: %d\n", commit.Stats.TotalFiles))
			prompt.WriteString(fmt.Sprintf("• Lines: +%d -%d\n", commit.Stats.Additions, commit.Stats.Deletions))

			if commit.Stats.PrimaryLang != "" {
				prompt.WriteString(fmt.Sprintf("• Primary language: %s\n", commit.Stats.PrimaryLang))
			}

			if len(commit.Stats.Languages) > 1 {
				prompt.WriteString(fmt.Sprintf("• Languages: %s\n", strings.Join(commit.Stats.Languages, ", ")))
			}
		}

		// Add file changes (always available from ListCommitSummarize)
		if len(commit.Files) > 0 {
			prompt.WriteString("• File changes:\n")
			for _, file := range commit.Files {
				prompt.WriteString(fmt.Sprintf("  - %s (%s)", file.Path, file.Status))
				if file.Additions > 0 || file.Deletions > 0 {
					prompt.WriteString(fmt.Sprintf(" [+%d -%d]", file.Additions, file.Deletions))
				}
				prompt.WriteString("\n")

				// Include code changes if available
				if file.Content != "" {
					prompt.WriteString("    Code changes:\n")
					// Limit the content for prompt efficiency
					contentLines := strings.Split(file.Content, "\n")
					maxLines := 15 // Reasonable limit for prompts
					if len(contentLines) > maxLines {
						contentLines = contentLines[:maxLines]
						contentLines = append(contentLines, "... (truncated)")
					}

					for _, line := range contentLines {
						if line != "" {
							prompt.WriteString(fmt.Sprintf("    %s\n", line))
						}
					}
					prompt.WriteString("\n")
				}
			}
		}
		prompt.WriteString("\n")
	}

	// Add platform-specific instructions
	prompt.WriteString("Please create a summary following these guidelines:")
	prompt.WriteString(getPlatformInstructions(request.Platform))

	// Add code-specific instructions (always relevant since we always have code changes)
	prompt.WriteString("\n\nCode Analysis Instructions:")
	prompt.WriteString("\n- Focus on the actual code changes and their impact")
	prompt.WriteString("\n- Identify new features, bug fixes, refactoring, or optimizations")
	prompt.WriteString("\n- Mention specific functions, classes, or modules when relevant")
	prompt.WriteString("\n- Highlight technical improvements or architectural changes")
	prompt.WriteString("\n- Consider the programming languages and technologies involved")

	return prompt.String()
}
