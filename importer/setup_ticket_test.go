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
 * Set up for ticket/issue parts of ticket tests.
 * Contains:
 * - ticket data types
 * - ticket and associated data (users, labels etc.)
 * - expectations for use with tickets.
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
	updated             int64
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
		updated:             allocateUnixTime(),
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
		Updated:        ticket.updated,
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
	setUpTicketMilestoneChanges(t)
	setUpTicketOwnershipChanges(t)
	setUpTicketStatusChanges(t)
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
	// expect to record original trac user where ticket owner has no Gitea mapping
	originalAuthorName := ""
	if ticket.owner.giteaUser == "" {
		originalAuthorName = ticket.owner.tracUser
	}

	mockGiteaAccessor.
		EXPECT().
		AddIssue(gomock.Any()).
		DoAndReturn(func(issue *gitea.Issue) (int64, error) {
			assertEquals(t, issue.Index, ticket.ticketID)
			assertEquals(t, issue.Summary, ticket.summary)
			assertEquals(t, issue.Description, ticket.descriptionMarkdown)
			assertEquals(t, issue.OriginalAuthorID, int64(0))
			assertEquals(t, issue.OriginalAuthorName, originalAuthorName)
			assertEquals(t, issue.ReporterID, ticket.reporter.giteaUserID)
			assertEquals(t, issue.Milestone, ticket.milestoneName)
			assertEquals(t, issue.Closed, ticket.closed)
			assertEquals(t, issue.Created, ticket.created)
			return ticket.issueID, nil
		})

	// reporter (or default user if no Gitea mapping) will always be set as issue participant
	expectIssueParticipantToBeAdded(t, ticket, ticket.reporter)
	if ticket.owner.giteaUser != "" {
		expectIssueAssigneeToBeAdded(t, ticket, ticket.owner)
		expectIssueParticipantToBeAdded(t, ticket, ticket.owner)
	}
}

func expectIssueUpdateTimeSetToLatestOf(t *testing.T, ticket *TicketImport, ticketComments ...*TicketChangeImport) {
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

func expectIssueCommentCountUpdate(t *testing.T, ticket *TicketImport) {
	mockGiteaAccessor.
		EXPECT().
		UpdateIssueCommentCount(gomock.Eq(ticket.issueID)).
		Return(nil)
}

func expectRepoIssueCountsUpdate(t *testing.T) {
	mockGiteaAccessor.
		EXPECT().
		UpdateRepoIssueCounts().
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
