package wiki

import (
	"github.com/go-git/go-git"
)

// Accessor provides access to a Gitea Wiki repository.
type Accessor struct {
	repoURL string
	repoDir string
	repo    *git.Repository
}

// CreateAccessor retirurns a new Gitea Wiiki accessor.
func CreateAccessor(wikiRepoURL string, wikiRepoDir string) *Accessor {
	accessor := Accessor{repoURL: wikiRepoURL, repoDir: wikiRepoDir, repo: nil}
	return &accessor
}
