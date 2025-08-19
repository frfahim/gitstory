package git

type CommitSummary struct {
	Hash    string
	Author  string
	Date    string
	Message string
	Files   []FileChange
	Stats   CommitStats
}

type CommitDiffDetails struct {
	Files []FileChange
	Stats CommitStats
}

type FileChange struct {
	Path      string
	Status    string // Added, Modified, Deleted
	Additions int
	Deletions int
}

type CommitStats struct {
	TotalFiles int
	TotalLines int
	Additions  int
	Deletions  int
}
