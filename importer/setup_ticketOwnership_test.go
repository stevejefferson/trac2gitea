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
 * Set up for ticket/issue ownership parts of ticket tests.
 * Contains:
 * - ticket ownership and associated data (users, labels etc.)
 * - expectations for use with ticket ownerships.
 */

var (
	ownershipChangeAuthor     *TicketUserImport
	ownershipChangePrevOwner  *TicketUserImport
	ownershipChangeNewOwner   *TicketUserImport
	ownershipRemovalAuthor    *TicketUserImport
	ownershipRemovalPrevOwner *TicketUserImport
	ownershipRemovalNewOwner  *TicketUserImport
)

func setUpTicketOwnershipChangeUsers(t *testing.T) {
	ownershipChangeAuthor = createTicketUserImport("trac-owner-change-author", "gitea-owner-change-author")
	ownershipChangePrevOwner = createTicketUserImport("trac-owner-change-prev-owner", "gitea-owner-change-prev-owner")
	ownershipChangeNewOwner = createTicketUserImport("trac-owner-change-new-owner", "gitea-owner-change-new-owner")
	ownershipRemovalAuthor = createTicketUserImport("trac-owner-removal-author", "gitea-owner-removal-author")
	ownershipRemovalPrevOwner = createTicketUserImport("trac-owner-removal-prev-owner", "gitea-owner-removal-prev-owner")
	ownershipRemovalNewOwner = createTicketUserImport("", "")
}

func createOwnershipTicketChangeImport(author *TicketUserImport, prevOwner *TicketUserImport, owner *TicketUserImport) *TicketChangeImport {
	return &TicketChangeImport{
		tracChangeType: trac.TicketOwnershipChange,
		issueCommentID: allocateID(),
		author:         author,
		owner:          owner,
		prevOwner:      prevOwner,
		text:           "",
		markdownText:   "",
		time:           allocateUnixTime(),
	}
}

var (
	assigneeTicketChange        *TicketChangeImport
	assigneeRemovalTicketChange *TicketChangeImport
)

func setUpTicketOwnershipChanges(t *testing.T) {
	setUpTicketOwnershipChangeUsers(t)
	assigneeTicketChange = createOwnershipTicketChangeImport(ownershipChangeAuthor, ownershipChangePrevOwner, ownershipChangeNewOwner)
	assigneeRemovalTicketChange = createOwnershipTicketChangeImport(ownershipRemovalAuthor, ownershipRemovalPrevOwner, ownershipRemovalNewOwner)
}

func expectIssueCommentCreationForOwnershipChange(t *testing.T, ticket *TicketImport, ticketOwnership *TicketChangeImport) {
	// expect to look up Gitea user corresponding to previous ticket owner
	expectUserLookup(t, ticketOwnership.prevOwner)

	// expect to look up Gitea user corresponding to new ticket owner
	expectUserLookup(t, ticketOwnership.owner)

	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			assertEquals(t, issueComment.CommentType, gitea.AssigneeIssueCommentType)
			assertEquals(t, issueComment.AuthorID, ticketOwnership.author.giteaUserID)
			if ticketOwnership.owner.tracUser != "" {
				assertEquals(t, issueComment.AssigneeID, ticketOwnership.owner.giteaUserID)
			} else {
				assertEquals(t, issueComment.RemovedAssignee, ticketOwnership.prevOwner.giteaUserID)
			}
			assertEquals(t, issueComment.Time, ticketOwnership.time)
			return ticketOwnership.issueCommentID, nil
		})
	expectIssueUserToBeAdded(t, ticket, ticketOwnership.author)
}

func expectAllTicketOwnershipActions(t *testing.T, ticket *TicketImport, ticketOwnership *TicketChangeImport) {
	// expect to lookup Gitea equivalent of author of Trac ticket ownership change
	expectUserLookup(t, ticketOwnership.author)

	// expect retrieval/creation of issue comment for ticket ownership change
	expectIssueCommentCreationForOwnershipChange(t, ticket, ticketOwnership)
}
