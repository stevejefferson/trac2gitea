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
	componentLabel1  *TicketLabelImport
	componentLabel2  *TicketLabelImport
	priorityLabel1   *TicketLabelImport
	priorityLabel2   *TicketLabelImport
	resolutionLabel1 *TicketLabelImport
	resolutionLabel2 *TicketLabelImport
	severityLabel1   *TicketLabelImport
	severityLabel2   *TicketLabelImport
	typeLabel1       *TicketLabelImport
	typeLabel2       *TicketLabelImport
	versionLabel1    *TicketLabelImport
	versionLabel2    *TicketLabelImport
)

func setUpTicketLabels(t *testing.T) {
	componentLabel1 = createTicketLabelImport("component1", componentMap)
	componentLabel2 = createTicketLabelImport("component2", componentMap)
	priorityLabel1 = createTicketLabelImport("priority1", priorityMap)
	priorityLabel2 = createTicketLabelImport("priority2", priorityMap)
	resolutionLabel1 = createTicketLabelImport("resolution1", resolutionMap)
	resolutionLabel2 = createTicketLabelImport("resolution2", resolutionMap)
	severityLabel1 = createTicketLabelImport("severity1", severityMap)
	severityLabel2 = createTicketLabelImport("severity2", severityMap)
	typeLabel1 = createTicketLabelImport("type1", typeMap)
	typeLabel2 = createTicketLabelImport("type2", typeMap)
	versionLabel1 = createTicketLabelImport("version1", versionMap)
	versionLabel2 = createTicketLabelImport("version2", versionMap)
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
		ComponentName:  ticket.componentLabel.tracName,
		PriorityName:   ticket.priorityLabel.tracName,
		ResolutionName: ticket.resolutionLabel.tracName,
		SeverityName:   ticket.severityLabel.tracName,
		TypeName:       ticket.typeLabel.tracName,
		VersionName:    ticket.versionLabel.tracName,
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
	setUpTicketLabelChanges(t)
	setUpTicketMilestoneChanges(t)
	setUpTicketOwnershipChanges(t)
	setUpTicketStatusChanges(t)
	setUpTicketSummaryChanges(t)
	setUpTicketAttachments(t)

	closedTicket = createTicketImport(
		"closed", true,
		closedTicketOwner, closedTicketReporter,
		componentLabel1, priorityLabel1, resolutionLabel1, severityLabel1, typeLabel1, versionLabel1)
	openTicket = createTicketImport(
		"open", false,
		openTicketOwner, openTicketReporter,
		componentLabel2, priorityLabel2, resolutionLabel2, severityLabel2, typeLabel2, versionLabel2)
	noTracUserTicket = createTicketImport(
		"noTracUser", false,
		noTracUserTicketOwner, noTracUserTicketReporter,
		componentLabel1, priorityLabel1, resolutionLabel1, severityLabel1, typeLabel1, versionLabel1)
	unmappedTracUserTicket = createTicketImport(
		"unmappedTracUser", false,
		unmappedTracUserTicketOwner, unmappedTracUserTicketReporter,
		componentLabel1, priorityLabel1, resolutionLabel1, severityLabel1, typeLabel1, versionLabel1)
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
			assertEquals(t, issue.OriginalAuthorID, gitea.NullID)
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

	// expect creation of all labels from Trac ticket appearing in the Gitea issue
	expectIssueLabelCreation(t, ticket, ticket.componentLabel)
	expectIssueLabelCreation(t, ticket, ticket.priorityLabel)
	expectIssueLabelCreation(t, ticket, ticket.resolutionLabel)
	expectIssueLabelCreation(t, ticket, ticket.severityLabel)
	expectIssueLabelCreation(t, ticket, ticket.typeLabel)
	expectIssueLabelCreation(t, ticket, ticket.versionLabel)
}
