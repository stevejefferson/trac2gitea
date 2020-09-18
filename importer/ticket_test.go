// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"
)

func TestImportClosedTicketOnly(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, closedTicket)

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, closedTicket)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)

	// expect all issue counts to be updated
	expectIssueCountUpdates(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, openTicket)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect all issue counts to be updated
	expectIssueCountUpdates(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, closedTicket)
	expectTracChangeRetrievals(t, openTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket)
	expectIssueUpdateTimeSetToLatestOf(t, openTicket)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)
	expectIssueCommentCountUpdate(t, openTicket)

	// expect all issue counts to be updated
	expectIssueCountUpdates(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, noTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, noTracUserTicket)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, noTracUserTicket)

	// expect all issue counts to be updated
	expectIssueCountUpdates(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, unmappedTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, unmappedTracUserTicket)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, unmappedTracUserTicket)

	// expect all issue counts to be updated
	expectIssueCountUpdates(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
