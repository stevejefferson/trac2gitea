// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import "testing"

func TestImportTicketWithAttachments(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, closedTicket)

	// expect trac to return us attachments
	expectTracAttachmentRetrievals(t, closedTicket, closedTicketAttachment1, closedTicketAttachment2)

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, closedTicket)

	// expect all actions for creating Gitea issue attachments from Trac ticket attachments
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment1)
	expectAllTicketAttachmentActions(t, closedTicket, closedTicketAttachment2)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketAttachment1.comment, closedTicketAttachment2.comment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, closedTicket)
	expectTracChangeRetrievals(t, openTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, closedTicket, closedTicketAttachment1.comment, closedTicketAttachment2.comment)
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, openTicketAttachment1.comment, openTicketAttachment2.comment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, closedTicket)
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, noTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, noTracUserTicket, noTracUserTicketAttachment.comment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, noTracUserTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

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

	// expect trac to return us no changes
	expectTracChangeRetrievals(t, unmappedTracUserTicket)

	// expect issues update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, unmappedTracUserTicket, unmappedTracUserTicketAttachment.comment)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, unmappedTracUserTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
