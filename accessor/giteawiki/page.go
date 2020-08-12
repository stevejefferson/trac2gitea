package giteawiki

import (
	"log"
	"os"
	"path/filepath"
)

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

// TranslatePageName translates a Trac wiki page name into a Gitea one
func (accessor *Accessor) TranslatePageName(pageName string) string {
	// special case: Trac "WikiStart" page is Gitea "Home" page...
	if pageName == "WikiStart" {
		return "Home"
	}

	return pageName
}
