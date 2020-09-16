// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

/*
 * Set up for ticket/issue summary parts of ticket tests.
 * Contains:
 * - ticket summary change and associated data (users, labels etc.)
 * - expectations for use with ticket summary changes.
 */

var (
	summaryChangeAuthor *TicketUserImport
)

func setUpTicketSummaryChangeUsers(t *testing.T) {
	summaryChangeAuthor = createTicketUserImport("trac-summary-change-author", "gitea-summary-change-author")
}

func createSummaryTicketChangeImport(author *TicketUserImport, prevSummary string, summary string) *TicketChangeImport {
	return &TicketChangeImport{
		tracChangeType: trac.TicketSummaryChange,
		issueCommentID: allocateID(),
		author:         author,
		prevSummary:    prevSummary,
		summary:        summary,
		time:           allocateUnixTime(),
	}
}

var (
	summaryTicketChange *TicketChangeImport
)

const (
	ticketSummary1 = "summary#1-of-ticket"
	ticketSummary2 = "summary#2-of-ticket"
)

func setUpTicketSummaryChanges(t *testing.T) {
	setUpTicketSummaryChangeUsers(t)
	summaryTicketChange = createSummaryTicketChangeImport(summaryChangeAuthor, ticketSummary1, ticketSummary2)
}

func expectIssueCommentCreationForSummaryChange(t *testing.T, ticket *TicketImport, ticketSummary *TicketChangeImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			assertEquals(t, issueComment.CommentType, gitea.TitleIssueCommentType)
			assertEquals(t, issueComment.AuthorID, ticketSummary.author.giteaUserID)
			assertEquals(t, issueComment.OldTitle, ticketSummary.prevSummary)
			assertEquals(t, issueComment.Title, ticketSummary.summary)
			assertEquals(t, issueComment.Time, ticketSummary.time)
			return ticketSummary.issueCommentID, nil
		})

	if ticketSummary.author.giteaUser != "" {
		expectIssueParticipantToBeAdded(t, ticket, ticketSummary.author)
	}
}

func expectAllTicketSummaryActions(t *testing.T, ticket *TicketImport, ticketSummary *TicketChangeImport) {
	// expect to lookup Gitea equivalent of author of Trac ticket change
	expectUserLookup(t, ticketSummary.author)

	// expect creation of issue comment for ticket summary change
	expectIssueCommentCreationForSummaryChange(t, ticket, ticketSummary)
}
