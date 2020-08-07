package wiki

import (
	"log"
	"os"
	"path/filepath"
)

// Accessor provides access to a Gitea Wiki repository.
type Accessor struct {
	wikiRepoDir string
}

// CreateAccessor retirurns a new Gitea Wiiki accessor./
func CreateAccessor(repoDir string) *Accessor {
	stat, err := os.Stat(repoDir)
	if err != nil {
		log.Fatal(err)
	}
	if stat.IsDir() != true {
		log.Fatal("Gitea wiki repo directory is not a directory: ", repoDir)
	}

	accessor := Accessor{wikiRepoDir: repoDir}
	return &accessor
}

// WritePageVersion writes a version of a given page to the Gitea Wiki.
func (accessor *Accessor) WritePageVersion(pageName string, markdownText string, version int64, comment string, time int64) {
	pagePath := filepath.Join(accessor.wikiRepoDir, pageName+".md")
	file, err := os.Create(pagePath)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(markdownText)
}
