package gitea

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"stevejefferson.co.uk/trac2gitea/log"
)

// GetWikiAttachmentRelPath returns the location of an attachment to Trac a wiki page when stored in the Gitea wiki repository.
// The returned path is relative to the root of the Gitea wiki repository.
func (accessor *DefaultAccessor) GetWikiAttachmentRelPath(pageName string, filename string) string {
	return filepath.Join("attachments", pageName, filename)
}

// GetWikiHtdocRelPath returns the location of a given Trac 'htdocs' file when stored in the Gitea wiki repository.
// The returned path is relative to the root of the Gitea wiki repository.
func (accessor *DefaultAccessor) GetWikiHtdocRelPath(filename string) string {
	return filepath.Join("htdocs", filename)
}

// GetWikiFileURL returns a URL for viewing a file stored in the Gitea wiki repository.
func (accessor *DefaultAccessor) GetWikiFileURL(relpath string) string {
	//FIXME: we want a path to the "raw" wiki repository here - this is my best guess at what this should be but sadly it does not work
	return "../raw/" + relpath
}

// CloneWiki clones our wiki repo to the provided directory.
func (accessor *DefaultAccessor) CloneWiki() {
	isBare := false
	repository, err := git.PlainClone(accessor.wikiRepoDir, isBare, &git.CloneOptions{
		URL:               accessor.wikiRepoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Fatal(err)
	}

	accessor.wikiRepo = repository
}

// CommitWiki stages any files added or updated since the last commit then commits them to our cloned wiki repo.
// We package the staging and commit together here because it is easier than embedding hooks to do the git staging
// deep into the wiki parsing process where files from the Trac worksapce can get copied over on-the-fly.
func (accessor *DefaultAccessor) CommitWiki(author string, authorEMail string, message string) {
	worktree, err := accessor.wikiRepo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	status, err := worktree.Status()
	for file, filestatus := range status {
		worktreeStatus := filestatus.Worktree
		if worktreeStatus == git.Untracked || worktreeStatus == git.Modified {
			_, err = worktree.Add(file)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	_, err = worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  author,
			Email: authorEMail,
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

// PushWiki pushes all changes to the local wiki repository back to the remote.
func (accessor *DefaultAccessor) PushWiki() {
	// TODO: not working yet

	// auth := &http.BasicAuth{
	// 	Username: accessor.userName,
	// 	Password: "[git_basic_auth_token]",
	// }

	// err := accessor.wikiRepo.Push(&git.PushOptions{
	// 	RemoteName: "origin",
	// 	Auth:       auth,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Infof("Trac wiki has been imported into cloned wiki repository at %s. Please review changes and push back to remote when done.\n",
		accessor.wikiRepoDir)
}

// CopyFileToWiki copies an external file into the Gitea Wiki, returning a URL through which the file can be viewed/
func (accessor *DefaultAccessor) CopyFileToWiki(externalFilePath string, giteaWikiRelPath string) {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warnf("cannot copy non-existant file referenced from Wiki: \"%s\"\n", externalFilePath)
		return
	}

	giteaPath := filepath.Join(accessor.wikiRepoDir, giteaWikiRelPath)
	err = os.MkdirAll(path.Dir(giteaPath), 0775)
	if err != nil {
		log.Fatal(err)
	}

	// determine whether file already exists - if it does we'll just assume we've already copied it earlier in the conversion
	_, err = os.Stat(giteaPath)
	if !os.IsExist(err) {
		accessor.copyFile(externalFilePath, giteaPath)
		log.Debugf("Copied file %s to wiki path %s\n", externalFilePath, giteaWikiRelPath)
	}
}

// WriteWikiPage writes (a version of) a wiki page to the checked-out wiki repository, returning the path to the written file.
func (accessor *DefaultAccessor) WriteWikiPage(pageName string, markdownText string) string {
	pagePath := filepath.Join(accessor.wikiRepoDir, pageName+".md")
	file, err := os.Create(pagePath)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(markdownText)

	log.Debugf("Wrote version of wiki page %s\n", pageName)

	return pagePath
}

// TranslateWikiPageName translates a Trac wiki page name into a Gitea one
func (accessor *DefaultAccessor) TranslateWikiPageName(pageName string) string {
	// special case: Trac "WikiStart" page is Gitea "Home" page...
	if pageName == "WikiStart" {
		return "Home"
	}

	return pageName
}
