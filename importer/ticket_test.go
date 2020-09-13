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
 * The tests in here are complex.
 *
 * They are driven from data expessed in terms of the `Ticket*Import` types below, set up in the `setUp*` functions.
 * The `Ticket*Import` types contain both data to be returned from Trac and data to be returned from Gitea in response to that Trac data.
 *
 * The data created by the `setUp*` functions is passed into the `expect*` functions to set up expections on the mock trac and gitea accessors
 * and on the mock markdown converter. These expectations are governed by the fields in the `Ticket*Import` data.
 *
 * The actual test methods call the top-level `setUp` method to set up all data followed by a selection of `expect` methods to create mock expections
 * for some or all of that data. Finally `importer.ImportTickets(...)` is called to trigger the actual importation we are testing.
 */

/*
 * Data set-up.
 */
// allocators - we give all items unique values so that we can spot any misallocations
func resetAllocators() {
	idCounter = 1000
	intCounter = 20000
	unixTimeCounter = 300000
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
	componentMap  map[string]string
	priorityMap   map[string]string
	resolutionMap map[string]string
	severityMap   map[string]string
	typeMap       map[string]string
	versionMap    map[string]string
)

func initMaps() {
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

func createTicketUserImport(tracUser string, giteaUser string) *TicketUserImport {
	// if we have a gitea user mapping use a unique user id otherwise expect default user id to be used
	var giteaUserID int64
	if giteaUser != "" {
		giteaUserID = allocateID()
	} else {
		giteaUserID = defaultUserID
	}

	user := TicketUserImport{
		tracUser:    tracUser,
		giteaUser:   giteaUser,
		giteaUserID: giteaUserID,
	}

	if tracUser != "" {
		userMap[user.tracUser] = user.giteaUser
	}

	return &user
}

var (
	closedTicketOwner              *TicketUserImport
	closedTicketReporter           *TicketUserImport
	openTicketOwner                *TicketUserImport
	openTicketReporter             *TicketUserImport
	noTracUserTicketOwner          *TicketUserImport
	noTracUserTicketReporter       *TicketUserImport
	unmappedTracUserTicketOwner    *TicketUserImport
	unmappedTracUserTicketReporter *TicketUserImport
)

func setUpTicketUsers(t *testing.T) {
	closedTicketOwner = createTicketUserImport("trac-closed-ticket-owner", "gitea-closed-ticket-owner")
	closedTicketReporter = createTicketUserImport("trac-closed-ticket-reporter", "gitea-closed-ticket-reporter")
	openTicketOwner = createTicketUserImport("trac-open-ticket-owner", "gitea-open-ticket-owner")
	openTicketReporter = createTicketUserImport("trac-open-ticket-reporter", "gitea-open-ticket-reporter")
	noTracUserTicketOwner = createTicketUserImport("", "")
	noTracUserTicketReporter = createTicketUserImport("", "")
	unmappedTracUserTicketOwner = createTicketUserImport("trac-unmapped-user-ticket-owner", "")
	unmappedTracUserTicketReporter = createTicketUserImport("trac-unmapped-user-ticket-reporter", "")
}

var (
	closedTicketComment1Author          *TicketUserImport
	closedTicketComment2Author          *TicketUserImport
	openTicketComment1Author            *TicketUserImport
	openTicketComment2Author            *TicketUserImport
	noTracUserTicketCommentAuthor       *TicketUserImport
	unmappedTracUserTicketCommentAuthor *TicketUserImport
)

func setUpTicketCommentUsers(t *testing.T) {
	closedTicketComment1Author = createTicketUserImport("trac-closed-ticket-comment1-author", "gitea-closed-ticket-comment1-author")
	closedTicketComment2Author = createTicketUserImport("trac-closed-ticket-comment2-author", "gitea-losed-ticket-comment2-author")
	openTicketComment1Author = createTicketUserImport("trac-open-ticket-comment1-author", "gitea-open-ticket-comment1-author")
	openTicketComment2Author = createTicketUserImport("trac-open-ticket-comment2-author", "gitea-open-ticket-comment2-author")
	noTracUserTicketCommentAuthor = createTicketUserImport("", "")
	unmappedTracUserTicketCommentAuthor = createTicketUserImport("trac-unmapped-user-ticket-comment-author", "")
}

var (
	closedTicketAttachment1Author          *TicketUserImport
	closedTicketAttachment2Author          *TicketUserImport
	openTicketAttachment1Author            *TicketUserImport
	openTicketAttachment2Author            *TicketUserImport
	noTracUserTicketAttachmentAuthor       *TicketUserImport
	unmappedTracUserTicketAttachmentAuthor *TicketUserImport
)

func setUpTicketAttachmentUsers(t *testing.T) {
	closedTicketAttachment1Author = createTicketUserImport("trac-closed-ticket-attachment1-author", "gitea-closed-ticket-attachment1-author")
	closedTicketAttachment2Author = createTicketUserImport("trac-closed-ticket-attachment2-author", "gitea-closed-ticket-attachment2-author")
	openTicketAttachment1Author = createTicketUserImport("trac-open-ticket-attachment1-author", "gitea-open-ticket-attachment1-author")
	openTicketAttachment2Author = createTicketUserImport("trac-open-ticket-attachment2-author", "gitea-open-ticket-attachment2-author")
	noTracUserTicketAttachmentAuthor = createTicketUserImport("", "")
	unmappedTracUserTicketAttachmentAuthor = createTicketUserImport("trac-unmapped-user-attachment-author", "")
}

// TicketLabelImport holds the data on a label associated with an imported ticket
type TicketLabelImport struct {
	name         string
	labelName    string
	labelID      int64
	issueLabelID int64
}

func createTicketLabelImport(prefix string, ticketLabelMap map[string]string) *TicketLabelImport {
	ticketLabel := TicketLabelImport{
		name:         prefix + "-name",
		labelName:    prefix + "-label",
		labelID:      allocateID(),
		issueLabelID: allocateID(),
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

	noTracUserTicketComponent  *TicketLabelImport
	noTracUserTicketPriority   *TicketLabelImport
	noTracUserTicketResolution *TicketLabelImport
	noTracUserTicketSeverity   *TicketLabelImport
	noTracUserTicketType       *TicketLabelImport
	noTracUserTicketVersion    *TicketLabelImport

	unmappedTracUserTicketComponent  *TicketLabelImport
	unmappedTracUserTicketPriority   *TicketLabelImport
	unmappedTracUserTicketResolution *TicketLabelImport
	unmappedTracUserTicketSeverity   *TicketLabelImport
	unmappedTracUserTicketType       *TicketLabelImport
	unmappedTracUserTicketVersion    *TicketLabelImport
)

func setUpTicketLabels(t *testing.T) {
	closedTicketComponent = createTicketLabelImport("closed-component", componentMap)
	closedTicketPriority = createTicketLabelImport("closed-priority", priorityMap)
	closedTicketResolution = createTicketLabelImport("closed-resolution", resolutionMap)
	closedTicketSeverity = createTicketLabelImport("closed-severity", severityMap)
	closedTicketType = createTicketLabelImport("closed-type", typeMap)
	closedTicketVersion = createTicketLabelImport("closed-version", versionMap)

	openTicketComponent = createTicketLabelImport("open-component", componentMap)
	openTicketPriority = createTicketLabelImport("open-priority", priorityMap)
	openTicketResolution = createTicketLabelImport("open-resolution", resolutionMap)
	openTicketSeverity = createTicketLabelImport("open-severity", severityMap)
	openTicketType = createTicketLabelImport("open-type", typeMap)
	openTicketVersion = createTicketLabelImport("open-version", versionMap)

	noTracUserTicketComponent = createTicketLabelImport("no-trac-user-component", componentMap)
	noTracUserTicketPriority = createTicketLabelImport("no-trac-user-priority", priorityMap)
	noTracUserTicketResolution = createTicketLabelImport("no-trac-user-resolution", resolutionMap)
	noTracUserTicketSeverity = createTicketLabelImport("no-trac-user-severity", severityMap)
	noTracUserTicketType = createTicketLabelImport("no-trac-user-type", typeMap)
	noTracUserTicketVersion = createTicketLabelImport("no-trac-user-version", versionMap)

	unmappedTracUserTicketComponent = createTicketLabelImport("unmapped-trac-user-component", componentMap)
	unmappedTracUserTicketPriority = createTicketLabelImport("unmapped-trac-user-priority", priorityMap)
	unmappedTracUserTicketResolution = createTicketLabelImport("unmapped-trac-user-resolution", resolutionMap)
	unmappedTracUserTicketSeverity = createTicketLabelImport("unmapped-trac-user-severity", severityMap)
	unmappedTracUserTicketType = createTicketLabelImport("unmapped-trac-user-type", typeMap)
	unmappedTracUserTicketVersion = createTicketLabelImport("unmapped-trac-user-version", versionMap)
}

// TicketCommentImport holds the data on a ticket comment import operation
type TicketCommentImport struct {
	issueCommentID int64
	author         *TicketUserImport
	text           string
	markdownText   string
	time           int64
}

func createTicketCommentImport(prefix string, author *TicketUserImport) *TicketCommentImport {
	return &TicketCommentImport{
		issueCommentID: allocateID(),
		author:         author,
		text:           prefix + " ticket comment text",
		markdownText:   prefix + " ticket comment text after conversion to markdown",
		time:           allocateUnixTime(),
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
	closedTicketComment1          *TicketCommentImport
	closedTicketComment2          *TicketCommentImport
	openTicketComment1            *TicketCommentImport
	openTicketComment2            *TicketCommentImport
	noTracUserTicketComment       *TicketCommentImport
	unmappedTracUserTicketComment *TicketCommentImport
)

func setUpTicketComments(t *testing.T) {
	setUpTicketCommentUsers(t)
	closedTicketComment1 = createTicketCommentImport("closed-ticket-comment1", closedTicketComment1Author)
	closedTicketComment2 = createTicketCommentImport("closed-ticket-comment2", closedTicketComment2Author)

	openTicketComment1 = createTicketCommentImport("open-ticket-comment1", openTicketComment1Author)
	openTicketComment2 = createTicketCommentImport("open-ticket-comment2", openTicketComment2Author)

	noTracUserTicketComment = createTicketCommentImport("no-trac-user-ticket-comment", noTracUserTicketCommentAuthor)
	unmappedTracUserTicketComment = createTicketCommentImport("unmapped-trac-user-ticket-comment", unmappedTracUserTicketCommentAuthor)
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
	comment := createTicketCommentImport(prefix+"-comment-", author)

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
	closedTicketAttachment1          *TicketAttachmentImport
	closedTicketAttachment2          *TicketAttachmentImport
	openTicketAttachment1            *TicketAttachmentImport
	openTicketAttachment2            *TicketAttachmentImport
	noTracUserTicketAttachment       *TicketAttachmentImport
	unmappedTracUserTicketAttachment *TicketAttachmentImport
)

func setUpTicketAttachments(t *testing.T) {
	setUpTicketAttachmentUsers(t)
	closedTicketAttachment1 = createTicketAttachmentImport("closed-ticket-attachment1", closedTicketAttachment1Author)
	closedTicketAttachment2 = createTicketAttachmentImport("closed-ticket-attachment2", closedTicketAttachment2Author)
	openTicketAttachment1 = createTicketAttachmentImport("open-ticket-attachment1", openTicketAttachment1Author)
	openTicketAttachment2 = createTicketAttachmentImport("open-ticket-attachment2", openTicketAttachment2Author)
	noTracUserTicketAttachment = createTicketAttachmentImport("no-trac-user-ticket-attachment", noTracUserTicketAttachmentAuthor)
	unmappedTracUserTicketAttachment = createTicketAttachmentImport("unmapped-trac-user-ticket-attachment", unmappedTracUserTicketAttachmentAuthor)
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
	closedTicket           *TicketImport
	openTicket             *TicketImport
	noTracUserTicket       *TicketImport
	unmappedTracUserTicket *TicketImport
)

// setUpTickets is the top-level setUp method for the ticket tests.
// It should be called by all tests - it is the mock expectations that determines which parts of the set up data are actually used in any test
func setUpTickets(t *testing.T) {
	setUp(t)
	resetAllocators()
	initMaps()
	setUpTicketUsers(t)
	setUpTicketLabels(t)
	setUpTicketComments(t)
	setUpTicketAttachments(t)

	closedTicket = createTicketImport(
		"closed", true,
		closedTicketOwner, closedTicketReporter,
		closedTicketComponent, closedTicketPriority, closedTicketResolution,
		closedTicketSeverity, closedTicketType, closedTicketVersion)
	openTicket = createTicketImport(
		"open", false,
		openTicketOwner, openTicketReporter,
		openTicketComponent, openTicketPriority, openTicketResolution,
		openTicketSeverity, openTicketType, openTicketVersion)
	noTracUserTicket = createTicketImport(
		"noTracUser", false,
		noTracUserTicketOwner, noTracUserTicketReporter,
		noTracUserTicketComponent, noTracUserTicketPriority, noTracUserTicketResolution,
		noTracUserTicketSeverity, noTracUserTicketType, noTracUserTicketVersion)
	unmappedTracUserTicket = createTicketImport(
		"unmappedTracUser", false,
		unmappedTracUserTicketOwner, unmappedTracUserTicketReporter,
		unmappedTracUserTicketComponent, unmappedTracUserTicketPriority, unmappedTracUserTicketResolution,
		unmappedTracUserTicketSeverity, unmappedTracUserTicketType, unmappedTracUserTicketVersion)
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
	// only expect user lookup if we have a trac -> gitea user mapping
	if user.tracUser == "" || user.giteaUser == "" {
		return
	}

	mockGiteaAccessor.
		EXPECT().
		GetUserID(gomock.Eq(user.giteaUser)).
		Return(user.giteaUserID, nil)
}

func expectLabelRetrieval(t *testing.T, label *TicketLabelImport) {
	mockGiteaAccessor.
		EXPECT().
		AddLabel(gomock.Eq(label.labelName), gomock.Any()).
		Return(label.labelID, nil)
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

func expectIssueUserToBeAdded(t *testing.T, ticket *TicketImport, user *TicketUserImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueUser(gomock.Eq(ticket.issueID), gomock.Eq(user.giteaUserID)).
		Return(nil)
}

func expectIssueAssigneeToBeAdded(t *testing.T, ticket *TicketImport, user *TicketUserImport) {
	if user.giteaUser != "" {
		mockGiteaAccessor.
			EXPECT().
			AddIssueAssignee(gomock.Eq(ticket.issueID), gomock.Eq(user.giteaUserID)).
			Return(nil)
	}
}

func expectIssueCreation(t *testing.T, ticket *TicketImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssue(gomock.Any()).
		DoAndReturn(func(issue *gitea.Issue) (int64, error) {
			assertEquals(t, issue.Index, ticket.ticketID)
			assertEquals(t, issue.Summary, ticket.summary)
			assertEquals(t, issue.Description, ticket.descriptionMarkdown)
			assertEquals(t, issue.OriginalAuthorID, int64(0))
			assertEquals(t, issue.OriginalAuthorName, ticket.owner.tracUser)
			assertEquals(t, issue.ReporterID, ticket.reporter.giteaUserID)
			assertEquals(t, issue.Milestone, ticket.milestoneName)
			assertEquals(t, issue.Closed, ticket.closed)
			assertEquals(t, issue.Created, ticket.created)
			return ticket.issueID, nil
		})

	expectIssueAssigneeToBeAdded(t, ticket, ticket.owner)
	expectIssueUserToBeAdded(t, ticket, ticket.reporter)
	if ticket.owner.giteaUser != "" {
		expectIssueUserToBeAdded(t, ticket, ticket.owner)
	}
}

func expectIssueLabelRetrieval(t *testing.T, ticket *TicketImport, ticketLabel *TicketLabelImport) {
	// expect retrieval/creation of underlying label first
	expectLabelRetrieval(t, ticketLabel)

	mockGiteaAccessor.
		EXPECT().
		AddIssueLabel(gomock.Eq(ticket.issueID), gomock.Eq(ticketLabel.labelID)).
		Return(ticketLabel.issueLabelID, nil)
}

func expectIssueCommentCreation(t *testing.T, ticket *TicketImport, ticketComment *TicketCommentImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			assertEquals(t, issueComment.AuthorID, ticketComment.author.giteaUserID)
			assertEquals(t, issueComment.Text, ticketComment.markdownText)
			assertEquals(t, issueComment.Time, ticketComment.time)
			return ticketComment.issueCommentID, nil
		})
	expectIssueUserToBeAdded(t, ticket, ticketComment.author)
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
	expectIssueCommentCreation(t, ticket, ticketComment)
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

func expectIssueAttachmentAddition(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueAttachment(gomock.Eq(ticket.issueID), gomock.Any(), gomock.Eq(ticketAttachment.attachmentPath)).
		DoAndReturn(func(issueID int64, issueAttachment *gitea.IssueAttachment, filePath string) (int64, error) {
			assertEquals(t, issueAttachment.CommentID, ticketAttachment.comment.issueCommentID)
			assertEquals(t, issueAttachment.FileName, ticketAttachment.filename)
			assertEquals(t, issueAttachment.Time, ticketAttachment.comment.time)
			return ticketAttachment.issueAttachmentID, nil
		})
}

func expectAllTicketAttachmentActions(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	expectAllTicketCommentActions(t, ticket, ticketAttachment.comment)
	expectTracAttachmentPathRetrieval(t, ticket, ticketAttachment)
	expectIssueAttachmentAddition(t, ticket, ticketAttachment)
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

func TestImportTicketWithNoTracUser(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, noTracUserTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, noTracUserTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, noTracUserTicket)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, noTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, noTracUserTicket)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithAttachmentButNoTracUser(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, noTracUserTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, noTracUserTicket)

	// expect trac to return us an attachment
	expectTracAttachmentRetrievals(t, noTracUserTicket, noTracUserTicketAttachment)

	// expect all actions for creating Gitea issue attachment from Trac ticket attachment
	expectAllTicketAttachmentActions(t, noTracUserTicket, noTracUserTicketAttachment)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, noTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, noTracUserTicket, noTracUserTicketAttachment.comment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithCommentButNoTracUser(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, noTracUserTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, noTracUserTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, noTracUserTicket)

	// expect trac to return us a comment
	expectTracCommentRetrievals(t, noTracUserTicket, noTracUserTicketComment)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, noTracUserTicket, noTracUserTicketComment)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, noTracUserTicket, noTracUserTicketComment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithUnmappedTracUser(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, unmappedTracUserTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, unmappedTracUserTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, unmappedTracUserTicket)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, unmappedTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, unmappedTracUserTicket)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithAttachmentButUnmappedTracUser(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, unmappedTracUserTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, unmappedTracUserTicket)

	// expect trac to return us an attachment
	expectTracAttachmentRetrievals(t, unmappedTracUserTicket, unmappedTracUserTicketAttachment)

	// expect all actions for creating Gitea issue attachment from Trac ticket attachment
	expectAllTicketAttachmentActions(t, unmappedTracUserTicket, unmappedTracUserTicketAttachment)

	// expect trac to return us no comments
	expectTracCommentRetrievals(t, unmappedTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, unmappedTracUserTicket, unmappedTracUserTicketAttachment.comment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketWithCommentButUnmappedTracUser(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, unmappedTracUserTicket)

	// expect all actions for creating Gitea issues from Trac tickets
	expectAllTicketActions(t, unmappedTracUserTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, unmappedTracUserTicket)

	// expect trac to return us a comment
	expectTracCommentRetrievals(t, unmappedTracUserTicket, unmappedTracUserTicketComment)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, unmappedTracUserTicket, unmappedTracUserTicketComment)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, unmappedTracUserTicket, unmappedTracUserTicketComment)

	// expect repository issue count to be updated
	expectRepoIssueCountUpdate(t, 1, 0)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
