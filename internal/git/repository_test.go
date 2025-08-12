package git

import (
	"testing"

	"github.com/frfahim/gitstory/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenRepository_Success(t *testing.T) {
	repo := testutil.CreateTestRepo(t)
	defer repo.Cleanup()

	gitRepo, err := OpenRepository(repo.Dir)
	require.NoError(t, err)
	assert.NotNil(t, gitRepo)
}

func TestOpenRepository_NotGitRepo(t *testing.T) {
	tempDir := t.TempDir()

	gitRepo, err := OpenRepository(tempDir)
	assert.Error(t, err)
	assert.Nil(t, gitRepo)
}

func TestOpenRepository_InvalidPath(t *testing.T) {
	gitRepo, err := OpenRepository("/invalid/path")
	assert.Error(t, err)
	assert.Nil(t, gitRepo)
}
