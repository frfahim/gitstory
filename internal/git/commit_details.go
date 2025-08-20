package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func (r *Repository) GetCommitDiffDetails(commit *object.Commit, includeDiff bool) (CommitDiffDetails, error) {
	commitDiff, stats, err := r.extractFileChanges(commit, includeDiff)
	return CommitDiffDetails{
		Files: commitDiff,
		Stats: stats,
	}, err

}

func (r *Repository) extractFileChanges(commit *object.Commit, includeDiff bool) ([]FileChange, CommitStats, error) {
	var files []FileChange
	var stats CommitStats

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
	}
	return files, stats, nil
}

// Process a single file change
func (r *Repository) processFileChange(change *object.Change) FileChange {
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
		return FileChange{
			Path:   path,
			Status: status,
		}
	}
	// Get the file statistics from the patch
	for _, fs := range patch.Stats() {
		additionCount += fs.Addition
		deletionCount += fs.Deletion
	}

	fileChange := FileChange{
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
