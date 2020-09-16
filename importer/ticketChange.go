// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importTicketChange imports a single ticket change from Trac to Gitea, returns ID of created Gitea comment or -1 if comment already exists
func (importer *Importer) importTicketChange(issueID int64, change *trac.TicketChange, userMap map[string]string) (int64, error) {
	originalAuthorName := ""
	authorID, err := importer.getUser(change.Author, userMap)
	if err != nil {
		return -1, err
	}
	if authorID != -1 {
		// change author has a Gitea mapping: make this user  a participant in issue
		err = importer.giteaAccessor.AddIssueParticipant(issueID, authorID)
		if err != nil {
			return -1, err
		}
	} else {
		// change author cannot be mapped onto Gitea: use default user as author but record original Trac user on the change
		authorID = importer.defaultAuthorID
		originalAuthorName = change.Author
	}

	// perform change-specific issue operations
	issueComment := gitea.IssueComment{
		AssigneeID:         0,
		AuthorID:           authorID,
		OriginalAuthorID:   0,
		OriginalAuthorName: originalAuthorName,
		RemovedAssigneeID:  0,
		Text:               "",
		Time:               change.Time,
	}
	var issueCommentID int64
	switch change.ChangeType {
	case trac.TicketCommentChange:
		issueCommentID, err = importer.importCommentIssueComment(issueID, change, &issueComment, userMap)
	case trac.TicketMilestoneChange:
		issueCommentID, err = importer.importMilestoneIssueComment(issueID, change, &issueComment)
	case trac.TicketOwnerChange:
		issueCommentID, err = importer.importOwnershipIssueComment(issueID, change, &issueComment, userMap)
	case trac.TicketStatusChange:
		issueCommentID, err = importer.importStatusChangeIssueComment(issueID, change, &issueComment)
	}
	if err != nil {
		return -1, err
	}

	return issueCommentID, nil
}

func (importer *Importer) importTicketChanges(ticketID int64, issueID int64, lastUpdate int64, userMap map[string]string) (int64, error) {
	commentLastUpdate := lastUpdate
	err := importer.tracAccessor.GetTicketChanges(ticketID, func(change *trac.TicketChange) error {
		commentID, err := importer.importTicketChange(issueID, change, userMap)
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
