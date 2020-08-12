package giteawiki

import (
	"gopkg.in/src-d/go-git.v4"
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
