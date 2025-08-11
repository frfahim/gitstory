package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type CommitSummary struct {
	Hash    string
	Author  string
	Date    string
	Message string
}

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
func (repo *Repository) ListCommitSummaries(n int) ([]CommitSummary, error) {
	commits, err := repo.ListCommits(n)
	if err != nil {
		return nil, err
	}
	var summaries []CommitSummary
	for _, c := range commits {
		summaries = append(summaries, CommitSummary{
			Hash:    c.Hash.String()[:7],
			Author:  c.Author.Name,
			Date:    c.Author.When.Format(time.RFC3339),
			Message: c.Message,
		})
	}
	return summaries, nil
}
