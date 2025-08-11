package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// FindMergeBase returns the hash of the common ancestor between two commits
func FindMergeBase(repo *git.Repository, hash1, hash2 plumbing.Hash) (plumbing.Hash, error) {
	seen := map[plumbing.Hash]struct{}{}
	queue := []plumbing.Hash{hash1}
	for len(queue) > 0 {
		h := queue[0]
		queue = queue[1:]
		seen[h] = struct{}{}
		commit, err := repo.CommitObject(h)
		if err != nil {
			continue
		}
		queue = append(queue, commit.ParentHashes...)
	}

	queue2 := []plumbing.Hash{hash2}
	for len(queue2) > 0 {
		h := queue2[0]
		queue2 = queue2[1:]
		if _, ok := seen[h]; ok {
			return h, nil
		}
		commit, err := repo.CommitObject(h)
		if err != nil {
			continue
		}
		queue2 = append(queue2, commit.ParentHashes...)
	}
	return plumbing.ZeroHash, fmt.Errorf("no common ancestor found")
}
