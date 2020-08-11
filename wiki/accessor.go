package wiki

import (
	"log"
	"os"
	"path/filepath"

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

// WritePage writes (a version of) a wiki page to the checked-out wiki repository, returning the path to the written file.
func (accessor *Accessor) WritePage(pageName string, markdownText string) string {
	pagePath := filepath.Join(accessor.repoDir, pageName+".md")
	file, err := os.Create(pagePath)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(markdownText)

	return pagePath
}
