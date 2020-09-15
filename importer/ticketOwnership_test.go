// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import "testing"

func TestImportTicketOwnershipChange(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one ownership change
	expectTracChangeRetrievals(t, openTicket, assigneeTicketChange)

	// expect all actions for creating Gitea assignments from Trac ticket ownership changes
	expectAllTicketOwnershipActions(t, openTicket, assigneeTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, assigneeTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketOwnershipRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one ownership change
	expectTracChangeRetrievals(t, openTicket, assigneeRemovalTicketChange)

	// expect all actions for creating Gitea assignments from Trac ticket ownership changes
	expectAllTicketOwnershipActions(t, openTicket, assigneeRemovalTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, assigneeRemovalTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
