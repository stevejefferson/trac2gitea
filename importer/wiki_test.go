// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

const (
	tracWikiPage1 = "Page1"
	tracWikiPage2 = "Page2"

	tracWikiPage1Attachment1Path = "/path/page1/attachment1"
	tracWikiPage1Attachment2Path = "/path/page1/attachment2"
	tracWikiPage2Attachment1Path = "/path/page2/attachment1"
	tracWikiPage2Attachment2Path = "/path/page2/attachment2"

	giteaWikiPage1Attachment1Path = "/giteapath/gitea_page1/gitea_attachment1"
	giteaWikiPage1Attachment2Path = "/giteapath/gitea_page1/gitea_attachment2"
	giteaWikiPage2Attachment1Path = "/giteapath/gitea_page2/gitea_attachment1"
	giteaWikiPage2Attachment2Path = "/giteapath/gitea_page2/gitea_attachment2"
)

var (
	tracWikiPage1Attachment1 *trac.WikiAttachment
	tracWikiPage1Attachment2 *trac.WikiAttachment
	tracWikiPage2Attachment1 *trac.WikiAttachment
	tracWikiPage2Attachment2 *trac.WikiAttachment
)

func setUpWikiAttachments(t *testing.T) {
	setUp(t)

	tracWikiPage1Attachment1 = &trac.WikiAttachment{PageName: tracWikiPage1, FileName: "attachment1.file"}
	tracWikiPage1Attachment2 = &trac.WikiAttachment{PageName: tracWikiPage1, FileName: "attachment2.file"}
	tracWikiPage2Attachment1 = &trac.WikiAttachment{PageName: tracWikiPage2, FileName: "attachment3.file"}
	tracWikiPage2Attachment2 = &trac.WikiAttachment{PageName: tracWikiPage2, FileName: "attachment4.file"}
}

const (
	tracWikiPage1v1Author = "user1"
	tracWikiPage1v2Author = "user2 <user2@abc.def>"
	tracWikiPage2v1Author = "user3 <user3@bcd.efg>"
	tracWikiPage2v2Author = "user4"

	giteaWikiPage1v1Author = "gitea_usr1"
	giteaWikiPage1v2Author = "gitea_usr2"
	giteaWikiPage2v1Author = "gitea_usr3"
	giteaWikiPage2v2Author = "gitea_usr4"

	giteaWikiPage1v1AuthorEmail = "g_usr1@abcd.efg"
	giteaWikiPage1v2AuthorEmail = "g_usr2@bcde.fgh"
	giteaWikiPage2v1AuthorEmail = "g_usr3@cdef.ghi"
	giteaWikiPage2v2AuthorEmail = "g_usr4@defg.hij"

	giteaWikiPage1 = "Gitea_Page1"
	giteaWikiPage2 = "Gitea_Page_2"
)

var (
	tracWikiPage1v1 *trac.WikiPage
	tracWikiPage1v2 *trac.WikiPage
	tracWikiPage2v1 *trac.WikiPage
	tracWikiPage2v2 *trac.WikiPage
)

func setUpWiki(t *testing.T) {
	setUpWikiAttachments(t)

	tracWikiPage1v1 = &trac.WikiPage{
		Name:       "Page1",
		Author:     tracWikiPage1v1Author,
		Text:       "this is the text of version 1 of wiki page #1",
		Comment:    "Page#1 was created",
		Version:    1,
		UpdateTime: 12345}
	tracWikiPage1v2 = &trac.WikiPage{
		Name:       "Page1",
		Author:     tracWikiPage1v2Author,
		Text:       "this is the text of version 2 of wiki page #1",
		Comment:    "Page#1 was updated",
		Version:    2,
		UpdateTime: 23456}

	tracWikiPage2v1 = &trac.WikiPage{
		Name:       "Page2",
		Author:     tracWikiPage2v1Author,
		Text:       "this is the text of version 1 of wiki page number 2",
		Comment:    "Page#2 was created",
		Version:    1,
		UpdateTime: 56789}
	tracWikiPage2v2 = &trac.WikiPage{
		Name:       "Page2",
		Author:     tracWikiPage2v2Author,
		Text:       "this is the text of version 2 of wiki page number 2",
		Comment:    "Page#2 was updated",
		Version:    2,
		UpdateTime: 67890}

	userMap[tracWikiPage1v1Author] = giteaWikiPage1v1Author
	userMap[tracWikiPage1v2Author] = giteaWikiPage1v2Author
	userMap[tracWikiPage2v1Author] = giteaWikiPage2v1Author
	userMap[tracWikiPage2v2Author] = giteaWikiPage2v2Author
}

func expectCloneWiki(t *testing.T) {
	mockGiteaAccessor.
		EXPECT().
		CloneWiki().
		Return(nil)
}

func expectTracToReturnWikiAttachments(t *testing.T, wikiAttachments ...*trac.WikiAttachment) {
	mockTracAccessor.
		EXPECT().
		GetWikiAttachments(gomock.Any()).
		DoAndReturn(func(handler func(attachment *trac.WikiAttachment) error) error {
			for _, wikiAttachment := range wikiAttachments {
				handler(wikiAttachment)
			}
			return nil
		})
}

func expectTracToReturnWikiPages(t *testing.T, wikiPages ...*trac.WikiPage) {
	mockTracAccessor.
		EXPECT().
		GetWikiPages(gomock.Any()).
		DoAndReturn(func(handler func(page *trac.WikiPage) error) error {
			for _, wikiPage := range wikiPages {
				handler(wikiPage)
			}
			return nil
		})
}

func expectToCopyTracWikiAttachmentToGitea(t *testing.T, wikiAttachment *trac.WikiAttachment, tracWikiAttachmentPath string, giteaWikiAttachmentPath string) {
	// expect to retrieve path to Trac wiki attachment
	mockTracAccessor.
		EXPECT().
		GetWikiAttachmentPath(wikiAttachment).
		Return(tracWikiAttachmentPath)

	// expect to retrieve path to store Wiki attachment in Gitea
	mockGiteaAccessor.
		EXPECT().
		GetWikiAttachmentRelPath(wikiAttachment.PageName, wikiAttachment.FileName).
		Return(giteaWikiAttachmentPath)

	// expect to copy attachment
	mockGiteaAccessor.
		EXPECT().
		CopyFileToWiki(tracWikiAttachmentPath, giteaWikiAttachmentPath)
}

func expectToTestForPredefinedWikiPage(t *testing.T, tracWikiPage *trac.WikiPage, isPredefined bool) {
	mockTracAccessor.
		EXPECT().
		IsPredefinedPage(tracWikiPage.Name).
		Return(isPredefined)
}

func expectToTranslateWikiPageName(t *testing.T, tracWikiPage *trac.WikiPage, giteaWikiPage string) {
	mockGiteaAccessor.
		EXPECT().
		TranslateWikiPageName(tracWikiPage.Name).
		Return(giteaWikiPage)
}

func expectTracWikiPagesToHaveAlreadyBeenCommitted(t *testing.T, giteaWikiPage string, tracWikiPagesAlreadyCommittedToGitea ...*trac.WikiPage) {
	log := []string{"initial version of page " + giteaWikiPage}
	for _, tracWikiPage := range tracWikiPagesAlreadyCommittedToGitea {
		// note: log message header is used to determine whether a Trac wiki page has already been committed to Gitea
		// - format of this is entirely dependent on the implementation so we are forced to copy the relevant code here
		updateTimeStr := time.Unix(tracWikiPage.UpdateTime, 0)
		message := fmt.Sprintf("[Imported from Trac: page %s, version %d at %s]\n%s", tracWikiPage.Name, tracWikiPage.Version, updateTimeStr, tracWikiPage.Comment)
		log = append(log, message)
	}

	mockGiteaAccessor.
		EXPECT().
		LogWiki(giteaWikiPage).
		Return(log, nil)
}

func expectToWriteAndCommitGiteaWikiPage(
	t *testing.T,
	tracWikiPage *trac.WikiPage,
	giteaWikiPage string, giteaPageAuthor string, giteaAuthorEmail string) {
	// expect to convert Trac page to markdown
	markdownText := "trac wiki " + tracWikiPage.Text + "converted to markdown"
	mockMarkdownConverter.
		EXPECT().
		WikiConvert(tracWikiPage.Name, tracWikiPage.Text).
		Return(markdownText)

	// expect to write translated page to Gitea
	mockGiteaAccessor.
		EXPECT().
		WriteWikiPage(giteaWikiPage, markdownText).
		Return("path-to-wiki-file", nil)

	// expect to lookup email address of Gitea user as author of commit
	mockGiteaAccessor.
		EXPECT().
		GetUserEMailAddress(giteaPageAuthor).
		Return(giteaAuthorEmail, nil)

	// expect to commit Gitea wiki page including Trac page name, version and commit comment in Gitea comment
	mockGiteaAccessor.
		EXPECT().
		CommitWiki(giteaPageAuthor, giteaAuthorEmail, gomock.Any()).
		DoAndReturn(func(author string, email string, comment string) error {
			assertTrue(t, strings.Contains(comment, tracWikiPage.Name))
			assertTrue(t, strings.Contains(comment, fmt.Sprintf("%d", tracWikiPage.Version)))
			assertTrue(t, strings.Contains(comment, tracWikiPage.Comment))
			return nil
		})
}

func expectToPushGiteaWiki(t *testing.T) {
	mockGiteaAccessor.
		EXPECT().
		PushWiki().
		Return(nil)
}

func TestImportWithoutPushOfPredefinedSingleVersionWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us a single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1)
	expectTracToReturnWikiAttachments(t)

	// trac wiki page is a predefined one
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, true)

	// ...do not expect any other actions

	dataImporter.ImportWiki(userMap, false)
}

func TestImportWithoutPushOfPredefinedSingleVersionWikiPageWhenConvertingPredefinedPages(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us a single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1)
	expectTracToReturnWikiAttachments(t)

	// do not test whether trac wiki page is a predefined one

	// translate to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)

	// no trac wiki page versions have yet been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1)

	// commit wiki page
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)

	// ...do not expect wiki to be pushed

	predefinedPageDataImporter.ImportWiki(userMap, false)
}

func TestImportWithoutPushOfSingleVersionWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us a single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1)
	expectTracToReturnWikiAttachments(t)

	// trac wiki page is not a predefined one
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)

	// translate to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)

	// no trac wiki page versions have yet been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1)

	// commit wiki page
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportWithPushOfSingleVersionWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us a single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1)
	expectTracToReturnWikiAttachments(t)

	// trac wiki page is not a predefined one
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)

	// translate to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)

	// no trac wiki page versions have yet been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1)

	// commit wiki page
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)

	// push wiki page
	expectToPushGiteaWiki(t)

	dataImporter.ImportWiki(userMap, true)
}

func TestImportOfMultiVersionWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us two versions of single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1, tracWikiPage1v2)
	expectTracToReturnWikiAttachments(t)

	// trac wiki pages are not predefined
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v2, false)

	// translate each version of page to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)
	expectToTranslateWikiPageName(t, tracWikiPage1v2, giteaWikiPage1)

	// no trac wiki page versions have yet been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1)

	// commit wiki pages
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, giteaWikiPage1v2Author, giteaWikiPage1v2Author)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportOfMultipleMultiVersionWikiPages(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us two versions of two wiki pages and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1, tracWikiPage1v2, tracWikiPage2v1, tracWikiPage2v2)
	expectTracToReturnWikiAttachments(t)

	// trac wiki pages are not predefined
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v2, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage2v1, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage2v2, false)

	// translate each version of page to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)
	expectToTranslateWikiPageName(t, tracWikiPage1v2, giteaWikiPage1)
	expectToTranslateWikiPageName(t, tracWikiPage2v1, giteaWikiPage2)
	expectToTranslateWikiPageName(t, tracWikiPage2v2, giteaWikiPage2)

	// no trac wiki page versions have yet been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1)
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage2)

	// commit wiki pages
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, giteaWikiPage1v2Author, giteaWikiPage1v2Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage2v1, giteaWikiPage2, giteaWikiPage2v1Author, giteaWikiPage2v1Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage2v2, giteaWikiPage2, giteaWikiPage2v2Author, giteaWikiPage2v2Author)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportOfAlreadyImportedWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us a single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1)
	expectTracToReturnWikiAttachments(t)

	// trac wiki page is not a predefined one
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)

	// translate to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)

	// trac wiki page version has already been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1, tracWikiPage1v1)

	// ...do not expect to commit wiki page

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportOfMultiVersionWikiPageWithOneAlreadyImportedVersion(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us two versions single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1, tracWikiPage1v2)
	expectTracToReturnWikiAttachments(t)

	// trac wiki page is not a predefined one
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v2, false)

	// translate to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)
	expectToTranslateWikiPageName(t, tracWikiPage1v2, giteaWikiPage1)

	// one trac wiki page version has already been committed to produce the Gitea wiki page
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1, tracWikiPage1v1)

	// commit version of wiki page that has not already been imported
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, giteaWikiPage1v2Author, giteaWikiPage1v2Author)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportOfSingleAttachmentToSingleWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us one attachment and no pages
	expectTracToReturnWikiPages(t)
	expectTracToReturnWikiAttachments(t, tracWikiPage1Attachment1)

	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment1, tracWikiPage1Attachment1Path, giteaWikiPage1Attachment1Path)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportOfMultipleAttachmentsToSingleWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us two attachments and no pages
	expectTracToReturnWikiPages(t)
	expectTracToReturnWikiAttachments(t, tracWikiPage1Attachment1, tracWikiPage1Attachment2)

	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment1, tracWikiPage1Attachment1Path, giteaWikiPage1Attachment1Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment2, tracWikiPage1Attachment2Path, giteaWikiPage1Attachment2Path)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportOfMultipleAttachmentsToMultipleWikiPages(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us four attachments and no pages
	expectTracToReturnWikiPages(t)
	expectTracToReturnWikiAttachments(t, tracWikiPage1Attachment1, tracWikiPage1Attachment2, tracWikiPage2Attachment1, tracWikiPage2Attachment2)

	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment1, tracWikiPage1Attachment1Path, giteaWikiPage1Attachment1Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment2, tracWikiPage1Attachment2Path, giteaWikiPage1Attachment2Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage2Attachment1, tracWikiPage2Attachment1Path, giteaWikiPage2Attachment1Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage2Attachment2, tracWikiPage2Attachment2Path, giteaWikiPage2Attachment2Path)

	// ...do not expect wiki to be pushed

	dataImporter.ImportWiki(userMap, false)
}

func TestImportAndPushOfMultipleVersionsOfMultipleWikiPagesWithMultipleAttachments(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us the full set of pages and attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1, tracWikiPage2v1, tracWikiPage2v2, tracWikiPage1v2)
	expectTracToReturnWikiAttachments(t, tracWikiPage1Attachment1, tracWikiPage1Attachment2, tracWikiPage2Attachment1, tracWikiPage2Attachment2)

	// attachments are copied
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment1, tracWikiPage1Attachment1Path, giteaWikiPage1Attachment1Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage1Attachment2, tracWikiPage1Attachment2Path, giteaWikiPage1Attachment2Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage2Attachment1, tracWikiPage2Attachment1Path, giteaWikiPage2Attachment1Path)
	expectToCopyTracWikiAttachmentToGitea(t, tracWikiPage2Attachment2, tracWikiPage2Attachment2Path, giteaWikiPage2Attachment2Path)

	// trac wiki pages are not predefined
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v2, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage2v1, false)
	expectToTestForPredefinedWikiPage(t, tracWikiPage2v2, false)

	// translate each version of page to markdown
	expectToTranslateWikiPageName(t, tracWikiPage1v1, giteaWikiPage1)
	expectToTranslateWikiPageName(t, tracWikiPage1v2, giteaWikiPage1)
	expectToTranslateWikiPageName(t, tracWikiPage2v1, giteaWikiPage2)
	expectToTranslateWikiPageName(t, tracWikiPage2v2, giteaWikiPage2)

	// no trac wiki page versions have yet been committed to produce any Gitea wiki pages
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage1)
	expectTracWikiPagesToHaveAlreadyBeenCommitted(t, giteaWikiPage2)

	// commit wiki pages
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, giteaWikiPage1v2Author, giteaWikiPage1v2Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage2v1, giteaWikiPage2, giteaWikiPage2v1Author, giteaWikiPage2v1Author)
	expectToWriteAndCommitGiteaWikiPage(t, tracWikiPage2v2, giteaWikiPage2, giteaWikiPage2v2Author, giteaWikiPage2v2Author)

	// push changes
	expectToPushGiteaWiki(t)

	dataImporter.ImportWiki(userMap, true)
}
