// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// cache of commit message list keyed by page name - use this because retrieving the git commit log is potentially slow
var commitMessagesByPage map[string][]string

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
	log.Info("cloning wiki repository %s into directory %s", accessor.wikiRepoURL, accessor.wikiRepoDir)

	repository, err := git.PlainClone(accessor.wikiRepoDir, isBare, &git.CloneOptions{
		URL:               accessor.wikiRepoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		err = errors.Wrapf(err, "cloning repository %s into directory %s", accessor.wikiRepoURL, accessor.wikiRepoDir)
		return err
	}

	accessor.wikiRepo = repository

	// reset the commit log cache
	commitMessagesByPage = make(map[string][]string)

	return nil
}

// CommitWikiToRepo stages any files added or updated since the last commit then commits them to our cloned wiki repo.
// We package the staging and commit together here because it is easier than embedding hooks to do the git staging
// deep into the wiki parsing process where files from the Trac worksapce can get copied over on-the-fly.
func (accessor *DefaultAccessor) CommitWikiToRepo(author string, authorEMail string, message string) error {
	worktree, err := accessor.wikiRepo.Worktree()
	if err != nil {
		err = errors.Wrapf(err, "retrieving git work tree for cloned wiki")
		return err
	}

	status, err := worktree.Status()
	for file, filestatus := range status {
		worktreeStatus := filestatus.Worktree
		if worktreeStatus == git.Untracked || worktreeStatus == git.Modified {
			_, err = worktree.Add(file)
			if err != nil {
				err = errors.Wrapf(err, "adding file %s to git work tree", file)
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
		err = errors.Wrapf(err, "committing changes to git for cloned wiki")
		return err
	}

	return nil
}

// CopyFileToWiki copies an external file into the Gitea Wiki, returning a URL through which the file can be viewed/
func (accessor *DefaultAccessor) CopyFileToWiki(externalFilePath string, giteaWikiRelPath string) error {
	_, err := os.Stat(externalFilePath)
	if os.IsNotExist(err) {
		log.Warn("cannot copy non-existant file referenced from Wiki: \"%s\"", externalFilePath)
		return nil
	}

	giteaPath := filepath.Join(accessor.wikiRepoDir, giteaWikiRelPath)
	giteaDir := path.Dir(giteaPath)
	err = os.MkdirAll(giteaDir, 0775)
	if err != nil {
		err = errors.Wrapf(err, "creating directory %s in cloned wiki", giteaDir)
		return err
	}

	_, err = os.Stat(giteaPath)
	if accessor.overwrite || !os.IsExist(err) {
		copyFile(externalFilePath, giteaPath)
		log.Debug("copied file %s to wiki path %s", externalFilePath, giteaWikiRelPath)
	}

	return nil
}

// commitLog returns the log of commits for the given file.
func (accessor *DefaultAccessor) commitLog(pageName string) ([]string, error) {
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
		err = errors.Wrapf(err, "retrieving git log for file %s", wikiFilename)
		return nil, err
	}

	var commitMessages []string
	err = commitIter.ForEach(func(commit *object.Commit) error {
		commitMessages = append(commitMessages, commit.Message)
		return nil
	})

	return commitMessages, nil
}

// pageCommitExists determines whether or not a commit of the given page exists with a commit message containing the provided string
func (accessor *DefaultAccessor) pageCommitExists(pageName string, commitString string) (bool, error) {
	commitMessages, haveCommitMessages := commitMessagesByPage[pageName]
	if !haveCommitMessages {
		pageCommitMessages, err := accessor.commitLog(pageName)
		if err != nil {
			return false, err
		}
		commitMessagesByPage[pageName] = pageCommitMessages
		commitMessages = pageCommitMessages
	}

	for _, commitMessage := range commitMessages {
		if strings.Contains(commitMessage, commitString) {
			return true, nil
		}
	}

	return false, nil
}

// WriteWikiPage writes (a version of) a wiki page to the checked-out wiki repository, returning the path to the written file.
func (accessor *DefaultAccessor) WriteWikiPage(pageName string, markdownText string, commitMarker string) (bool, error) {
	// if we're not explicitly overwriting, look for conflicting previous commit of wiki page
	if !accessor.overwrite {
		hasCommit, err := accessor.pageCommitExists(pageName, commitMarker)
		if err != nil {
			return false, err
		}
		if hasCommit {
			return false, nil
		}
	}

	pagePath := filepath.Join(accessor.wikiRepoDir, wikiPageFileName(pageName))
	pageDir := path.Dir(pagePath)
	err := os.MkdirAll(pageDir, 0775)
	if err != nil {
		err = errors.Wrapf(err, "creating directory %s for page in cloned wiki", pageDir)
		return false, err
	}

	file, err := os.Create(pagePath)
	if err != nil {
		err = errors.Wrapf(err, "creating wiki page file %s", pagePath)
		return false, err
	}

	file.WriteString(markdownText)
	log.Debug("wrote version of wiki page %s", pageName)
	return true, nil
}

// TranslateWikiPageName translates a Trac wiki page name into a Gitea one
func (accessor *DefaultAccessor) TranslateWikiPageName(pageName string) string {
	// special case: Trac "WikiStart" page is Gitea "Home" page...
	if pageName == "WikiStart" {
		return "Home"
	}

	return pageName
}

// commitWikiRepo commits all wiki repository changes by pushing all changes to the local wiki repository back to the remote.
// (Ff pushing the wiki is disabled, the local repository is left and a message is output)
func (accessor *DefaultAccessor) commitWikiRepo() error {
	if !accessor.pushWiki {
		log.Info("wiki updates have been committed to cloned repository %s; please review changes and push back to remote when done.", accessor.wikiRepoDir)
		return nil
	}

	auth := &http.BasicAuth{
		Username: accessor.userName,
		Password: accessor.wikiRepoToken,
	}

	log.Debug("pushing wiki to remote")
	err := accessor.wikiRepo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		err = errors.Wrapf(err, "pushing cloned wiki to remote")
		return err
	}

	log.Debug("deleting cloned wiki directory %s after pushing it", accessor.wikiRepoDir)
	return os.RemoveAll(accessor.wikiRepoDir)
}

// rollbackWikiRepo rolls back all changes to the wiki repository by deleting the local cloned repository without pushing it.
func (accessor *DefaultAccessor) rollbackWikiRepo() error {
	log.Debug("rolling back all wiki repo changes by deleting cloned wiki directory %s", accessor.wikiRepoDir)
	return os.RemoveAll(accessor.wikiRepoDir)
}
