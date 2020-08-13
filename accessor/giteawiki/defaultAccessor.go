package giteawiki

import (
	"gopkg.in/src-d/go-git.v4"
)

// DefaultAccessor is the default implementation of the Gitea Wiki Accessor, accessing the wiki via its git repository.
type DefaultAccessor struct {
	repoURL string
	repoDir string
	repo    *git.Repository
}

// CreateDefaultAccessor retirurns a new Gitea Wiiki accessor.
func CreateDefaultAccessor(wikiRepoURL string, wikiRepoDir string) *DefaultAccessor {
	accessor := DefaultAccessor{repoURL: wikiRepoURL, repoDir: wikiRepoDir, repo: nil}
	return &accessor
}
