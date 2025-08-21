package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/frfahim/gitstory/internal/types"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitDiffDetails struct {
	Files []types.FileChange
	Stats types.CommitStats
}

func (r *Repository) GetCommitDiffDetails(commit *object.Commit, includeDiff bool) (CommitDiffDetails, error) {
	commitDiff, stats, err := r.extractFileChanges(commit, includeDiff)
	return CommitDiffDetails{
		Files: commitDiff,
		Stats: stats,
	}, err

}

func (r *Repository) extractFileChanges(commit *object.Commit, includeDiff bool) ([]types.FileChange, types.CommitStats, error) {
	var files []types.FileChange
	var stats types.CommitStats
	languageCount := make(map[string]int)

	// Get the current commit tree object
	currentTree, err := commit.Tree()
	if err != nil {
		return files, stats, fmt.Errorf("failed to get current commit (%s) tree: %w", commit.Hash, err)
	}
	// Get the parent commit and it's tree
	parentCommit, err := r.getParentCommit(commit)
	if err != nil {
		// TODO: handle when parent commit is empty
		return nil, stats, fmt.Errorf("failed to get parent commit from commit(%s): %w", commit.Hash, err)
	}
	parentTree, err := parentCommit.Tree()
	if err != nil {
		return files, stats, fmt.Errorf("failed to get parent commit (%s) tree: %w", parentCommit.Hash, err)
	}

	// Get the file changes between the parent and current commit
	fileChanges, err := parentTree.Diff(currentTree)
	if err != nil {
		return files, stats, fmt.Errorf("failed to get commit diff: %w", err)
	}

	stats.TotalFiles = len(fileChanges)
	// Collect file change statistics
	for _, change := range fileChanges {
		fileChange := r.processFileChange(change)
		files = append(files, fileChange)
		stats.TotalLines += fileChange.Additions + fileChange.Deletions
		stats.Additions += fileChange.Additions
		stats.Deletions += fileChange.Deletions
		lang := r.detectLanguage(fileChange.Path)
		if lang != "" {
			languageCount[lang]++
		}
	}

	// Find the primary language and add all used languages
	maxLangCount := 0
	for lang, count := range languageCount {
		stats.Languages = append(stats.Languages, lang)
		if count > maxLangCount {
			maxLangCount = count
			stats.PrimaryLang = lang
		}
	}
	return files, stats, nil
}

// Process a single file change
func (r *Repository) processFileChange(change *object.Change) types.FileChange {
	var path, status string
	var additionCount, deletionCount int = 0, 0

	action, _ := change.Action()
	status = action.String()
	if status == "Delete" {
		path = change.From.Name
	} else {
		path = change.To.Name
	}
	patch, err := change.Patch()
	if err != nil {
		return types.FileChange{
			Path:   path,
			Status: status,
		}
	}
	// Get the file statistics from the patch
	for _, fs := range patch.Stats() {
		additionCount += fs.Addition
		deletionCount += fs.Deletion
	}

	fileChange := types.FileChange{
		Path:      path,
		Status:    status,
		Additions: additionCount,
		Deletions: deletionCount,
		Content:   r.extractChangesOnly(patch),
	}

	return fileChange
}

// Extract the changes from a patch String
func (r *Repository) extractChangesOnly(patch *object.Patch) string {
	var additions, deletions, result strings.Builder

	lines := strings.Split(patch.String(), "\n")

	for _, line := range lines {
		// Only extract actual changes, skip headers and context
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			additions.WriteString(line + "\n")
		}
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletions.WriteString(line + "\n")
		}
	}

	if deletions.Len() > 0 {
		result.WriteString("DELETIONS:\n")
		result.WriteString(deletions.String())
		result.WriteString("\n")
	}
	if additions.Len() > 0 {
		result.WriteString("ADDITIONS:\n")
		result.WriteString(additions.String())
	}

	return strings.TrimSpace(result.String())
}

// Get the parent commit
func (r *Repository) getParentCommit(commit *object.Commit) (*object.Commit, error) {
	if commit.NumParents() == 0 {
		return nil, fmt.Errorf("commit (%s) has no parent", commit.Hash)
	}
	parentCommit, err := commit.Parent(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent commit (%s): %w", commit.Hash, err)
	}
	return parentCommit, nil
}

func (r *Repository) detectLanguage(filePath string) string {
	// Use the file extension to determine the programming language
	ext := strings.ToLower(filepath.Ext(filePath))

	// Map file extensions to programming languages
	languageMap := map[string]string{
		".go":         "Go",
		".js":         "JavaScript",
		".ts":         "TypeScript",
		".py":         "Python",
		".java":       "Java",
		".c":          "C",
		".cpp":        "C++",
		".rs":         "Rust",
		".php":        "PHP",
		".rb":         "Ruby",
		".swift":      "Swift",
		".kt":         "Kotlin",
		".dart":       "Dart",
		".cs":         "C#",
		".scala":      "Scala",
		".clj":        "Clojure",
		".html":       "HTML",
		".css":        "CSS",
		".scss":       "SCSS",
		".sass":       "Sass",
		".sql":        "SQL",
		".sh":         "Shell",
		".yaml":       "YAML",
		".yml":        "YAML",
		".json":       "JSON",
		".xml":        "XML",
		".md":         "Markdown",
		".dockerfile": "Docker",
	}

	if lang, exists := languageMap[ext]; exists {
		return lang
	}

	// Check for specific filenames
	filename := strings.ToLower(filepath.Base(filePath))
	if filename == "dockerfile" {
		return "Docker"
	}
	if filename == "makefile" {
		return "Make"
	}

	return ""
}
