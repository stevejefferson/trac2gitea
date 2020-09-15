// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import "testing"

func TestImportTicketWithComments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, closedTicket)

	// expect trac to return us comment changes
	expectTracChangeRetrievals(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment1)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment2)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us comment changes
	expectTracChangeRetrievals(t, openTicket, openTicketComment1, openTicketComment2)
	expectTracChangeRetrievals(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, openTicket, openTicketComment1)
	expectAllTicketCommentActions(t, openTicket, openTicketComment2)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment1)
	expectAllTicketCommentActions(t, closedTicket, closedTicketComment2)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, openTicketComment1, openTicketComment2)
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketComment1, closedTicketComment2)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us a comment change
	expectTracChangeRetrievals(t, noTracUserTicket, noTracUserTicketComment)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, noTracUserTicket, noTracUserTicketComment)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, noTracUserTicket, noTracUserTicketComment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, noTracUserTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us a comment change
	expectTracChangeRetrievals(t, unmappedTracUserTicket, unmappedTracUserTicketComment)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, unmappedTracUserTicket, unmappedTracUserTicketComment)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, unmappedTracUserTicket, unmappedTracUserTicketComment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, unmappedTracUserTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
