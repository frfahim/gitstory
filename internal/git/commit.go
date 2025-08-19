package git

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// ListCommits returns the latest N commits from the repository
func (r *Repository) ListCommits(n int) ([]*object.Commit, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return nil, err
	}
	iter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	var commits []*object.Commit
	count := 0
	err = iter.ForEach(func(c *object.Commit) error {
		if count >= n {
			return storer.ErrStop
		}
		commits = append(commits, c)
		count++
		return nil
	})
	return commits, err
}

// ListCommitSummaries returns summary info for last N commits
func (repo *Repository) ListCommitSummaries(commits []*object.Commit) ([]CommitSummary, error) {
	var summaries []CommitSummary
	for _, commit := range commits {
		commitSummary := CommitSummary{
			Hash:    commit.Hash.String()[:7],
			Author:  commit.Author.Name,
			Date:    commit.Author.When.Format(time.RFC3339),
			Message: commit.Message,
		}
		details, _ := repo.GetCommitDiffDetails(commit, true)
		commitSummary.Files = details.Files
		commitSummary.Stats = details.Stats
		summaries = append(summaries, commitSummary)
	}
	return summaries, nil
}

// ListUniqueCommits returns commits unique to the current branch (not in baseBranch)
func (r *Repository) ListUniqueCommits(baseBranch string, n int) ([]*object.Commit, error) {
	ref, err := r.repo.Head()
	if err != nil {
		return nil, err
	}
	baseRef, err := r.repo.Reference(plumbing.NewBranchReferenceName(baseBranch), true)
	if err != nil {
		return nil, fmt.Errorf("base branch '%s' not found: %w", baseBranch, err)
	}
	mergeBase, err := FindMergeBase(r.repo, ref.Hash(), baseRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to find merge-base: %w", err)
	}

	iter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var commits []*object.Commit
	count := 0
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash == mergeBase {
			return storer.ErrStop
		}
		if count >= n {
			return storer.ErrStop
		}
		commits = append(commits, c)
		count++
		return nil
	})
	return commits, err
}
