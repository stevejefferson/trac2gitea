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
 * Set up for ticket/issue status parts of ticket tests.
 * Contains:
 * - ticket status change and associated data (users, labels etc.)
 * - expectations for use with ticket status changes.
 */

var (
	closeStatusChangeAuthor  *TicketUserImport
	reopenStatusChangeAuthor *TicketUserImport
)

func setUpTicketStatusChangeUsers(t *testing.T) {
	closeStatusChangeAuthor = createTicketUserImport("trac-close-status-change-author", "gitea-close-status-change-author")
	reopenStatusChangeAuthor = createTicketUserImport("trac-reopen-status-change-author", "gitea-reopen-status-change-author")
}

func createCloseTicketChangeImport(author *TicketUserImport, isClose bool) *TicketChangeImport {
	return &TicketChangeImport{
		tracChangeType: trac.TicketStatusChange,
		issueCommentID: allocateID(),
		author:         author,
		time:           allocateUnixTime(),
		isClose:        isClose,
	}
}

var (
	closeTicketChange  *TicketChangeImport
	reopenTicketChange *TicketChangeImport
)

func setUpTicketStatusChanges(t *testing.T) {
	setUpTicketStatusChangeUsers(t)
	closeTicketChange = createCloseTicketChangeImport(closeStatusChangeAuthor, true)
	reopenTicketChange = createCloseTicketChangeImport(reopenStatusChangeAuthor, false)
}

func expectIssueCommentCreationForStatusChange(t *testing.T, ticket *TicketImport, ticketStatus *TicketChangeImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			if ticketStatus.isClose {
				assertEquals(t, issueComment.CommentType, gitea.CloseIssueCommentType)
			} else {
				assertEquals(t, issueComment.CommentType, gitea.ReopenIssueCommentType)
			}
			assertEquals(t, issueComment.AuthorID, ticketStatus.author.giteaUserID)
			assertEquals(t, issueComment.Time, ticketStatus.time)
			return ticketStatus.issueCommentID, nil
		})
	if ticketStatus.author.giteaUser != "" {
		expectIssueParticipantToBeAdded(t, ticket, ticketStatus.author)
	}
}

func expectAllTicketStatusActions(t *testing.T, ticket *TicketImport, ticketStatus *TicketChangeImport) {
	// expect to lookup Gitea equivalent of author of Trac ticket change
	expectUserLookup(t, ticketStatus.author)

	// expect creation of issue comment for ticket status change
	expectIssueCommentCreationForStatusChange(t, ticket, ticketStatus)
}
