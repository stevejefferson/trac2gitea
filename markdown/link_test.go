// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package markdown_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
)

// functions returning trac and markdown formats for various types of link
func tracPlainLink(link string) string {
	return link
}

func tracSingleBracketLink(link string) string {
	return "[" + link + "]"
}

func tracSingleBracketLinkWithText(link string, text string) string {
	return "[" + link + " " + text + "]"
}

func tracDoubleBracketLink(link string) string {
	return "[[" + link + "]]"
}

func tracDoubleBracketLinkWithText(link string, text string) string {
	return "[[" + link + "|" + text + "]]"
}

func tracImageWithLink(image string, link string) string {
	return "[[Image(" + image + ", link=" + link + ")]]"
}

func tracImage(image string) string {
	return "[[Image(" + image + ")]]"
}

func markdownAutomaticLink(link string) string {
	return "<" + link + ">"
}

func markdownLinkWithText(link string, text string) string {
	return "[" + text + "](" + link + ")"
}

func markdownImage(image string) string {
	return "![](" + image + ")"
}

func markdownImageWithLink(image string, link string) string {
	return "[![](" + image + ")](" + link + ")"
}

func ticketConvert(tracText string) string {
	return converter.TicketConvert(ticketID, tracText)
}

func wikiConvert(tracText string) string {
	return converter.WikiConvert(wikiPage, tracText)
}

// verifyLink verifies that the provided trac formatting for a link + text results in the corresponding markdown format
func verifyLink(
	t *testing.T,
	setUpFn func(t *testing.T),
	tearDownFn func(t *testing.T),
	convertFn func(tracText string) string,
	tracFormatLink string,
	markdownFormatLink string) {
	setUpFn(t)
	defer tearDownFn(t)

	conversion := convertFn(leadingText + " " + tracFormatLink + " " + trailingText)
	assertEquals(t, conversion, leadingText+" "+markdownFormatLink+" "+trailingText)
}

const (
	linkText            = "text associated with link"
	additionalImageLink = "http://somewhere.com"
)

func verifyAllLinkTypes(
	t *testing.T,
	setUpFn func(t *testing.T),
	tearDownFn func(t *testing.T),
	convertFn func(tracText string) string,
	tracLinkStr string,
	markdownLinkStr string) {
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracPlainLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracSingleBracketLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracSingleBracketLinkWithText(tracLinkStr, linkText), markdownLinkWithText(markdownLinkStr, linkText))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracDoubleBracketLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracDoubleBracketLinkWithText(tracLinkStr, linkText), markdownLinkWithText(markdownLinkStr, linkText))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracImage(tracLinkStr), markdownImage(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracImageWithLink(tracLinkStr, additionalImageLink), markdownImageWithLink(markdownLinkStr, additionalImageLink))
	verifyLink(t, setUpFn, tearDownFn, convertFn, tracImageWithLink(additionalImageLink, tracLinkStr), markdownImageWithLink(additionalImageLink, markdownLinkStr))
}

const httpLink = "http://www.example.com"

func TestHttpLinks(t *testing.T) {
	verifyAllLinkTypes(t, setUp, tearDown, wikiConvert, httpLink, httpLink)
}

const httpsLink = "https://www.example.com"

func TestHttpsLink(t *testing.T) {
	verifyAllLinkTypes(t, setUp, tearDown, wikiConvert, httpsLink, httpsLink)
}

const (
	htdocFile      = "somefile.jpg"
	tracHtdocLink  = "htdocs:" + htdocFile
	giteaHtdocFile = "htdocs/" + htdocFile
	giteaHtdocURL  = "./somedir/" + htdocFile
)

func setUpHtdocs(t *testing.T) {
	setUp(t)

	// expect to have to retrieve full path to htdoc file within Trac workspace
	tracHtdocPath := filepath.Join("dir", "trac", "htdocs", htdocFile)
	mockTracAccessor.
		EXPECT().
		GetFullPath(gomock.Eq("htdocs"), gomock.Eq(htdocFile)).
		Return(tracHtdocPath)

	// expect to retrieve path where trac "htdocs" file will be stored in the Wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiHtdocRelPath(gomock.Eq(htdocFile)).
		Return(giteaHtdocFile)

	// expect to copy file into "htdocs" subdirectory of Wiki repo
	mockGiteaAccessor.
		EXPECT().
		CopyFileToWiki(gomock.Eq(tracHtdocPath), gomock.Eq(giteaHtdocFile))
	// expect to retrieve URL for viewing htdocs file in wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiFileURL(gomock.Eq(giteaHtdocFile)).
		Return(giteaHtdocURL)
}

func TestHtdocsLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpHtdocs,
		tearDown,
		wikiConvert,
		tracHtdocLink,
		giteaHtdocURL)
}

const (
	wikiPageName            = "SomeWikiPage"
	transformedWikiPageName = "TransformedWikiPage"
)

func setUpWikiLink(t *testing.T) {
	setUp(t)

	// expect call to translate name of wiki page
	mockGiteaAccessor.
		EXPECT().
		TranslateWikiPageName(gomock.Eq(wikiPageName)).
		Return(transformedWikiPageName)
}

func TestWikiUnprefixedLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		wikiConvert,
		wikiPageName,
		transformedWikiPageName)
}

func TestWikiPrefixedLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		wikiConvert,
		"wiki:"+wikiPageName,
		transformedWikiPageName)
}

const wikiPageAnchor = "page-anchor"

func TestWikiUnprefixedLinkWithAnchor(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		wikiConvert,
		wikiPageName+"#"+wikiPageAnchor,
		transformedWikiPageName+"#"+wikiPageAnchor)
}

func TestWikiPrefixedLinkWithAnchor(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		wikiConvert,
		"wiki:"+wikiPageName+"#"+wikiPageAnchor,
		transformedWikiPageName+"#"+wikiPageAnchor)
}

const (
	issueID  int64 = 26535
	issueURL       = "url-for-viewing-issue=26535"
)

var ticketIDStr = fmt.Sprintf("%d", ticketID)

func setUpAnyTicketLink(t *testing.T, tktID int64) {
	setUp(t)

	// expect call to lookup gitea issue for trac ticket
	mockGiteaAccessor.
		EXPECT().
		GetIssueID(gomock.Eq(tktID)).
		Return(issueID, nil)
}

func setUpTicketOnlyLink(t *testing.T) {
	setUpAnyTicketLink(t, ticketID)

	// expect call to lookup gitea issue URL
	mockGiteaAccessor.
		EXPECT().
		GetIssueURL(gomock.Eq(issueID)).
		Return(issueURL)
}

func TestTicketLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpTicketOnlyLink,
		tearDown,
		wikiConvert,
		"ticket:"+ticketIDStr,
		issueURL)
}

const (
	tracCommentNum    int64  = 12
	tracCommentNumStr        = "12"
	commentTime       int64  = 112233
	commentID         int64  = 54321
	commentURL        string = "url-of-comment-54321"
)

func setUpTicketCommentLink(t *testing.T, tktID int64) {
	setUpAnyTicketLink(t, tktID)

	// expect a call to lookup text of trac comment
	mockTracAccessor.
		EXPECT().
		GetTicketCommentTime(gomock.Eq(tktID), gomock.Eq(tracCommentNum)).
		Return(commentTime, nil)

	// expect call to lookup gitea ID for trac comment
	mockGiteaAccessor.
		EXPECT().
		GetIssueCommentIDsByTime(gomock.Eq(issueID), gomock.Eq(commentTime)).
		Return([]int64{commentID}, nil)

	// expect call to lookup URL of gitea comment
	mockGiteaAccessor.
		EXPECT().
		GetIssueCommentURL(gomock.Eq(issueID), gomock.Eq(commentID)).
		Return(commentURL)
}

func setUpImplicitTicketCommentLink(t *testing.T) {
	setUpTicketCommentLink(t, ticketID)
}

func TestImplicitTicketCommentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpImplicitTicketCommentLink,
		tearDown,
		ticketConvert,
		"comment:"+tracCommentNumStr,
		commentURL)
}

const (
	otherTicketID    int64  = 234567
	otherTicketIDStr string = "234567"
)

func setUpExplicitTicketCommentLink(t *testing.T) {
	setUpTicketCommentLink(t, otherTicketID)
}

func TestExplicitTicketCommentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpExplicitTicketCommentLink,
		tearDown,
		ticketConvert,
		"comment:"+tracCommentNumStr+":ticket:"+otherTicketIDStr,
		commentURL)
}

const (
	milestoneID   int64 = 678
	milestoneName       = "some-milestone"
	milestoneURL        = "url-for-viewing-milestone-678"
)

func setUpMilestoneLink(t *testing.T) {
	setUp(t)

	// expect call to lookup gitea milestone ID
	mockGiteaAccessor.
		EXPECT().
		GetMilestoneID(gomock.Eq(milestoneName)).
		Return(milestoneID, nil)

	// expect call to lookup URL for milestone
	mockGiteaAccessor.
		EXPECT().
		GetMilestoneURL(gomock.Eq(milestoneID)).
		Return(milestoneURL)
}

func TestMilestoneLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpMilestoneLink,
		tearDown,
		wikiConvert,
		"milestone:"+milestoneName,
		milestoneURL)
}

const (
	attachmentName        = "some-attachment.png"
	attachmentWikiRelPath = "attachments-dir/somepage/xyz.png"
	attachmentWikiURL     = "url-for-accessing-some-attachment-in-wiki"
)

func setUpWikiAttachmentLink(t *testing.T, page string) {
	setUp(t)

	// expect call to get relative path of attachment within wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiAttachmentRelPath(gomock.Eq(page), gomock.Eq(attachmentName)).
		Return(attachmentWikiRelPath)

	// expect call to lookup URL for attachment file
	mockGiteaAccessor.
		EXPECT().
		GetWikiFileURL(gomock.Eq(attachmentWikiRelPath)).
		Return(attachmentWikiURL)
}

func setUpImplicitWikiAttachmentLink(t *testing.T) {
	setUpWikiAttachmentLink(t, wikiPage)
}

func TestImplicitWikiAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpImplicitWikiAttachmentLink,
		tearDown,
		wikiConvert,
		"attachment:"+attachmentName,
		attachmentWikiURL)
}

const (
	otherWikiPage = "SomeOtherWikiPage"
)

func setUpExplicitWikiAttachmentLink(t *testing.T) {
	setUpWikiAttachmentLink(t, otherWikiPage)
}

func TestExplicitWikiAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpExplicitWikiAttachmentLink,
		tearDown,
		wikiConvert,
		"attachment:"+attachmentName+":wiki:"+otherWikiPage,
		attachmentWikiURL)
}

const (
	ticketAttachmentUUID = "UUID-of-ticket-attachment"
	ticketAttachmentURL  = "url-of-ticket-attachment"
)

func setUpTicketAttachmentLink(t *testing.T, tktID int64) {
	setUpAnyTicketLink(t, tktID)

	// expect call to get relative path of attachment within wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetIssueAttachmentUUID(gomock.Eq(issueID), gomock.Eq(attachmentName)).
		Return(ticketAttachmentUUID, nil)

	// expect call to lookup URL for attachment file
	mockGiteaAccessor.
		EXPECT().
		GetIssueAttachmentURL(gomock.Eq(issueID), gomock.Eq(ticketAttachmentUUID)).
		Return(ticketAttachmentURL)
}

func setUpImplicitTicketAttachmentLink(t *testing.T) {
	setUpTicketAttachmentLink(t, ticketID)
}

func TestImplicitTicketAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpImplicitTicketAttachmentLink,
		tearDown,
		ticketConvert,
		"attachment:"+attachmentName,
		ticketAttachmentURL)
}

func setUpExplicitTicketAttachmentLink(t *testing.T) {
	setUpTicketAttachmentLink(t, otherTicketID)
}

func TestExplicitTicketAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpExplicitTicketAttachmentLink,
		tearDown,
		ticketConvert,
		"attachment:"+attachmentName+":ticket:"+otherTicketIDStr,
		ticketAttachmentURL)
}

const (
	commitID  = "123abc456def7890"
	commitURL = "url-of-changeset-commit"
)

func setUpChangesetLink(t *testing.T) {
	setUp(t)

	// expect call to get commit URL
	mockGiteaAccessor.
		EXPECT().
		GetCommitURL(gomock.Eq(commitID)).
		Return(commitURL)
}

func TestChangesetLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpChangesetLink,
		tearDown,
		wikiConvert,
		"changeset:\""+commitID+"/repository-name\"",
		commitURL)
}

const (
	sourcePath = "path/to/some/source/file"
	sourceURL  = "url-of-source-file"
)

func setUpSourceLink(t *testing.T) {
	setUp(t)

	// expect call to get commit URL
	mockGiteaAccessor.
		EXPECT().
		GetSourceURL(gomock.Eq("master"), gomock.Eq(sourcePath)).
		Return(sourceURL)
}

func TestSourceLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpSourceLink,
		tearDown,
		wikiConvert,
		"source:\"repo-name/"+sourcePath+"\"",
		sourceURL)
}
