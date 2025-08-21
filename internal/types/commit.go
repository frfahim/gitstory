package types

// CommitStats represents statistical information about a commit
type CommitStats struct {
	TotalFiles  int      `json:"total_files"`
	TotalLines  int      `json:"total_lines"`
	Additions   int      `json:"additions"`
	Deletions   int      `json:"deletions"`
	Languages   []string `json:"languages"`
	PrimaryLang string   `json:"primary_lang"`
}

// FileChange represents a file modification in a commit
type FileChange struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	Content   string `json:"content,omitempty"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// CommitData represents standardized commit information for AI consumption
type CommitData struct {
	Hash    string       `json:"hash"`
	Message string       `json:"message"`
	Author  string       `json:"author"`
	Date    string       `json:"date"`
	Stats   CommitStats  `json:"stats"`
	Files   []FileChange `json:"files"`
}
