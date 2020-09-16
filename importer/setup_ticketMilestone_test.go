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
 * Set up for ticket/issue milestone parts of ticket tests.
 * Contains:
 * - ticket milestone change and associated data (users, labels etc.)
 * - expectations for use with ticket milestone changes.
 */

var (
	milestoneChangeAuthor *TicketUserImport
)

func setUpTicketMilestoneChangeUsers(t *testing.T) {
	milestoneChangeAuthor = createTicketUserImport("trac-milestone-change-author", "gitea-milestone-change-author")
}

// TicketMilestoneImport describes a milestone appearing in a ticket milestone change
type TicketMilestoneImport struct {
	milestoneName string
	milestoneID   int64
}

func createTicketMilestoneImport(name string) *TicketMilestoneImport {
	return &TicketMilestoneImport{milestoneName: name, milestoneID: allocateID()}
}

var (
	milestone1 *TicketMilestoneImport
	milestone2 *TicketMilestoneImport
)

func setUpTicketMilestones(t *testing.T) {
	milestone1 = createTicketMilestoneImport("milestone1")
	milestone2 = createTicketMilestoneImport("milestone2")
}

func createMilestoneTicketChangeImport(author *TicketUserImport, prevMilestone *TicketMilestoneImport, milestone *TicketMilestoneImport) *TicketChangeImport {
	return &TicketChangeImport{
		tracChangeType: trac.TicketMilestoneChange,
		issueCommentID: allocateID(),
		author:         author,
		prevMilestone:  prevMilestone,
		milestone:      milestone,
		time:           allocateUnixTime(),
	}
}

var (
	milestoneTicketChange *TicketChangeImport
)

func setUpTicketMilestoneChanges(t *testing.T) {
	setUpTicketMilestoneChangeUsers(t)
	setUpTicketMilestones(t)
	milestoneTicketChange = createMilestoneTicketChangeImport(milestoneChangeAuthor, milestone1, milestone2)
}

func expectTicketMilestoneRetrieval(t *testing.T, milestone *TicketMilestoneImport) {
	mockGiteaAccessor.
		EXPECT().
		GetMilestoneID(milestone.milestoneName).
		Return(milestone.milestoneID, nil)
}

func expectIssueCommentCreationForMilestoneChange(t *testing.T, ticket *TicketImport, ticketMilestone *TicketChangeImport) {

	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			assertEquals(t, issueComment.CommentType, gitea.MilestoneIssueCommentType)
			assertEquals(t, issueComment.AuthorID, ticketMilestone.author.giteaUserID)
			assertEquals(t, issueComment.OldMilestoneID, ticketMilestone.prevMilestone.milestoneID)
			assertEquals(t, issueComment.MilestoneID, ticketMilestone.milestone.milestoneID)
			assertEquals(t, issueComment.Time, ticketMilestone.time)
			return ticketMilestone.issueCommentID, nil
		})

	if ticketMilestone.author.giteaUser != "" {
		expectIssueParticipantToBeAdded(t, ticket, ticketMilestone.author)
	}
}

func expectAllTicketMilestoneActions(t *testing.T, ticket *TicketImport, ticketMilestone *TicketChangeImport) {
	// expect to lookup Gitea equivalent of author of Trac ticket change
	expectUserLookup(t, ticketMilestone.author)

	// expect lookups of old and new milestones
	expectTicketMilestoneRetrieval(t, ticketMilestone.prevMilestone)
	expectTicketMilestoneRetrieval(t, ticketMilestone.milestone)

	// expect creation of issue comment for ticket milestone change
	expectIssueCommentCreationForStatusChange(t, ticket, ticketMilestone)
}
