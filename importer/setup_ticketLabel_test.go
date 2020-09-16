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
 * Set up for ticket/issue label parts of ticket tests.
 * Contains:
 * - ticket label data types
 * - ticket label change data types and associated data
 * - expectations for use with ticket labels and ticket label changes
 */

// TicketLabelImport holds the data on a label associated with an imported ticket label
type TicketLabelImport struct {
	tracName          string
	giteaLabelName    string
	giteaLabelID      int64
	giteaIssueLabelID int64
}

func createTicketLabelImport(prefix string, ticketLabelMap map[string]string) *TicketLabelImport {
	ticketLabel := TicketLabelImport{
		tracName:          prefix + "-name",
		giteaLabelName:    prefix + "-label",
		giteaLabelID:      allocateID(),
		giteaIssueLabelID: allocateID(),
	}

	ticketLabelMap[ticketLabel.tracName] = ticketLabel.giteaLabelName
	return &ticketLabel
}

var (
	componentChangeAuthor  *TicketUserImport
	priorityChangeAuthor   *TicketUserImport
	resolutionChangeAuthor *TicketUserImport
	severityChangeAuthor   *TicketUserImport
	typeChangeAuthor       *TicketUserImport
	versionChangeAuthor    *TicketUserImport
)

func setUpTicketLabelChangeUsers(t *testing.T) {
	componentChangeAuthor = createTicketUserImport("trac-component-change-author", "gitea-component-change-author")
	priorityChangeAuthor = createTicketUserImport("trac-priority-change-author", "gitea-priority-change-author")
	resolutionChangeAuthor = createTicketUserImport("trac-resolution-change-author", "gitea-resolution-change-author")
	severityChangeAuthor = createTicketUserImport("trac-severity-change-author", "gitea-severity-change-author")
	typeChangeAuthor = createTicketUserImport("trac-type-change-author", "gitea-type-change-author")
	versionChangeAuthor = createTicketUserImport("trac-version-change-author", "gitea-version-change-author")
}

func createLabelTicketChangeImport(author *TicketUserImport, tracChangeType trac.TicketChangeType, prevLabel *TicketLabelImport, label *TicketLabelImport) *TicketChangeImport {
	return &TicketChangeImport{
		tracChangeType: tracChangeType,
		issueCommentID: allocateID(),
		author:         author,
		prevLabel:      prevLabel,
		label:          label,
		time:           allocateUnixTime(),
	}
}

var (
	componentAddTicketChange    *TicketChangeImport
	componentAmendTicketChange  *TicketChangeImport
	componentRemoveTicketChange *TicketChangeImport

	priorityAddTicketChange    *TicketChangeImport
	priorityAmendTicketChange  *TicketChangeImport
	priorityRemoveTicketChange *TicketChangeImport

	resolutionAddTicketChange    *TicketChangeImport
	resolutionAmendTicketChange  *TicketChangeImport
	resolutionRemoveTicketChange *TicketChangeImport

	severityAddTicketChange    *TicketChangeImport
	severityAmendTicketChange  *TicketChangeImport
	severityRemoveTicketChange *TicketChangeImport

	typeAddTicketChange    *TicketChangeImport
	typeAmendTicketChange  *TicketChangeImport
	typeRemoveTicketChange *TicketChangeImport

	versionAddTicketChange    *TicketChangeImport
	versionAmendTicketChange  *TicketChangeImport
	versionRemoveTicketChange *TicketChangeImport
)

func setUpTicketLabelChanges(t *testing.T) {
	setUpTicketLabelChangeUsers(t)

	componentAddTicketChange = createLabelTicketChangeImport(componentChangeAuthor, trac.TicketComponentChange, nil, componentLabel2)
	componentAmendTicketChange = createLabelTicketChangeImport(componentChangeAuthor, trac.TicketComponentChange, componentLabel1, componentLabel2)
	componentRemoveTicketChange = createLabelTicketChangeImport(componentChangeAuthor, trac.TicketComponentChange, componentLabel1, nil)

	priorityAddTicketChange = createLabelTicketChangeImport(priorityChangeAuthor, trac.TicketPriorityChange, nil, priorityLabel2)
	priorityAmendTicketChange = createLabelTicketChangeImport(priorityChangeAuthor, trac.TicketPriorityChange, priorityLabel1, priorityLabel2)
	priorityRemoveTicketChange = createLabelTicketChangeImport(priorityChangeAuthor, trac.TicketPriorityChange, priorityLabel1, nil)

	resolutionAddTicketChange = createLabelTicketChangeImport(resolutionChangeAuthor, trac.TicketResolutionChange, nil, resolutionLabel2)
	resolutionAmendTicketChange = createLabelTicketChangeImport(resolutionChangeAuthor, trac.TicketResolutionChange, resolutionLabel1, resolutionLabel2)
	resolutionRemoveTicketChange = createLabelTicketChangeImport(resolutionChangeAuthor, trac.TicketResolutionChange, resolutionLabel1, nil)

	severityAddTicketChange = createLabelTicketChangeImport(severityChangeAuthor, trac.TicketSeverityChange, nil, severityLabel2)
	severityAmendTicketChange = createLabelTicketChangeImport(severityChangeAuthor, trac.TicketSeverityChange, severityLabel1, severityLabel2)
	severityRemoveTicketChange = createLabelTicketChangeImport(severityChangeAuthor, trac.TicketSeverityChange, severityLabel1, nil)

	typeAddTicketChange = createLabelTicketChangeImport(typeChangeAuthor, trac.TicketTypeChange, nil, typeLabel2)
	typeAmendTicketChange = createLabelTicketChangeImport(typeChangeAuthor, trac.TicketTypeChange, typeLabel1, typeLabel2)
	typeRemoveTicketChange = createLabelTicketChangeImport(typeChangeAuthor, trac.TicketTypeChange, typeLabel1, nil)

	versionAddTicketChange = createLabelTicketChangeImport(versionChangeAuthor, trac.TicketVersionChange, nil, versionLabel2)
	versionAmendTicketChange = createLabelTicketChangeImport(versionChangeAuthor, trac.TicketVersionChange, versionLabel1, versionLabel2)
	versionRemoveTicketChange = createLabelTicketChangeImport(versionChangeAuthor, trac.TicketVersionChange, versionLabel1, nil)
}

func expectLabelCreation(t *testing.T, label *TicketLabelImport) {
	mockGiteaAccessor.
		EXPECT().
		AddLabel(gomock.Eq(label.giteaLabelName), gomock.Any()).
		Return(label.giteaLabelID, nil)
}

func expectLabelRetrieval(t *testing.T, label *TicketLabelImport) {
	mockGiteaAccessor.
		EXPECT().
		GetLabelID(gomock.Eq(label.giteaLabelName)).
		Return(label.giteaLabelID, nil)
}

func expectIssueLabelCreation(t *testing.T, ticket *TicketImport, ticketLabel *TicketLabelImport) {
	// expect creation of underlying label first
	expectLabelCreation(t, ticketLabel)

	mockGiteaAccessor.
		EXPECT().
		AddIssueLabel(gomock.Eq(ticket.issueID), gomock.Eq(ticketLabel.giteaLabelID)).
		Return(ticketLabel.giteaIssueLabelID, nil)
}

func expectIssueCommentCreationForLabelChange(t *testing.T, ticket *TicketImport, ticketLabelChange *TicketChangeImport, label *TicketLabelImport, isAdd bool) {
	expectLabelRetrieval(t, label)

	mockGiteaAccessor.
		EXPECT().
		AddIssueComment(gomock.Eq(ticket.issueID), gomock.Any()).
		DoAndReturn(func(issueID int64, issueComment *gitea.IssueComment) (int64, error) {
			assertEquals(t, issueComment.CommentType, gitea.LabelIssueCommentType)
			assertEquals(t, issueComment.AuthorID, ticketLabelChange.author.giteaUserID)
			assertEquals(t, issueComment.LabelID, label.giteaLabelID)
			if isAdd {
				assertEquals(t, issueComment.Text, "1")
			}
			assertEquals(t, issueComment.Time, ticketLabelChange.time)
			return ticketLabelChange.issueCommentID, nil
		})

	if ticketLabelChange.author.giteaUser != "" {
		expectIssueParticipantToBeAdded(t, ticket, ticketLabelChange.author)
	}
}

func expectAllTicketLabelActions(t *testing.T, ticket *TicketImport, ticketLabelChange *TicketChangeImport) {
	// expect to lookup Gitea equivalent of author of Trac ticket change
	expectUserLookup(t, ticketLabelChange.author)

	// expect issue comment to remove previous label (if any) and issue comment to add new label (if any)
	if ticketLabelChange.prevLabel != nil {
		expectIssueCommentCreationForLabelChange(t, ticket, ticketLabelChange, ticketLabelChange.prevLabel, false)
	}
	if ticketLabelChange.label != nil {
		expectIssueCommentCreationForLabelChange(t, ticket, ticketLabelChange, ticketLabelChange.label, true)
	}
}
