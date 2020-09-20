// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"fmt"
	"strings"
	"testing"

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

func expectToWriteGiteaWikiPage(
	t *testing.T,
	tracWikiPage *trac.WikiPage,
	giteaWikiPage string,
	pageWritten bool) {
	// expect to convert Trac page to markdown
	markdownText := "trac wiki " + tracWikiPage.Text + "converted to markdown"
	mockMarkdownConverter.
		EXPECT().
		WikiConvert(tracWikiPage.Name, tracWikiPage.Text).
		Return(markdownText)

	// expect to write translated page to Gitea, returning provided status
	mockGiteaAccessor.
		EXPECT().
		WriteWikiPage(giteaWikiPage, markdownText, gomock.Any()).
		Return(pageWritten, nil)
}

func expectToCommitGiteaWikiPage(
	t *testing.T,
	tracWikiPage *trac.WikiPage,
	giteaPageAuthor string,
	giteaAuthorEmail string) {
	// expect to lookup email address of Gitea user as author of commit
	mockGiteaAccessor.
		EXPECT().
		GetUserEMailAddress(giteaPageAuthor).
		Return(giteaAuthorEmail, nil)

	// expect to commit Gitea wiki page including Trac page name, version and commit comment in Gitea comment
	mockGiteaAccessor.
		EXPECT().
		CommitWikiToRepo(giteaPageAuthor, giteaAuthorEmail, gomock.Any()).
		DoAndReturn(func(author string, email string, comment string) error {
			assertTrue(t, strings.Contains(comment, tracWikiPage.Name))
			assertTrue(t, strings.Contains(comment, fmt.Sprintf("%d", tracWikiPage.Version)))
			assertTrue(t, strings.Contains(comment, tracWikiPage.Comment))
			return nil
		})
}

func TestImportOfPredefinedSingleVersionWikiPage(t *testing.T) {
	setUpWiki(t)
	defer tearDown(t)

	// clone existing Gitea wiki
	expectCloneWiki(t)

	// trac should return us a single wiki page and no attachments
	expectTracToReturnWikiPages(t, tracWikiPage1v1)
	expectTracToReturnWikiAttachments(t)

	// trac wiki page is a predefined one
	expectToTestForPredefinedWikiPage(t, tracWikiPage1v1, true)

	dataImporter.ImportWiki(userMap)
}

func TestImportOfPredefinedSingleVersionWikiPageWhenConvertingPredefinedPages(t *testing.T) {
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

	// write and commit wiki page
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)

	predefinedPageDataImporter.ImportWiki(userMap)
}

func TestImportOfSingleVersionWikiPage(t *testing.T) {
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

	// write and commit wiki page
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)

	dataImporter.ImportWiki(userMap)
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

	// write and commit wiki page
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1v2Author, giteaWikiPage1v2Author)

	dataImporter.ImportWiki(userMap)
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

	// write and commit wiki pages
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1v2Author, giteaWikiPage1v2Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage2v1, giteaWikiPage2, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage2v1, giteaWikiPage2v1Author, giteaWikiPage2v1Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage2v2, giteaWikiPage2, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage2v2, giteaWikiPage2v2Author, giteaWikiPage2v2Author)

	dataImporter.ImportWiki(userMap)
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

	// fail to write wiki page
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, false)

	// ...do not expect to commit wiki page

	dataImporter.ImportWiki(userMap)
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

	// fail top write first version of wiki page
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, false)

	// successfully write and commit second version of page
	expectToWriteGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1v2Author, giteaWikiPage1v2Author)

	dataImporter.ImportWiki(userMap)
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

	dataImporter.ImportWiki(userMap)
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

	dataImporter.ImportWiki(userMap)
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

	dataImporter.ImportWiki(userMap)
}

func TestImportOfMultipleVersionsOfMultipleWikiPagesWithMultipleAttachments(t *testing.T) {
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

	// write and commit wiki pages
	expectToWriteGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v1, giteaWikiPage1v1Author, giteaWikiPage1v1Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage1v2, giteaWikiPage1v2Author, giteaWikiPage1v2Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage2v1, giteaWikiPage2, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage2v1, giteaWikiPage2v1Author, giteaWikiPage2v1Author)
	expectToWriteGiteaWikiPage(t, tracWikiPage2v2, giteaWikiPage2, true)
	expectToCommitGiteaWikiPage(t, tracWikiPage2v2, giteaWikiPage2v2Author, giteaWikiPage2v2Author)

	dataImporter.ImportWiki(userMap)
}
