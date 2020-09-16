// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// createIssueComment creates a basic Gitea IssueComment structure to be populated by individual ticket change import functions
func (importer *Importer) createIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string) (*gitea.IssueComment, error) {
	originalAuthorName := ""
	authorID, err := importer.getUserID(change.Author, userMap)
	if err != nil {
		return nil, err
	}
	if authorID != -1 {
		// change author has a Gitea mapping: make this user  a participant in issue
		err = importer.giteaAccessor.AddIssueParticipant(issueID, authorID)
		if err != nil {
			return nil, err
		}
	} else {
		// change author cannot be mapped onto Gitea: use default user as author but record original Trac user on the change
		authorID = importer.defaultAuthorID
		originalAuthorName = change.Author
	}

	// perform change-specific issue operations
	issueComment := gitea.IssueComment{
		AuthorID:           authorID,
		OriginalAuthorID:   0,
		OriginalAuthorName: originalAuthorName,
		LabelID:            0,
		OldMilestoneID:     0,
		MilestoneID:        0,
		AssigneeID:         0,
		RemovedAssigneeID:  0,
		OldTitle:           "",
		Title:              "",
		Text:               "",
		Time:               change.Time,
	}

	return &issueComment, nil
}

// importTicketChange imports a single ticket change from Trac to Gitea, returns ID of created Gitea comment or -1 if comment already exists
func (importer *Importer) importTicketChange(
	issueID int64,
	change *trac.TicketChange,
	userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) (int64, error) {
	var issueCommentID int64
	var err error

	switch change.ChangeType {
	case trac.TicketCommentChange:
		issueCommentID, err = importer.importCommentIssueComment(issueID, change, userMap)
	case trac.TicketComponentChange:
		issueCommentID, err = importer.importLabelChangeIssueComment(issueID, change, userMap, componentMap)
	case trac.TicketMilestoneChange:
		issueCommentID, err = importer.importMilestoneIssueComment(issueID, change, userMap)
	case trac.TicketOwnerChange:
		issueCommentID, err = importer.importOwnershipIssueComment(issueID, change, userMap)
	case trac.TicketPriorityChange:
		issueCommentID, err = importer.importLabelChangeIssueComment(issueID, change, userMap, priorityMap)
	case trac.TicketResolutionChange:
		issueCommentID, err = importer.importLabelChangeIssueComment(issueID, change, userMap, resolutionMap)
	case trac.TicketSeverityChange:
		issueCommentID, err = importer.importLabelChangeIssueComment(issueID, change, userMap, severityMap)
	case trac.TicketTypeChange:
		issueCommentID, err = importer.importLabelChangeIssueComment(issueID, change, userMap, typeMap)
	case trac.TicketStatusChange:
		issueCommentID, err = importer.importStatusChangeIssueComment(issueID, change, userMap)
	case trac.TicketSummaryChange:
		issueCommentID, err = importer.importSummaryChangeIssueComment(issueID, change, userMap)
	case trac.TicketVersionChange:
		issueCommentID, err = importer.importLabelChangeIssueComment(issueID, change, userMap, versionMap)
	}
	if err != nil {
		return -1, err
	}

	return issueCommentID, nil
}

func (importer *Importer) importTicketChanges(
	ticketID int64,
	issueID int64,
	lastUpdate int64,
	userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) (int64, error) {
	commentLastUpdate := lastUpdate
	err := importer.tracAccessor.GetTicketChanges(ticketID, func(change *trac.TicketChange) error {
		commentID, err := importer.importTicketChange(issueID, change, userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap)
		if err != nil {
			return err
		}
		if commentID != -1 && commentLastUpdate < change.Time {
			commentLastUpdate = change.Time
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return commentLastUpdate, nil
}
