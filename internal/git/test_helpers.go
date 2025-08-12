package git

import (
	"testing"

	"github.com/frfahim/gitstory/internal/testutil"
	"github.com/stretchr/testify/require"
)

// setupTestRepo creates a test repo with commits and returns the Repository
func setupTestRepo(t *testing.T) (*Repository, *testutil.TestRepo) {
	testRepo := testutil.CreateTestRepo(t)
	testRepo.AddMultipleCommits(t)

	repo, err := OpenRepository(testRepo.Dir)
	require.NoError(t, err)

	return repo, testRepo
}

// setupEmptyTestRepo creates a test repo without commits
func setupEmptyTestRepo(t *testing.T) (*Repository, *testutil.TestRepo) {
	testRepo := testutil.CreateTestRepo(t)

	repo, err := OpenRepository(testRepo.Dir)
	require.NoError(t, err)

	return repo, testRepo
}
