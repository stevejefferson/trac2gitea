// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
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

func wikiPageFileName(pageName string) string {
	return pageName + ".md"
}

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
func (accessor *DefaultAccessor) CloneWiki() error {
	isBare := false
	log.Info("Cloning wiki repository %s into directory %s\n", accessor.wikiRepoURL, accessor.wikiRepoDir)

	repository, err := git.PlainClone(accessor.wikiRepoDir, isBare, &git.CloneOptions{
		URL:               accessor.wikiRepoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	accessor.wikiRepo = repository
	return nil
}

// LogWiki returns the log of commits for the given wiki file.
func (accessor *DefaultAccessor) LogWiki(pageName string) ([]string, error) {
	wikiFilename := wikiPageFileName(pageName)
	wikiFile := filepath.Join(accessor.wikiRepoDir, wikiFilename)

	// if file does not exist then we needn't look for its log...
	_, err := os.Stat(wikiFile)
	if os.IsNotExist(err) {
		noCommits := []string{}
		return noCommits, nil
	}

	commitIter, err := accessor.wikiRepo.Log(&git.LogOptions{FileName: &wikiFilename})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var commitMessages []string
	err = commitIter.ForEach(func(commit *object.Commit) error {
		commitMessages = append(commitMessages, commit.Message)
		return nil
	})

	return commitMessages, nil
}

// CommitWiki stages any files added or updated since the last commit then commits them to our cloned wiki repo.
// We package the staging and commit together here because it is easier than embedding hooks to do the git staging
// deep into the wiki parsing process where files from the Trac worksapce can get copied over on-the-fly.
func (accessor *DefaultAccessor) CommitWiki(author string, authorEMail string, message string) error {
	worktree, err := accessor.wikiRepo.Worktree()
	if err != nil {
		log.Error(err)
		return err
	}

	status, err := worktree.Status()
	for file, filestatus := range status {
		worktreeStatus := filestatus.Worktree
		if worktreeStatus == git.Untracked || worktreeStatus == git.Modified {
			_, err = worktree.Add(file)
			if err != nil {
				log.Error(err)
				return err
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
		log.Error(err)
		return err
	}

	return nil
}

// PushWiki pushes all changes to the local wiki repository back to the remote.
func (accessor *DefaultAccessor) PushWiki() error {
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
	// 	log.Error(err)
	//	return err
	// }

	log.Infof("Trac wiki has been imported into cloned wiki repository at %s. Please review changes and push back to remote when done.\n",
		accessor.wikiRepoDir)
	return nil
}

// CopyFileToWiki copies an external file into the Gitea Wiki, returning a URL through which the file can be viewed/
func (accessor *DefaultAccessor) CopyFileToWiki(externalFilePath string, giteaWikiRelPath string) error {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warnf("cannot copy non-existant file referenced from Wiki: \"%s\"\n", externalFilePath)
		return nil
	}

	giteaPath := filepath.Join(accessor.wikiRepoDir, giteaWikiRelPath)
	err = os.MkdirAll(path.Dir(giteaPath), 0775)
	if err != nil {
		log.Error(err)
		return err
	}

	// determine whether file already exists - if it does we'll just assume we've already copied it earlier in the conversion
	_, err = os.Stat(giteaPath)
	if !os.IsExist(err) {
		accessor.copyFile(externalFilePath, giteaPath)
		log.Debugf("Copied file %s to wiki path %s\n", externalFilePath, giteaWikiRelPath)
	}

	return nil
}

// WriteWikiPage writes (a version of) a wiki page to the checked-out wiki repository, returning the path to the written file.
func (accessor *DefaultAccessor) WriteWikiPage(pageName string, markdownText string) (string, error) {
	pagePath := filepath.Join(accessor.wikiRepoDir, wikiPageFileName(pageName))
	file, err := os.Create(pagePath)
	if err != nil {
		log.Error(err)
		return "", err
	}

	file.WriteString(markdownText)
	log.Debugf("Wrote version of wiki page %s\n", pageName)
	return pagePath, nil
}

// TranslateWikiPageName translates a Trac wiki page name into a Gitea one
func (accessor *DefaultAccessor) TranslateWikiPageName(pageName string) string {
	// special case: Trac "WikiStart" page is Gitea "Home" page...
	if pageName == "WikiStart" {
		return "Home"
	}

	return pageName
}
