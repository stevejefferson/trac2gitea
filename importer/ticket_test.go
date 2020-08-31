// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

/*
 * The tests in here are (of necessity) complex.
 *
 * They are driven from data expessed in terms of the Ticket*Import types below set up in the `setUp*` functions.
 * The Ticket*Import contain both data to be returned from Trac and data to be returned from Gitea in response to that Trac data.
 *
 * The data created by the `setUp*` functions is passed into the `expect*` functions to set up expections on the mock trac and gitea accessors
 * and on the mock markdown converter governed by the fields in the data.
 *
 * The actual test methods call the top-level `setUp` method to set up all data followed by a selection of `expect` methods to create mock expections
 * for some or all of that data. Finally `importer.ImportTickets(...)` is called to trigger the actual importation we are testing.
 */

/*
 * Data set-up.
 */
// allocators - we give all items unique values so that we can spot any misallocations
func resetAllocators(startID int64) {
	idCounter = startID
	unixTimeCounter = 1000
	intCounter = 2000
}

var idCounter int64

func allocateID() int64 {
	idCounter++
	return idCounter
}

var unixTimeCounter int64

func allocateUnixTime() int64 {
	unixTimeCounter++
	return unixTimeCounter
}

var intCounter int64

func allocateInt() int64 {
	intCounter++
	return intCounter
}

// Trac -> Gitea naming maps
var (
	userMap map[string]string

	componentMap  map[string]string
	priorityMap   map[string]string
	resolutionMap map[string]string
	severityMap   map[string]string
	typeMap       map[string]string
	versionMap    map[string]string
)

func initMaps() {
	userMap = make(map[string]string)

	componentMap = make(map[string]string)
	priorityMap = make(map[string]string)
	resolutionMap = make(map[string]string)
	severityMap = make(map[string]string)
	typeMap = make(map[string]string)
	versionMap = make(map[string]string)
}

// TicketUserImport holds the data on a user referenced by an imported ticket
type TicketUserImport struct {
	tracUser    string
	giteaUser   string
	giteaUserID int64
}

func createTicketUserImport(prefix string, uMap map[string]string) *TicketUserImport {
	user := TicketUserImport{
		tracUser:    prefix + "-trac-user",
		giteaUser:   prefix + "-gitea-user",
		giteaUserID: allocateID(),
	}
	uMap[user.tracUser] = user.giteaUser
	return &user
}

var (
	closedTicketOwner    *TicketUserImport
	closedTicketReporter *TicketUserImport
	openTicketOwner      *TicketUserImport
	openTicketReporter   *TicketUserImport
)

func setUpTicketUsers(t *testing.T) {
	closedTicketOwner = createTicketUserImport("closed-owner", userMap)
	closedTicketReporter = createTicketUserImport("closed-reporter", userMap)
	openTicketOwner = createTicketUserImport("open-owner", userMap)
	openTicketReporter = createTicketUserImport("open-reporter", userMap)
}

var (
	closedTicketComment1Author *TicketUserImport
	closedTicketComment2Author *TicketUserImport
	openTicketComment1Author   *TicketUserImport
	openTicketComment2Author   *TicketUserImport
)

func setUpTicketCommentUsers(t *testing.T) {
	closedTicketComment1Author = createTicketUserImport("closed-ticket-comment1-author", userMap)
	closedTicketComment2Author = createTicketUserImport("closed-ticket-comment2-author", userMap)
	openTicketComment1Author = createTicketUserImport("open-ticket-comment1-author", userMap)
	openTicketComment2Author = createTicketUserImport("open-ticket-comment2-author", userMap)
}

var (
	closedTicketAttachment1Author *TicketUserImport
	closedTicketAttachment2Author *TicketUserImport
	openTicketAttachment1Author   *TicketUserImport
	openTicketAttachment2Author   *TicketUserImport
)

func setUpTicketAttachmentUsers(t *testing.T) {
	closedTicketAttachment1Author = createTicketUserImport("closed-ticket-attachment1-author", userMap)
	closedTicketAttachment2Author = createTicketUserImport("closed-ticket-attachment2-author", userMap)
	openTicketAttachment1Author = createTicketUserImport("open-ticket-attachment1-author", userMap)
	openTicketAttachment2Author = createTicketUserImport("open-ticket-attachment2-author", userMap)
}

// TicketLabelImport holds the data on a label associated with an imported ticket
type TicketLabelImport struct {
	name                  string
	labelName             string
	labelID               int64
	issueLabelID          int64
	giteaLabelExists      bool
	giteaIssueLabelExists bool
}

func createTicketLabelImport(prefix string, giteaLabelExists bool, giteaIssueLabelExists bool, ticketLabelMap map[string]string) *TicketLabelImport {
	ticketLabel := TicketLabelImport{
		name:                  prefix + "-name",
		labelName:             prefix + "-label",
		labelID:               allocateID(),
		issueLabelID:          allocateID(),
		giteaLabelExists:      giteaLabelExists,
		giteaIssueLabelExists: giteaIssueLabelExists,
	}

	ticketLabelMap[ticketLabel.name] = ticketLabel.labelName
	return &ticketLabel
}

var (
	closedTicketComponent  *TicketLabelImport
	closedTicketPriority   *TicketLabelImport
	closedTicketResolution *TicketLabelImport
	closedTicketSeverity   *TicketLabelImport
	closedTicketType       *TicketLabelImport
	closedTicketVersion    *TicketLabelImport

	openTicketComponent  *TicketLabelImport
	openTicketPriority   *TicketLabelImport
	openTicketResolution *TicketLabelImport
	openTicketSeverity   *TicketLabelImport
	openTicketType       *TicketLabelImport
	openTicketVersion    *TicketLabelImport
)

func setUpTicketLabels(t *testing.T) {
	closedTicketComponent = createTicketLabelImport("closed-component", false, false, componentMap)
	closedTicketPriority = createTicketLabelImport("closed-priority", true, false, priorityMap)
	closedTicketResolution = createTicketLabelImport("closed-resolution", true, true, resolutionMap)
	closedTicketSeverity = createTicketLabelImport("closed-severity", false, false, severityMap)
	closedTicketType = createTicketLabelImport("closed-type", false, false, typeMap)
	closedTicketVersion = createTicketLabelImport("closed-version", true, false, versionMap)

	openTicketComponent = createTicketLabelImport("open-component", true, true, componentMap)
	openTicketPriority = createTicketLabelImport("open-priority", false, false, priorityMap)
	openTicketResolution = createTicketLabelImport("open-resolution", false, false, resolutionMap)
	openTicketSeverity = createTicketLabelImport("open-severity", true, false, severityMap)
	openTicketType = createTicketLabelImport("open-type", true, true, typeMap)
	openTicketVersion = createTicketLabelImport("open-version", false, false, versionMap)
}

// TicketCommentImport holds the data on a ticket comment import operation
type TicketCommentImport struct {
	issueCommentID          int64
	author                  *TicketUserImport
	text                    string
	markdownText            string
	time                    int64
	giteaIssueCommentExists bool
}

func createTicketCommentImport(prefix string, author *TicketUserImport, giteaIssueCommentExists bool) *TicketCommentImport {
	return &TicketCommentImport{
		issueCommentID:          allocateID(),
		author:                  author,
		text:                    prefix + " ticket comment text",
		markdownText:            prefix + " ticket comment text after conversion to markdown",
		time:                    allocateUnixTime(),
		giteaIssueCommentExists: giteaIssueCommentExists,
	}
}

func createTracTicketComment(ticket *TicketImport, ticketComment *TicketCommentImport) *trac.TicketComment {
	return &trac.TicketComment{
		TicketID: ticket.ticketID,
		Time:     ticketComment.time,
		Author:   ticketComment.author.tracUser,
		Text:     ticketComment.text,
	}
}

var (
	closedTicketComment1 *TicketCommentImport
	closedTicketComment2 *TicketCommentImport
	openTicketComment1   *TicketCommentImport
	openTicketComment2   *TicketCommentImport
)

func setUpTicketComments(t *testing.T) {
	setUpTicketCommentUsers(t)
	closedTicketComment1 = createTicketCommentImport("closed-ticket-comment1", closedTicketComment1Author, false)
	closedTicketComment2 = createTicketCommentImport("closed-ticket-comment2", closedTicketComment2Author, false)

	openTicketComment1 = createTicketCommentImport("open-ticket-comment1", openTicketComment1Author, true)
	openTicketComment2 = createTicketCommentImport("open-ticket-comment2", openTicketComment2Author, false)
}

// TicketAttachmentImport holds data on a ticket attachment import operation
type TicketAttachmentImport struct {
	issueAttachmentID int64
	comment           *TicketCommentImport
	filename          string
	attachmentPath    string
	size              int64
}

func createTicketAttachmentImport(prefix string, author *TicketUserImport) *TicketAttachmentImport {
	// express part of attachment data in terms of the comment that will appear in Gitea to describe it
	comment := createTicketCommentImport(prefix+"-comment-", author, false)

	// trac attachment path must have final directory of at least 12 chars (the trac UUID)
	attachmentFile := prefix + "-attachment.file"
	attachmentPath := "/path/to/attachment/" + prefix + "123456789012/" + attachmentFile

	return &TicketAttachmentImport{
		issueAttachmentID: allocateID(),
		comment:           comment,
		filename:          attachmentFile,
		attachmentPath:    attachmentPath,
		size:              allocateInt(),
	}
}

func createTracTicketAttachment(ticket *TicketImport, ticketAttachment *TicketAttachmentImport) *trac.TicketAttachment {
	return &trac.TicketAttachment{
		TicketID:    ticket.ticketID,
		Size:        ticketAttachment.size,
		Author:      ticketAttachment.comment.author.tracUser,
		FileName:    ticketAttachment.filename,
		Description: ticketAttachment.comment.text,
		Time:        ticketAttachment.comment.time,
	}
}

var (
	closedTicketAttachment1 *TicketAttachmentImport
	closedTicketAttachment2 *TicketAttachmentImport
	openTicketAttachment1   *TicketAttachmentImport
	openTicketAttachment2   *TicketAttachmentImport
)

func setUpTicketAttachments(t *testing.T) {
	setUpTicketAttachmentUsers(t)
	closedTicketAttachment1 = createTicketAttachmentImport("closed-ticket-attachment1", closedTicketAttachment1Author)
	closedTicketAttachment2 = createTicketAttachmentImport("closed-ticket-attachment2", closedTicketAttachment2Author)
	openTicketAttachment1 = createTicketAttachmentImport("open-ticket-attachment1", openTicketAttachment1Author)
	openTicketAttachment2 = createTicketAttachmentImport("open-ticket-attachment2", openTicketAttachment2Author)
}

// TicketImport holds the data on a ticket import operation
type TicketImport struct {
	ticketID            int64
	issueID             int64
	summary             string
	description         string
	descriptionMarkdown string
	owner               *TicketUserImport
	reporter            *TicketUserImport
	milestoneName       string
	componentLabel      *TicketLabelImport
	priorityLabel       *TicketLabelImport
	resolutionLabel     *TicketLabelImport
	severityLabel       *TicketLabelImport
	typeLabel           *TicketLabelImport
	versionLabel        *TicketLabelImport
	closed              bool
	status              string
	created             int64
}

func createTicketImport(
	prefix string,
	closed bool,
	owner *TicketUserImport,
	reporter *TicketUserImport,
	componentLabel *TicketLabelImport,
	priorityLabel *TicketLabelImport,
	resolutionLabel *TicketLabelImport,
	severityLabel *TicketLabelImport,
	typeLabel *TicketLabelImport,
	versionLabel *TicketLabelImport) *TicketImport {
	status := "open"
	if closed {
		status = "closed"
	}

	return &TicketImport{
		ticketID:            allocateID(),
		issueID:             allocateID(),
		summary:             prefix + "-summary",
		description:         prefix + "-description",
		descriptionMarkdown: prefix + "-markdown",
		owner:               owner,
		reporter:            reporter,
		milestoneName:       prefix + "-milestone",
		componentLabel:      componentLabel,
		priorityLabel:       priorityLabel,
		resolutionLabel:     resolutionLabel,
		severityLabel:       severityLabel,
		typeLabel:           typeLabel,
		versionLabel:        versionLabel,
		closed:              closed,
		status:              status,
		created:             allocateUnixTime(),
	}
}

func createTracTicket(ticket *TicketImport) *trac.Ticket {
	return &trac.Ticket{
		TicketID:       ticket.ticketID,
		Summary:        ticket.summary,
		Description:    ticket.description,
		Owner:          ticket.owner.tracUser,
		Reporter:       ticket.reporter.tracUser,
		MilestoneName:  ticket.milestoneName,
		ComponentName:  ticket.componentLabel.name,
		PriorityName:   ticket.priorityLabel.name,
		ResolutionName: ticket.resolutionLabel.name,
		SeverityName:   ticket.severityLabel.name,
		TypeName:       ticket.typeLabel.name,
		VersionName:    ticket.versionLabel.name,
		Status:         ticket.status,
		Created:        ticket.created,
	}
}

var (
	closedTicket *TicketImport
	openTicket   *TicketImport
)

// setUpTickets is the top-level setUp method for the ticket tests.
// It should be called by all tests - it is the mock expectations that determines which parts of the set up data are actually used in any test
func setUpTickets(t *testing.T) {
	setUp(t)
	resetAllocators(1000)
	initMaps()
	setUpTicketUsers(t)
	setUpTicketLabels(t)
	setUpTicketComments(t)
	setUpTicketAttachments(t)

	closedTicket = createTicketImport(
		"closed", true,
		closedTicketOwner, closedTicketReporter,
		closedTicketComponent, closedTicketPriority, closedTicketResolution, closedTicketSeverity, closedTicketType, closedTicketVersion)
	openTicket = createTicketImport(
		"open", false,
		openTicketOwner, openTicketReporter,
		openTicketComponent, openTicketPriority, openTicketResolution, openTicketSeverity, openTicketType, openTicketVersion)
}

/*
 * Mock expectations
 */
func expectTracCommentRetrievals(t *testing.T, ticket *TicketImport, ticketComments ...*TicketCommentImport) {
	// expect trac accessor to return each of our trac ticket comments
	mockTracAccessor.
		EXPECT().
		GetTicketComments(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, handlerFn func(comment *trac.TicketComment) error) error {
			for _, ticketComment := range ticketComments {
				tracComment := createTracTicketComment(ticket, ticketComment)
				handlerFn(tracComment)
			}
			return nil
		})
}

func expectTracAttachmentRetrievals(t *testing.T, ticket *TicketImport, ticketAttachments ...*TicketAttachmentImport) {
	// expect trac accessor to return each of our trac ticket attachments
	mockTracAccessor.
		EXPECT().
		GetTicketAttachments(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, handlerFn func(attachment *trac.TicketAttachment) error) error {
			for _, ticketAttachment := range ticketAttachments {
				tracAttachment := createTracTicketAttachment(ticket, ticketAttachment)
				handlerFn(tracAttachment)
			}
			return nil
		})
}

func expectTracTicketRetrievals(t *testing.T, tickets ...*TicketImport) {
	// expect trac accessor to return each of our trac tickets
	mockTracAccessor.
		EXPECT().
		GetTickets(gomock.Any()).
		DoAndReturn(func(handlerFn func(ticket *trac.Ticket) error) error {
			for _, ticket := range tickets {
				tracTicket := createTracTicket(ticket)
				handlerFn(tracTicket)
			}
			return nil
		})
}

func expectUserLookup(t *testing.T, user *TicketUserImport) {
	mockGiteaAccessor.
		EXPECT().
		GetUserID(gomock.Eq(user.giteaUser)).
		Return(user.giteaUserID, nil)
}

func expectIssueLookup(t *testing.T, ticket *TicketImport) {
	// expect to look for existing Gitea issue for Trac ticket - we currently assume there is none
	mockGiteaAccessor.
		EXPECT().
		GetIssueID(gomock.Eq(ticket.ticketID)).
		Return(int64(-1), nil)
}

func expectLabelRetrieval(t *testing.T, label *TicketLabelImport) {
	// expect to lookup label by name - return -1 if we expect to create label
	returnedLabelID := label.labelID
	if !label.giteaLabelExists {
		returnedLabelID = int64(-1)
	}
	mockGiteaAccessor.
		EXPECT().
		GetLabelID(gomock.Eq(label.labelName)).
		Return(returnedLabelID, nil)

	if !label.giteaLabelExists {
		mockGiteaAccessor.
			EXPECT().
			AddLabel(gomock.Eq(label.labelName), gomock.Any()).
			Return(label.labelID, nil)
	}
}

func expectDescriptionMarkdownConversion(t *testing.T, ticket *TicketImport) {
	mockMarkdownConverter.
		EXPECT().
		TicketConvert(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, text string) string {
			assertTrue(t, strings.Contains(text, ticket.description))
			return ticket.descriptionMarkdown
		})
}

func expectIssueCreation(t *testing.T, ticket *TicketImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssue(gomock.Any()).
		DoAndReturn(func(issue *gitea.Issue) (int64, error) {
			assertEquals(t, issue.Index, ticket.ticketID)
			assertEquals(t, issue.Summary, ticket.summary)
			assertTrue(t, strings.Contains(issue.Description, ticket.descriptionMarkdown))
			assertEquals(t, issue.OwnerID, ticket.owner.giteaUserID)
			assertEquals(t, issue.Owner, ticket.owner.giteaUser)
			assertEquals(t, issue.ReporterID, ticket.reporter.giteaUserID)
			assertEquals(t, issue.Milestone, ticket.milestoneName)
			assertEquals(t, issue.Closed, ticket.closed)
			assertEquals(t, issue.Created, ticket.created)
			return ticket.issueID, nil
		})
}

func expectIssueLabelRetrieval(t *testing.T, ticket *TicketImport, ticketLabel *TicketLabelImport) {
	// expect retrieval/creation of underlying label first
	expectLabelRetrieval(t, ticketLabel)

	// expect to lookup issue label by id - return -1 if we expect to create issue label
	returnedIssueLabelID := ticketLabel.issueLabelID
	if !ticketLabel.giteaIssueLabelExists {
		returnedIssueLabelID = int64(-1)
	}
	mockGiteaAccessor.
		EXPECT().
		GetIssueLabelID(gomock.Eq(ticket.issueID), gomock.Eq(ticketLabel.labelID)).
		Return(returnedIssueLabelID, nil)

	if !ticketLabel.giteaIssueLabelExists {
		mockGiteaAccessor.
			EXPECT().
			AddIssueLabel(gomock.Eq(ticket.issueID), gomock.Eq(ticketLabel.labelID)).
			Return(ticketLabel.issueLabelID, nil)
	}
}

func expectIssueCommentRetrieval(t *testing.T, ticket *TicketImport, ticketComment *TicketCommentImport) {
	// expect to lookup issue comment by id - return -1 if we expect to create issue comment
	returnedIssueCommentID := ticketComment.issueCommentID
	if !ticketComment.giteaIssueCommentExists {
		returnedIssueCommentID = int64(-1)
	}

	mockGiteaAccessor.
		EXPECT().
		GetIssueCommentID(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, text string) (int64, error) {
			assertEquals(t, issueID, ticket.issueID)
			assertTrue(t, strings.Contains(text, ticketComment.markdownText))
			return returnedIssueCommentID, nil
		})

	if !ticketComment.giteaIssueCommentExists {
		mockGiteaAccessor.
			EXPECT().
			AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
			DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
				assertEquals(t, issueComment.AuthorID, ticketComment.author.giteaUserID)
				assertTrue(t, strings.Contains(issueComment.Text, ticketComment.markdownText))
				assertEquals(t, issueComment.Time, ticketComment.time)
				return ticketComment.issueCommentID, nil
			})
	}
}

func expectTicketCommentMarkdownConversion(t *testing.T, ticket *TicketImport, ticketComment *TicketCommentImport) {
	mockMarkdownConverter.
		EXPECT().
		TicketConvert(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, text string) string {
			assertTrue(t, strings.Contains(text, ticketComment.text))
			return ticketComment.markdownText
		})
}

func expectAllTicketCommentActions(t *testing.T, ticket *TicketImport, ticketComment *TicketCommentImport) {
	// expect to lookup Gitea equivalents of Trac ticket comment author
	expectUserLookup(t, ticketComment.author)

	// expect to convert ticket comment text to markdown
	expectTicketCommentMarkdownConversion(t, ticket, ticketComment)

	// expect retrieval/creation of issue comment for ticket comment
	expectIssueCommentRetrieval(t, ticket, ticketComment)
}

func expectTracAttachmentPathRetrieval(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	mockTracAccessor.
		EXPECT().
		GetTicketAttachmentPath(gomock.Any()).
		DoAndReturn(func(tracAttachment *trac.TicketAttachment) string {
			assertEquals(t, tracAttachment.TicketID, ticket.ticketID)
			assertEquals(t, tracAttachment.FileName, ticketAttachment.filename)
			assertEquals(t, tracAttachment.Size, ticketAttachment.size)
			assertEquals(t, tracAttachment.Time, ticketAttachment.comment.time)
			assertEquals(t, tracAttachment.Author, ticketAttachment.comment.author.tracUser)
			assertTrue(t, strings.Contains(tracAttachment.Description, ticketAttachment.comment.text))
			return ticketAttachment.attachmentPath
		})
}

func expectAttachmentUUIDRetrieval(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	mockGiteaAccessor.
		EXPECT().
		GetIssueAttachmentUUID(gomock.Eq(ticket.issueID), gomock.Eq(ticketAttachment.filename)).
		Return("", nil)
}

func expectIssueAttachmentCreation(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueAttachment(gomock.Eq(ticket.issueID), gomock.Eq(ticketAttachment.filename), gomock.Any()).
		DoAndReturn(func(issueID int64, filename string, issueAttachment *gitea.IssueAttachment) (int64, error) {
			assertEquals(t, issueAttachment.CommentID, ticketAttachment.comment.issueCommentID)
			assertEquals(t, issueAttachment.FilePath, ticketAttachment.attachmentPath)
			assertEquals(t, issueAttachment.Time, ticketAttachment.comment.time)
			return ticketAttachment.issueAttachmentID, nil
		})
}

func expectAllTicketAttachmentActions(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	expectAllTicketCommentActions(t, ticket, ticketAttachment.comment)
	expectTracAttachmentPathRetrieval(t, ticket, ticketAttachment)
	expectAttachmentUUIDRetrieval(t, ticket, ticketAttachment)
	expectIssueAttachmentCreation(t, ticket, ticketAttachment)
}

func expectIssueUpdateTimeSetToLatestOf(t *testing.T, ticket *TicketImport, ticketComments ...*TicketCommentImport) {
	latestUpdateTime := ticket.created
	for _, ticketComment := range ticketComments {
		if ticketComment.time > latestUpdateTime {
			latestUpdateTime = ticketComment.time
		}
	}

	mockGiteaAccessor.
		EXPECT().
		SetIssueUpdateTime(gomock.Eq(ticket.issueID), gomock.Eq(latestUpdateTime)).
		Return(nil)
}

func expectRepoIssueCountUpdate(t *testing.T, numIssues int, numClosedIssues int) {
	mockGiteaAccessor.
		EXPECT().
		UpdateRepoIssueCount(gomock.Eq(numIssues), gomock.Eq(numClosedIssues)).
		Return(nil)
}

func expectAllTicketActions(t *testing.T, ticket *TicketImport) {
	// expect to check whether Gitea issues already exist for our tickets
	expectIssueLookup(t, ticket)

	// expect to lookup Gitea equivalents of Trac ticket owner and reporter
	expectUserLookup(t, ticket.owner)
	expectUserLookup(t, ticket.reporter)

	// expect to convert ticket description to markdown
	expectDescriptionMarkdownConversion(t, ticket)

	// expect to create Gitea issue
	expectIssueCreation(t, ticket)

	// expect retrieval/creation of all labels from Trac ticket appearing in the Gitea issue
	expectIssueLabelRetrieval(t, ticket, ticket.componentLabel)
	expectIssueLabelRetrieval(t, ticket, ticket.priorityLabel)
	expectIssueLabelRetrieval(t, ticket, ticket.resolutionLabel)
	expectIssueLabelRetrieval(t, ticket, ticket.severityLabel)
	expectIssueLabelRetrieval(t, ticket, ticket.typeLabel)
	expectIssueLabelRetrieval(t, ticket, ticket.versionLabel)
}

func TestImportClosedTicketOnly(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, closedTicket)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, closedTicket)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportOpenTicketOnly(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, openTicket)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportMultipleTicketsOnly(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of tickets from Trac
	expectTracTicketRetrievals(t, closedTicket, openTicket)

	// expect all actions for creating Gitea issue from Trac tickets
	expectAllTicketActions(t, closedTicket)
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, closedTicket)
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, closedTicket)
	expectTracCommentRetrievals(t, openTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket)
	expectIssueUpdateTimeSetToLatestOf(t, openTicket)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 2, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithAttachments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us attachments
	expectTracAttachmentRetrievals(t, closedTicket, closedTicketAttachment1, closedTicketAttachment2)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue attachments from Trac ticket attachments
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment1)
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment2)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketAttachment1.comment, closedTicketAttachment2.comment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportMultipleTicketsWithAttachments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket, openTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, closedTicket)
	expectAllTicketActions(t, openTicket)

	// expect trac to return us attachments
	expectTracAttachmentRetrievals(t, closedTicket, closedTicketAttachment1, closedTicketAttachment2)
	expectTracAttachmentRetrievals(t, openTicket, openTicketAttachment1, openTicketAttachment2)

	// expect all actions for creating Gitea issue attachments from Trac ticket attachments
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment1)
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment2)
	expectAllTicketAttachmentActions(t, openTicket, openTicketAttachment1)
	expectAllTicketAttachmentActions(t, openTicket, openTicketAttachment2)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, closedTicket)
	expectTracCommentRetrievals(t, openTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketAttachment1.comment, closedTicketAttachment2.comment)
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, openTicketAttachment1.comment, openTicketAttachment2.comment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 2, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithComments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, closedTicket)

	// expect trac to return us comments
	expectTracCommentRetrievals(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment1)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment2)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportMultipleTicketsWithComments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of tickets from Trac
	expectTracTicketRetrievals(t, openTicket, closedTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, openTicket)
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)
	expectTracAttachmentRetrievals(t, closedTicket)

	// expect trac to return us comments
	expectTracCommentRetrievals(t, openTicket, openTicketComment1, openTicketComment2)
	expectTracCommentRetrievals(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, openTicket, openTicketComment1)
	expectAllTicketCommentActions(t, openTicket, openTicketComment2)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment1)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment2)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, openTicketComment1, openTicketComment2)
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 2, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithAttachmentsAndComments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us attachments
	expectTracAttachmentRetrievals(t, openTicket, openTicketAttachment1, openTicketAttachment2)

	// expect all actions for creating Gitea issue attachments from Trac ticket attachments
	expectAllTicketAttachmentActions(t, openTicket, openTicketAttachment1)
	expectAllTicketAttachmentActions(t, openTicket, openTicketAttachment2)

	// expect trac to return us comments
	expectTracCommentRetrievals(t, openTicket, openTicketComment1, openTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, openTicket, openTicketComment1)
	expectAllTicketCommentActions(t, openTicket, openTicketComment2)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket,
		openTicketComment1, openTicketComment2, openTicketAttachment1.comment, openTicketAttachment2.comment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportMultipleTicketsWithAttachmentsAndComments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of tickets from Trac
	expectTracTicketRetrievals(t, openTicket, closedTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, openTicket)
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us attachments
	expectTracAttachmentRetrievals(t, openTicket, openTicketAttachment1, openTicketAttachment2)
	expectTracAttachmentRetrievals(t, closedTicket, closedTicketAttachment1, closedTicketAttachment2)

	// expect all actions for creating Gitea issue attachments from Trac ticket attachments
	expectAllTicketAttachmentActions(t, openTicket, openTicketAttachment1)
	expectAllTicketAttachmentActions(t, openTicket, openTicketAttachment2)
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment1)
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment2)

	// expect trac to return us comments
	expectTracCommentRetrievals(t, openTicket, openTicketComment1, openTicketComment2)
	expectTracCommentRetrievals(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, openTicket, openTicketComment1)
	expectAllTicketCommentActions(t, openTicket, openTicketComment2)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment1)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment2)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket,
		openTicketComment1, openTicketComment2, openTicketAttachment1.comment, openTicketAttachment2.comment)
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket,
		closedTicketComment1, closedTicketComment2, closedTicketAttachment1.comment, closedTicketAttachment2.comment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 2, 1)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
