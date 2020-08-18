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

// verifyLink verifies that the provided trac formatting for a link + text results in the corresponding markdown format
func verifyLink(t *testing.T, setUpFn func(t *testing.T), tearDownFn func(t *testing.T), tracFormatLink string, markdownFormatLink string) {
	setUpFn(t)
	defer tearDownFn(t)

	conversion := converter.Convert(leadingText + " " + tracFormatLink + " " + trailingText)
	assertEquals(t, conversion, leadingText+" "+markdownFormatLink+" "+trailingText)
}

const (
	linkText            = "text associated with link"
	additionalImageLink = "http://somewhere.com"
)

func verifyAllLinkTypes(t *testing.T, setUpFn func(t *testing.T), tearDownFn func(t *testing.T), tracLinkStr string, markdownLinkStr string) {
	verifyLink(t, setUpFn, tearDownFn, tracPlainLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracSingleBracketLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracSingleBracketLinkWithText(tracLinkStr, linkText), markdownLinkWithText(markdownLinkStr, linkText))
	verifyLink(t, setUpFn, tearDownFn, tracDoubleBracketLink(tracLinkStr), markdownAutomaticLink(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracDoubleBracketLinkWithText(tracLinkStr, linkText), markdownLinkWithText(markdownLinkStr, linkText))
	verifyLink(t, setUpFn, tearDownFn, tracImage(tracLinkStr), markdownImage(markdownLinkStr))
	verifyLink(t, setUpFn, tearDownFn, tracImageWithLink(tracLinkStr, additionalImageLink), markdownImageWithLink(markdownLinkStr, additionalImageLink))
	verifyLink(t, setUpFn, tearDownFn, tracImageWithLink(additionalImageLink, tracLinkStr), markdownImageWithLink(additionalImageLink, markdownLinkStr))
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
		Return(tracHtdocPath).
		AnyTimes()

	// expect to retrieve path where trac "htdocs" file will be stored in the Wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiHtdocRelPath(gomock.Eq(htdocFile)).
		Return(giteaHtdocFile).
		AnyTimes()

	// expect to copy file into "htdocs" subdirectory of Wiki repo
	mockGiteaAccessor.
		EXPECT().
		CopyFileToWiki(gomock.Eq(tracHtdocPath), gomock.Eq(giteaHtdocFile)).
		AnyTimes()

	// expect to retrieve URL for viewing htdocs file in wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiFileURL(gomock.Eq(giteaHtdocFile)).
		Return(giteaHtdocURL).
		AnyTimes()
}

func TestHtdocsLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpHtdocs,
		tearDown,
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
		Return(transformedWikiPageName).
		AnyTimes()
}

func TestWikiUnprefixedLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		wikiPageName,
		transformedWikiPageName)
}

func TestWikiPrefixedLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		"wiki:"+wikiPageName,
		transformedWikiPageName)
}

const wikiPageAnchor = "page-anchor"

func TestWikiUnprefixedLinkWithAnchor(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		wikiPageName+"#"+wikiPageAnchor,
		transformedWikiPageName+"#"+wikiPageAnchor)
}

func TestWikiPrefixedLinkWithAnchor(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiLink,
		tearDown,
		"wiki:"+wikiPageName+"#"+wikiPageAnchor,
		transformedWikiPageName+"#"+wikiPageAnchor)
}

const (
	ticketID    int64 = 314159
	ticketIDStr       = "314159"
	issueID     int64 = 26535
	issueURL          = "url-for-viewing-issue=26535"
)

func setUpAnyTicketLink(t *testing.T) {
	setUp(t)

	// expect call to lookup gitea issue for trac ticket
	mockGiteaAccessor.
		EXPECT().
		GetIssueID(gomock.Eq(ticketID)).
		Return(issueID).
		AnyTimes()
}

func setUpTicketOnlyLink(t *testing.T) {
	setUpAnyTicketLink(t)

	// expect call to lookup gitea issue URL
	mockGiteaAccessor.
		EXPECT().
		GetIssueURL(gomock.Eq(issueID)).
		Return(issueURL).
		AnyTimes()
}

func TestTicketLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpTicketOnlyLink,
		tearDown,
		"ticket:"+ticketIDStr,
		issueURL)
}

const (
	tracCommentNum    int64  = 12
	tracCommentNumStr        = "12"
	commentStr               = "this is the text of a comment"
	commentID         int64  = 54321
	commentURL        string = "url-of-comment-54321"
)

func setUpTicketCommentLink(t *testing.T) {
	setUpAnyTicketLink(t)

	// expect a call to lookup text of trac comment
	mockTracAccessor.
		EXPECT().
		GetCommentString(gomock.Eq(ticketID), gomock.Eq(tracCommentNum)).
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
		setUpMilestoneLink,
		tearDown,
		"milestone:"+milestoneName,
		milestoneURL)
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
		Return(milestoneID).
		AnyTimes()

	// expect call to lookup URL for milestone
	mockGiteaAccessor.
		EXPECT().
		GetMilestoneURL(gomock.Eq(milestoneID)).
		Return(milestoneURL).
		AnyTimes()
}

func TestMilestoneLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpTicketCommentLink,
		tearDown,
		"comment:"+tracCommentNumStr+":ticket:"+ticketIDStr,
		commentURL)
}

const (
	attachmentName        = "some-attachment.png"
	attachmentWikiRelPath = "attachments-dir/somepage/xyz.png"
	attachmentWikiURL     = "url-for-accessing-some-attachment-in-wiki"
)

func setUpAttachmentLink(t *testing.T) {
	setUp(t)

	// expect call to get relative path of attachment within wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiAttachmentRelPath(gomock.Eq(wikiPage), gomock.Eq(attachmentName)).
		Return(attachmentWikiRelPath).
		AnyTimes()

	// expect call to lookup URL for attachment file
	mockGiteaAccessor.
		EXPECT().
		GetWikiFileURL(gomock.Eq(attachmentWikiRelPath)).
		Return(attachmentWikiURL).
		AnyTimes()
}

func TestAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpAttachmentLink,
		tearDown,
		"attachment:"+attachmentName,
		attachmentWikiURL)
}

const (
	otherWikiPage = "SomeOtherWikiPage"
)

func setUpWikiAttachmentLink(t *testing.T) {
	setUp(t)

	// expect call to get relative path of attachment within wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetWikiAttachmentRelPath(gomock.Eq(otherWikiPage), gomock.Eq(attachmentName)).
		Return(attachmentWikiRelPath).
		AnyTimes()

	// expect call to lookup URL for attachment file
	mockGiteaAccessor.
		EXPECT().
		GetWikiFileURL(gomock.Eq(attachmentWikiRelPath)).
		Return(attachmentWikiURL).
		AnyTimes()
}

func TestWikiAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpWikiAttachmentLink,
		tearDown,
		"attachment:"+attachmentName+":wiki:"+otherWikiPage,
		attachmentWikiURL)
}

const (
	ticketAttachmentUUID = "UUID-of-ticket-attachment"
	ticketAttachmentURL  = "url-of-ticket-attachment"
)

func setUpTicketAttachmentLink(t *testing.T) {
	setUpAnyTicketLink(t)

	// expect call to get relative path of attachment within wiki repo
	mockGiteaAccessor.
		EXPECT().
		GetAttachmentUUID(gomock.Eq(issueID), gomock.Eq(attachmentName)).
		Return(ticketAttachmentUUID).
		AnyTimes()

	// expect call to lookup URL for attachment file
	mockGiteaAccessor.
		EXPECT().
		GetAttachmentURL(gomock.Eq(ticketAttachmentUUID)).
		Return(ticketAttachmentURL).
		AnyTimes()
}

func TestTicketAttachmentLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpTicketAttachmentLink,
		tearDown,
		"attachment:"+attachmentName+":ticket:"+ticketIDStr,
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
		Return(commitURL).
		AnyTimes()
}

func TestChangesetLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpChangesetLink,
		tearDown,
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
		Return(sourceURL).
		AnyTimes()
}

func TestSourceLink(t *testing.T) {
	verifyAllLinkTypes(
		t,
		setUpSourceLink,
		tearDown,
		"source:\"repo-name/"+sourcePath+"\"",
		sourceURL)
}
