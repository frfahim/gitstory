package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type RepoInfo struct {
	Path          string
	IsGitRepo     bool
	CurrentBranch string
	RemoteURL     string
	CommitCount   int
	Error         string
}

// GetInfo returns basic information about the repository
func (r *Repository) GetInfo() (*RepoInfo, error) {
	info := &RepoInfo{
		Path:      r.path,
		IsGitRepo: true,
	}

	// Current branch
	head, err := r.repo.Head()
	if err == nil {
		info.CurrentBranch = head.Name().Short()
	}

	// Remote URL
	remotes, err := r.repo.Remotes()
	if err == nil && len(remotes) > 0 {
		urls := remotes[0].Config().URLs
		if len(urls) > 0 {
			info.RemoteURL = urls[0]
		}
	}

	// Commit count
	commitIter, err := r.repo.Log(&git.LogOptions{})
	if err == nil {
		count := 0
		commitIter.ForEach(func(c *object.Commit) error {
			count++
			return nil
		})
		info.CommitCount = count
	}

	return info, nil
}
