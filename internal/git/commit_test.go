package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCommits_Success(t *testing.T) {
	repo, testRepo := setupTestRepo(t)
	defer testRepo.Cleanup()

	commits, err := repo.ListCommits(10)
	require.NoError(t, err)
	assert.NotEmpty(t, commits)
}

func TestListCommits_WithLimit(t *testing.T) {
	repo, testRepo := setupTestRepo(t)
	defer testRepo.Cleanup()

	commits, err := repo.ListCommits(1)
	require.NoError(t, err)
	assert.Len(t, commits, 1)
}
