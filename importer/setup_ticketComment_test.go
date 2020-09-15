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
 * Set up for ticket/issue comment parts of ticket tests.
 * Contains:
 * - ticketcomment and associated data (users, labels etc.)
 * - expectations for use with ticket comments.
 */

var (
	closedTicketComment1Author          *TicketUserImport
	closedTicketComment2Author          *TicketUserImport
	openTicketComment1Author            *TicketUserImport
	openTicketComment2Author            *TicketUserImport
	noTracUserTicketCommentAuthor       *TicketUserImport
	unmappedTracUserTicketCommentAuthor *TicketUserImport
)

func setUpTicketCommentUsers(t *testing.T) {
	closedTicketComment1Author = createTicketUserImport("trac-closed-ticket-comment1-author", "gitea-closed-ticket-comment1-author")
	closedTicketComment2Author = createTicketUserImport("trac-closed-ticket-comment2-author", "gitea-losed-ticket-comment2-author")
	openTicketComment1Author = createTicketUserImport("trac-open-ticket-comment1-author", "gitea-open-ticket-comment1-author")
	openTicketComment2Author = createTicketUserImport("trac-open-ticket-comment2-author", "gitea-open-ticket-comment2-author")
	noTracUserTicketCommentAuthor = createTicketUserImport("", "")
	unmappedTracUserTicketCommentAuthor = createTicketUserImport("trac-unmapped-user-ticket-comment-author", "")
}

var (
	closedTicketComment1          *TicketChangeImport
	closedTicketComment2          *TicketChangeImport
	openTicketComment1            *TicketChangeImport
	openTicketComment2            *TicketChangeImport
	noTracUserTicketComment       *TicketChangeImport
	unmappedTracUserTicketComment *TicketChangeImport
)

func createCommentTicketChangeImport(prefix string, author *TicketUserImport) *TicketChangeImport {
	return &TicketChangeImport{
		tracChangeType: trac.TicketCommentChange,
		issueCommentID: allocateID(),
		author:         author,
		owner:          nil,
		prevOwner:      nil,
		text:           prefix + " ticket comment text",
		markdownText:   prefix + " ticket comment text after conversion to markdown",
		time:           allocateUnixTime(),
	}
}

func setUpTicketComments(t *testing.T) {
	setUpTicketCommentUsers(t)
	closedTicketComment1 = createCommentTicketChangeImport("closed-ticket-comment1", closedTicketComment1Author)
	closedTicketComment2 = createCommentTicketChangeImport("closed-ticket-comment2", closedTicketComment2Author)

	openTicketComment1 = createCommentTicketChangeImport("open-ticket-comment1", openTicketComment1Author)
	openTicketComment2 = createCommentTicketChangeImport("open-ticket-comment2", openTicketComment2Author)

	noTracUserTicketComment = createCommentTicketChangeImport("no-trac-user-ticket-comment", noTracUserTicketCommentAuthor)
	unmappedTracUserTicketComment = createCommentTicketChangeImport("unmapped-trac-user-ticket-comment", unmappedTracUserTicketCommentAuthor)
}

func expectIssueCommentCreationForComment(t *testing.T, ticket *TicketImport, ticketComment *TicketChangeImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			assertEquals(t, issueComment.CommentType, gitea.CommentIssueCommentType)
			assertEquals(t, issueComment.AuthorID, ticketComment.author.giteaUserID)
			assertEquals(t, issueComment.Text, ticketComment.markdownText)
			assertEquals(t, issueComment.Time, ticketComment.time)
			return ticketComment.issueCommentID, nil
		})
	expectIssueUserToBeAdded(t, ticket, ticketComment.author)
}

func expectTicketCommentMarkdownConversion(t *testing.T, ticket *TicketImport, ticketComment *TicketChangeImport) {
	mockMarkdownConverter.
		EXPECT().
		TicketConvert(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, text string) string {
			assertTrue(t, strings.Contains(text, ticketComment.text))
			return ticketComment.markdownText
		})
}

func expectAllTicketCommentActions(t *testing.T, ticket *TicketImport, ticketComment *TicketChangeImport) {
	// expect to lookup Gitea equivalents of Trac ticket comment author
	expectUserLookup(t, ticketComment.author)

	// expect to convert ticket comment text to markdown
	expectTicketCommentMarkdownConversion(t, ticket, ticketComment)

	// expect retrieval/creation of issue comment for ticket comment
	expectIssueCommentCreationForComment(t, ticket, ticketComment)
}
