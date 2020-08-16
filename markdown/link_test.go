package markdown_test

import (
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

func tracImageLink(link string) string {
	return "[[Image(" + link + ")]]"
}

func markdownAutomaticLink(link string) string {
	return "<" + link + ">"
}

func markdownLinkWithText(link string, text string) string {
	return "[" + text + "](" + link + ")"
}

func markdownImageLink(link string) string {
	return "![](" + link + ")"
}

// verifyLink verifies that the provided trac formatting for a link + text results in the corresponding markdown format
func verifyLink(t *testing.T, setUpFn func(t *testing.T), tearDownFn func(t *testing.T), tracFormatLink string, markdownFormatLink string) {
	setUpFn(t)
	defer tearDownFn(t)

	conversion := converter.Convert(leadingText + " " + tracFormatLink + " " + trailingText)
	assertEquals(t, conversion, leadingText+" "+markdownFormatLink+" "+trailingText)
}

const linkText = "text associated with link"

func verifyAllLinkTypes(t *testing.T, setUpFn func(t *testing.T), tearDownFn func(t *testing.T), tracLinkStr string, markdownLinkStr string) {
	verifyLink(t, setUpFn, tearDownFn, tracPlainLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracSingleBracketLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracSingleBracketLinkWithText(tracLinkStr, linkText), markdownLinkWithText(markdownLinkStr, linkText))
	verifyLink(t, setUpFn, tearDownFn, tracDoubleBracketLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracDoubleBracketLinkWithText(tracLinkStr, linkText), markdownLinkWithText(markdownLinkStr, linkText))
	verifyLink(t, setUpFn, tearDownFn, tracImageLink(tracLinkStr), markdownImageLink(markdownLinkStr))
}

const httpLink = "http://www.example.com"

func TestHttpLinks(t *testing.T) {
	verifyAllLinkTypes(t, setUp, tearDown, httpLink, httpLink)
}

const httpsLink = "https://www.example.com"

func TestHttpsLink(t *testing.T) {
	verifyAllLinkTypes(t, setUp, tearDown, httpsLink, httpsLink)
}

const (
	htdocsFile        = "somefile.jpg"
	tracHtdocsLink    = "htdocs:" + htdocsFile
	markdownHtdocsURL = "./somedir/" + htdocsFile
)

func setUpHtdocs(t *testing.T) {
	setUp(t)

	// expect to have to retrieve full path to htdocs file within Trac workspace
	tracHtdocsPath := filepath.Join("dir", "trac", "htdocs", htdocsFile)
	mockTracAccessor.
		EXPECT().
		GetFullPath(gomock.Eq("htdocs"), gomock.Eq(htdocsFile)).
		Return(tracHtdocsPath).
		AnyTimes()

	// expect to copy file into "htdocs" subdirectory of Wiki repo
	giteaWikiHtdocsPath := filepath.Join("htdocs", htdocsFile)
	mockGiteaWikiAccessor.
		EXPECT().
		CopyFile(gomock.Eq(tracHtdocsPath), gomock.Eq(giteaWikiHtdocsPath)).
		Return(markdownHtdocsURL).
		AnyTimes()
}

func TestHtdocsLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpHtdocs,
		tearDown,
		tracHtdocsLink,
		markdownHtdocsURL)
}

const (
	wikiPageName            = "SomeWikiPage"
	transformedWikiPageName = "TransformedWikiPage"
)

func setUpWiki(t *testing.T) {
	setUp(t)

	// expect call to translate name of wiki page
	mockGiteaWikiAccessor.
		EXPECT().
		TranslatePageName(gomock.Eq(wikiPageName)).
		Return(transformedWikiPageName).
		AnyTimes()
}

func TestWikiUnprefixedLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWiki,
		tearDown,
		wikiPageName,
		transformedWikiPageName)
}

func TestWikiPrefixedLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWiki,
		tearDown,
		"wiki:"+wikiPageName,
		transformedWikiPageName)
}

const (
	ticketID    int64 = 314159
	ticketIDStr       = "314159"
	issueID     int64 = 26535
)

func setUpTicket(t *testing.T) {
	setUp(t)

	// expect call to lookup gitea issue for trac ticket
	mockGiteaAccessor.
		EXPECT().
		GetIssueID(gomock.Eq(ticketID)).
		Return(issueID).
		AnyTimes()
}

func TestTicketLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpTicket,
		tearDown,
		"ticket:"+ticketIDStr,
		"#"+ticketIDStr)
}

const (
	tracCommentNum    int64  = 12
	tracCommentNumStr        = "12"
	commentStr               = "this is the text of a comment"
	commentID         int64  = 54321
	commentURL        string = "url-of-comment-54321"
)

func setUpTicketComment(t *testing.T) {
	setUpTicket(t)

	// expect a call to lookup text of trac comment
	mockTracAccessor.
		EXPECT().
		GetCommentString(gomock.Eq(issueID), gomock.Eq(tracCommentNum)).
		Return(commentStr).
		AnyTimes()

	// expect call to lookup gitea ID for trac comment
	mockGiteaAccessor.
		EXPECT().
		GetCommentID(gomock.Eq(issueID), gomock.Eq(commentStr)).
		Return(commentID).
		AnyTimes()

	// expect call to lookup URL of gitea comment
	mockGiteaAccessor.
		EXPECT().
		GetCommentURL(gomock.Eq(issueID), gomock.Eq(commentID)).
		Return(commentURL).
		AnyTimes()
}

func TestTicketCommentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpTicketComment,
		tearDown,
		"comment:"+tracCommentNumStr+":ticket:"+ticketIDStr,
		commentURL)
}
