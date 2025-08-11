package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repository struct {
	repo *git.Repository
	path string
}

type RepoInfo struct {
	Path          string
	IsGitRepo     bool
	CurrentBranch string
	RemoteURL     string
	CommitCount   int
	Error         string
}

// OpenRepository attempts to open a Git repository at the given path
func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("not a git repository (or any of the parent directories): %w", err)
	}

	return &Repository{
		repo: repo,
		path: path,
	}, nil
}

// GetInfo returns basic information about the repository
func (r *Repository) GetInfo() (*RepoInfo, error) {
	info := &RepoInfo{
		Path:      r.path,
		IsGitRepo: true,
	}

	// Get current branch
	head, err := r.repo.Head()
	if err == nil {
		info.CurrentBranch = head.Name().Short()
	}

	// Get remote URL (if exists)
	remotes, err := r.repo.Remotes()
	if err == nil && len(remotes) > 0 {
		urls := remotes[0].Config().URLs
		if len(urls) > 0 {
			info.RemoteURL = urls[0]
		}
	}

	// Get commit count (basic implementation)
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

// IsGitRepository checks if the given path is a Git repository
func IsGitRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}
