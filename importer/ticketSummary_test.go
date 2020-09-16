// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import "testing"

func TestImportTicketSummary(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one summary change
	expectTracChangeRetrievals(t, openTicket, summaryTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket summary changes
	expectAllTicketSummaryActions(t, openTicket, summaryTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, summaryTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
