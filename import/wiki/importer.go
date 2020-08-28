// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package wiki

import (
	"fmt"
	"strings"
	"time"

	"github.com/stevejefferson/trac2gitea/log"
	"github.com/stevejefferson/trac2gitea/markdown"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// Importer imports Trac Wiki data into a Gitea wiki repository.
type Importer struct {
	tracAccessor       trac.Accessor
	giteaAccessor      gitea.Accessor
	convertPredefineds bool
}

// CreateImporter creates a Trac wiki to Gitea wiki repository importer.
func CreateImporter(
	tAccessor trac.Accessor,
	gAccessor gitea.Accessor,
	convertPredefs bool) (*Importer, error) {

	importer := Importer{
		tracAccessor:       tAccessor,
		giteaAccessor:      gAccessor,
		convertPredefineds: convertPredefs}
	return &importer, nil
}

// ImportWiki imports a Trac wiki into a Gitea wiki repository.
func (importer *Importer) ImportWiki(userMap map[string]string, push bool) error {
	err := importer.giteaAccessor.CloneWiki()
	if err != nil {
		return err
	}

	importer.importWikiAttachments()
	importer.importWikiPages(userMap)

	if push {
		return importer.giteaAccessor.PushWiki()
	}

	log.Info("trac wiki has been imported into cloned wiki repository. Please review changes and push back to remote when done.")
	return nil
}

func (importer *Importer) importWikiAttachments() {
	importer.tracAccessor.GetWikiAttachments(func(attachment *trac.WikiAttachment) error {
		tracAttachmentPath := importer.tracAccessor.GetWikiAttachmentPath(attachment)
		giteaAttachmentPath := importer.giteaAccessor.GetWikiAttachmentRelPath(attachment.PageName, attachment.FileName)
		return importer.giteaAccessor.CopyFileToWiki(tracAttachmentPath, giteaAttachmentPath)
	})
}

// cache of commit message list keyed by page name - use this because 'LogWiki' is potentially slow
var commitMessagesByPage = make(map[string][]string)

// pageCommitExists determines whether or not a commit of the given page exists with a commit message containing the provided string
func (importer *Importer) pageCommitExists(pageName string, commitString string) (bool, error) {
	commitMessages, haveCommitMessages := commitMessagesByPage[pageName]
	if !haveCommitMessages {
		pageCommitMessages, err := importer.giteaAccessor.LogWiki(pageName)
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

func (importer *Importer) importWikiPages(userMap map[string]string) {
	importer.tracAccessor.GetWikiPages(func(page *trac.WikiPage) error {
		// skip predefined pages
		if !importer.convertPredefineds && importer.tracAccessor.IsPredefinedPage(page.Name) {
			log.Debug("skipping predefined Trac page %s", page.Name)
			return nil
		}

		// have we already converted this version of the trac wiki page?
		// - if so, skip it on the assumption that this is a re-import and that the only thing that is likely to have changed
		// is the addition of later trac versions of wiki pages - these will get added to the wiki repo as later versions
		updateTimeStr := time.Unix(page.UpdateTime, 0)
		tracPageVersionIdentifier := fmt.Sprintf("[Imported from Trac: page %s, version %d at %s]", page.Name, page.Version, updateTimeStr)
		translatedPageName := importer.giteaAccessor.TranslateWikiPageName(page.Name)
		hasCommit, err := importer.pageCommitExists(translatedPageName, tracPageVersionIdentifier)
		if err != nil {
			return err
		}
		if hasCommit {
			log.Info("wiki page %s: %s is already present in wiki - skipping...", translatedPageName, tracPageVersionIdentifier)
			return nil
		}

		// convert and write wiki page
		tracToMarkdownConverter := markdown.CreateWikiDefaultConverter(
			importer.tracAccessor, importer.giteaAccessor, page.Name)
		markdownText := tracToMarkdownConverter.Convert(page.Text)
		importer.giteaAccessor.WriteWikiPage(translatedPageName, markdownText)

		// find Gitea equivalent of Trac author
		giteaAuthor := userMap[page.Author]
		if giteaAuthor == "" {
			// can only happen if provided with faulty user-supplied map
			return fmt.Errorf("cannot find Gitea equivalent for trac author %s of wiki page %s", page.Author, page.Name)
		}
		giteaAuthorEmail, err := importer.giteaAccessor.GetUserEMailAddress(giteaAuthor)
		if err != nil {
			return err
		}

		// commit version of wiki page to local repository
		fullComment := tracPageVersionIdentifier + "\n\n" + page.Comment
		err = importer.giteaAccessor.CommitWiki(giteaAuthor, giteaAuthorEmail, fullComment)
		log.Info("wiki page %s: converted from Trac page %s, version %d", translatedPageName, page.Name, page.Version)
		return err
	})
}
