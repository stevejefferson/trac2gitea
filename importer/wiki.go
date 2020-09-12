// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"strings"
	"time"

	"github.com/stevejefferson/trac2gitea/log"

	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

func (importer *Importer) importWikiAttachments() {
	importer.tracAccessor.GetWikiAttachments(func(attachment *trac.WikiAttachment) error {
		tracAttachmentPath := importer.tracAccessor.GetWikiAttachmentPath(attachment)
		giteaAttachmentPath := importer.giteaAccessor.GetWikiAttachmentRelPath(attachment.PageName, attachment.FileName)
		return importer.giteaAccessor.CopyFileToWiki(tracAttachmentPath, giteaAttachmentPath)
	})
}

// cache of commit message list keyed by page name - use this because 'LogWiki' is potentially slow
var commitMessagesByPage map[string][]string

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
	// reset page commit log cache
	commitMessagesByPage = make(map[string][]string)

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
			log.Info("wiki page %s, version %d is already present in wiki - skipping...", translatedPageName, page.Version)
			return nil
		}

		// convert and write wiki page
		markdownText := importer.markdownConverter.WikiConvert(page.Name, page.Text)
		importer.giteaAccessor.WriteWikiPage(translatedPageName, markdownText)

		// find Gitea equivalent of Trac author if any
		author := page.Author
		authorEmail := ""
		giteaAuthor := userMap[page.Author]
		if giteaAuthor != "" {
			author = giteaAuthor
			authorEmail, err = importer.giteaAccessor.GetUserEMailAddress(giteaAuthor)
			if err != nil {
				return err
			}
		}

		// commit version of wiki page to local repository
		fullComment := tracPageVersionIdentifier + "\n\n" + page.Comment
		err = importer.giteaAccessor.CommitWiki(author, authorEmail, fullComment)
		log.Info("wiki page %s: converted from Trac page %s, version %d", translatedPageName, page.Name, page.Version)
		return err
	})
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
