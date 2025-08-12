package git

import (
	"testing"

	"github.com/frfahim/gitstory/internal/testutil"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfo_Success(t *testing.T) {
	repo, testRepo := setupTestRepo(t)
	defer testRepo.Cleanup()

	info, err := repo.GetInfo()
	require.NoError(t, err)

	assert.Equal(t, testRepo.Dir, info.Path)
	assert.True(t, info.IsGitRepo)
	assert.NotEmpty(t, info.CurrentBranch)
	assert.Greater(t, info.CommitCount, 0)
	assert.Empty(t, info.Error)
}

func TestGetInfo_EmptyRepository(t *testing.T) {
	testRepo := testutil.CreateTestRepo(t)
	defer testRepo.Cleanup()
	// Don't add any commits

	repo, err := OpenRepository(testRepo.Dir)
	require.NoError(t, err)

	info, err := repo.GetInfo()
	require.NoError(t, err)

	assert.Equal(t, testRepo.Dir, info.Path)
	assert.True(t, info.IsGitRepo)
	assert.Equal(t, 0, info.CommitCount)
	// CurrentBranch might be empty for repository without commits
}

func TestDetectDefaultBranch_MainExists(t *testing.T) {
	repo, testRepo := setupTestRepo(t)
	defer testRepo.Cleanup()

	// Create a "main" branch
	worktree, err := testRepo.Repo.Worktree()
	require.NoError(t, err)

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
		Create: true,
	})
	require.NoError(t, err)

	defaultBranch, err := repo.DetectDefaultBranch()
	require.NoError(t, err)

	assert.Equal(t, "main", defaultBranch)
}

func TestDetectDefaultBranch_MultipleBranches(t *testing.T) {
	repo, testRepo := setupTestRepo(t)
	defer testRepo.Cleanup()

	worktree, err := testRepo.Repo.Worktree()
	require.NoError(t, err)

	// Create multiple branches without main/master
	branches := []string{"development", "feature", "staging"}
	for _, branchName := range branches {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branchName),
			Create: true,
		})
		require.NoError(t, err)
	}

	defaultBranch, err := repo.DetectDefaultBranch()
	require.NoError(t, err)

	// Should return one of the branches (fallback behavior)
	assert.NotEmpty(t, defaultBranch)
	assert.Contains(t, append(branches, "master"), defaultBranch) // master might be the initial branch
}
