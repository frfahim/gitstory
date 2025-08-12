package testutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

type TestRepo struct {
	Dir  string
	Repo *git.Repository
}

func CreateTestRepo(t *testing.T) *TestRepo {
	tempDir, err := os.MkdirTemp("", "gitstory-test-*")
	require.NoError(t, err)

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	return &TestRepo{
		Dir:  tempDir,
		Repo: repo,
	}
}

func (tr *TestRepo) Cleanup() {
	if tr.Dir != "" {
		os.RemoveAll(tr.Dir)
	}
}

func (tr *TestRepo) AddCommit(t *testing.T, filename, content, message string) {
	worktree, err := tr.Repo.Worktree()
	require.NoError(t, err)

	filePath := filepath.Join(tr.Dir, filename)
	err = os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	_, err = worktree.Add(filename)
	require.NoError(t, err)

	signature := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
		When:  time.Now(),
	}

	_, err = worktree.Commit(message, &git.CommitOptions{
		Author: signature,
	})
	require.NoError(t, err)
}

func (tr *TestRepo) AddMultipleCommits(t *testing.T) {
	commits := []struct {
		filename string
		content  string
		message  string
	}{
		{"README.md", "# Test Project", "Initial commit"},
		{"main.go", "package main\n\nfunc main() {}", "Add main.go"},
		{"config.yaml", "version: 1.0", "Add configuration"},
	}

	for _, commit := range commits {
		tr.AddCommit(t, commit.filename, commit.content, commit.message)
	}
}
