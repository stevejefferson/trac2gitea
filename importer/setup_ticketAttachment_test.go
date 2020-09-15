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
* Set up for ticket/issue attachment parts of ticket tests.
* Contains:
* - ticket attachment and associated data (users, labels etc.)
* - expectations for use with ticket attachments.
 */

var (
	closedTicketAttachment1Author          *TicketUserImport
	closedTicketAttachment2Author          *TicketUserImport
	openTicketAttachment1Author            *TicketUserImport
	openTicketAttachment2Author            *TicketUserImport
	noTracUserTicketAttachmentAuthor       *TicketUserImport
	unmappedTracUserTicketAttachmentAuthor *TicketUserImport
)

func setUpTicketAttachmentUsers(t *testing.T) {
	closedTicketAttachment1Author = createTicketUserImport("trac-closed-ticket-attachment1-author", "gitea-closed-ticket-attachment1-author")
	closedTicketAttachment2Author = createTicketUserImport("trac-closed-ticket-attachment2-author", "gitea-closed-ticket-attachment2-author")
	openTicketAttachment1Author = createTicketUserImport("trac-open-ticket-attachment1-author", "gitea-open-ticket-attachment1-author")
	openTicketAttachment2Author = createTicketUserImport("trac-open-ticket-attachment2-author", "gitea-open-ticket-attachment2-author")
	noTracUserTicketAttachmentAuthor = createTicketUserImport("", "")
	unmappedTracUserTicketAttachmentAuthor = createTicketUserImport("trac-unmapped-user-attachment-author", "")
}

// TicketAttachmentImport holds data on a ticket attachment import operation
type TicketAttachmentImport struct {
	issueAttachmentID int64
	comment           *TicketChangeImport
	filename          string
	attachmentPath    string
	size              int64
}

func createTicketAttachmentImport(prefix string, author *TicketUserImport) *TicketAttachmentImport {
	// express part of attachment data in terms of the comment that will appear in Gitea to describe it
	comment := createCommentTicketChangeImport(prefix+"-comment-", author)

	// trac attachment path must have final directory of at least 12 chars (the trac UUID)
	attachmentFile := prefix + "-attachment.file"
	attachmentPath := "/path/to/attachment/" + prefix + "123456789012/" + attachmentFile

	return &TicketAttachmentImport{
		issueAttachmentID: allocateID(),
		comment:           comment,
		filename:          attachmentFile,
		attachmentPath:    attachmentPath,
		size:              allocateInt(),
	}
}

func createTracTicketAttachment(ticket *TicketImport, ticketAttachment *TicketAttachmentImport) *trac.TicketAttachment {
	return &trac.TicketAttachment{
		TicketID:    ticket.ticketID,
		Size:        ticketAttachment.size,
		Author:      ticketAttachment.comment.author.tracUser,
		FileName:    ticketAttachment.filename,
		Description: ticketAttachment.comment.text,
		Time:        ticketAttachment.comment.time,
	}
}

var (
	closedTicketAttachment1          *TicketAttachmentImport
	closedTicketAttachment2          *TicketAttachmentImport
	openTicketAttachment1            *TicketAttachmentImport
	openTicketAttachment2            *TicketAttachmentImport
	noTracUserTicketAttachment       *TicketAttachmentImport
	unmappedTracUserTicketAttachment *TicketAttachmentImport
)

func setUpTicketAttachments(t *testing.T) {
	setUpTicketAttachmentUsers(t)
	closedTicketAttachment1 = createTicketAttachmentImport("closed-ticket-attachment1", closedTicketAttachment1Author)
	closedTicketAttachment2 = createTicketAttachmentImport("closed-ticket-attachment2", closedTicketAttachment2Author)
	openTicketAttachment1 = createTicketAttachmentImport("open-ticket-attachment1", openTicketAttachment1Author)
	openTicketAttachment2 = createTicketAttachmentImport("open-ticket-attachment2", openTicketAttachment2Author)
	noTracUserTicketAttachment = createTicketAttachmentImport("no-trac-user-ticket-attachment", noTracUserTicketAttachmentAuthor)
	unmappedTracUserTicketAttachment = createTicketAttachmentImport("unmapped-trac-user-ticket-attachment", unmappedTracUserTicketAttachmentAuthor)
}

func expectTracAttachmentRetrievals(t *testing.T, ticket *TicketImport, ticketAttachments ...*TicketAttachmentImport) {
	// expect trac accessor to return each of our trac ticket attachments
	mockTracAccessor.
		EXPECT().
		GetTicketAttachments(gomock.Eq(ticket.ticketID), gomock.Any()).
		DoAndReturn(func(ticketID int64, handlerFn func(attachment *trac.TicketAttachment) error) error {
			for _, ticketAttachment := range ticketAttachments {
				tracAttachment := createTracTicketAttachment(ticket, ticketAttachment)
				handlerFn(tracAttachment)
			}
			return nil
		})
}

func expectTracAttachmentPathRetrieval(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	mockTracAccessor.
		EXPECT().
		GetTicketAttachmentPath(gomock.Any()).
		DoAndReturn(func(tracAttachment *trac.TicketAttachment) string {
			assertEquals(t, tracAttachment.TicketID, ticket.ticketID)
			assertEquals(t, tracAttachment.FileName, ticketAttachment.filename)
			assertEquals(t, tracAttachment.Size, ticketAttachment.size)
			assertEquals(t, tracAttachment.Time, ticketAttachment.comment.time)
			assertEquals(t, tracAttachment.Author, ticketAttachment.comment.author.tracUser)
			assertTrue(t, strings.Contains(tracAttachment.Description, ticketAttachment.comment.text))
			return ticketAttachment.attachmentPath
		})
}

func expectIssueAttachmentAddition(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueAttachment(gomock.Eq(ticket.issueID), gomock.Any(), gomock.Eq(ticketAttachment.attachmentPath)).
		DoAndReturn(func(issueID int64, issueAttachment *gitea.IssueAttachment, filePath string) (int64, error) {
			assertEquals(t, issueAttachment.CommentID, ticketAttachment.comment.issueCommentID)
			assertEquals(t, issueAttachment.FileName, ticketAttachment.filename)
			assertEquals(t, issueAttachment.Time, ticketAttachment.comment.time)
			return ticketAttachment.issueAttachmentID, nil
		})
}

func expectAllTicketAttachmentActions(t *testing.T, ticket *TicketImport, ticketAttachment *TicketAttachmentImport) {
	expectAllTicketCommentActions(t, ticket, ticketAttachment.comment)
	expectTracAttachmentPathRetrieval(t, ticket, ticketAttachment)
	expectIssueAttachmentAddition(t, ticket, ticketAttachment)
}
