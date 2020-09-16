// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import "testing"

func TestImportTicketComponentAddition(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one component addition
	expectTracChangeRetrievals(t, openTicket, componentAddTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, componentAddTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, componentAddTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketComponentAmend(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one component amend
	expectTracChangeRetrievals(t, openTicket, componentAmendTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, componentAmendTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, componentAmendTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketComponentRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one component removal
	expectTracChangeRetrievals(t, openTicket, componentRemoveTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, componentRemoveTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, componentRemoveTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketPriorityAddition(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one priority addition
	expectTracChangeRetrievals(t, openTicket, priorityAddTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, priorityAddTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, priorityAddTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketPriorityAmend(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one priority amend
	expectTracChangeRetrievals(t, openTicket, priorityAmendTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, priorityAmendTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, priorityAmendTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketPriorityRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one priority removal
	expectTracChangeRetrievals(t, openTicket, priorityRemoveTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, priorityRemoveTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, priorityRemoveTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketResolutionAddition(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one resolution addition
	expectTracChangeRetrievals(t, openTicket, resolutionAddTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, resolutionAddTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, resolutionAddTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketResolutionAmend(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one resolution amend
	expectTracChangeRetrievals(t, openTicket, resolutionAmendTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, resolutionAmendTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, resolutionAmendTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketResolutionRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one resolution removal
	expectTracChangeRetrievals(t, openTicket, resolutionRemoveTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, resolutionRemoveTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, resolutionRemoveTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketSeverityAddition(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one severity addition
	expectTracChangeRetrievals(t, openTicket, severityAddTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, severityAddTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, severityAddTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketSeverityAmend(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one severity amend
	expectTracChangeRetrievals(t, openTicket, severityAmendTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, severityAmendTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, severityAmendTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketSeverityRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one severity removal
	expectTracChangeRetrievals(t, openTicket, severityRemoveTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, severityRemoveTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, severityRemoveTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketTypeAddition(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one type addition
	expectTracChangeRetrievals(t, openTicket, typeAddTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, typeAddTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, typeAddTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketTypeAmend(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one type amend
	expectTracChangeRetrievals(t, openTicket, typeAmendTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, typeAmendTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, typeAmendTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketTypeRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one type removal
	expectTracChangeRetrievals(t, openTicket, typeRemoveTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, typeRemoveTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, typeRemoveTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketVersionAddition(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one version addition
	expectTracChangeRetrievals(t, openTicket, versionAddTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, versionAddTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, versionAddTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketVersionAmend(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one version amend
	expectTracChangeRetrievals(t, openTicket, versionAmendTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, versionAmendTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, versionAmendTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}

func TestImportTicketVersionRemoval(t *testing.T) {
	setUpTickets(t)
	defer tearDown(t)

	// first thing to expect is retrieval of ticket from Trac
	expectTracTicketRetrievals(t, openTicket)

	// expect all actions for creating Gitea issue from Trac ticket
	expectAllTicketActions(t, openTicket)

	// expect trac to return us no attachments
	expectTracAttachmentRetrievals(t, openTicket)

	// expect trac to return us one version removal
	expectTracChangeRetrievals(t, openTicket, versionRemoveTicketChange)

	// expect all actions for creating Gitea comments from Trac ticket label changes
	expectAllTicketLabelActions(t, openTicket, versionRemoveTicketChange)

	// expect issue update time to be updated
	expectIssueUpdateTimeSetToLatestOf(t, openTicket, versionRemoveTicketChange)

	// expect issue comment count to be updated
	expectIssueCommentCountUpdate(t, openTicket)

	// expect repository issue counts to be updated
	expectRepoIssueCountsUpdate(t)

	dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
}
