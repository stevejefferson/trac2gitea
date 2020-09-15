// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import "testing"

/*
 * Ticket "compound" tests.
 * These tests address tickets with combinations of attachments,comments,ownership changes etc.
 */

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

	// expect trac to return us comment changes
	expectTracChangeRetrievals(t, openTicket, openTicketComment1, openTicketComment2)

	// expect all actions for creating Gitea issue comments from Trac ticket comments
	expectAllTicketCommentActions(t, openTicket, openTicketComment1)
	expectAllTicketCommentActions(t, openTicket, openTicketComment2)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket,
		openTicketComment1, openTicketComment2, openTicketAttachment1.comment, openTicketAttachment2.comment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us comment changes
	expectTracChangeRetrievals(t, openTicket, openTicketComment1, openTicketComment2)
	expectTracChangeRetrievals(t, closedTicket, closedTicketComment1, closedTicketComment2)

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

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
