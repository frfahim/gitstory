package git

import (
	"github.com/go-git/go-git/v5"
)

// Repository holds a git.Repository and its path
type Repository struct {
	repo *git.Repository
	path string
}

// OpenRepository attempts to open a Git repository at the given path
func OpenRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repository{repo: repo, path: path}, nil
}

// IsGitRepository checks if the given path is a Git repository
func IsGitRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

func (r *Repository) CurrentBranchName() string {
	head, err := r.repo.Head()
	if err != nil {
		return "(unknown)"
	}
	return head.Name().Short()
}
