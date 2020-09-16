// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importOwnershipIssueComment imports a Trac ticket ownership change as a Gitea issue assignee change, returns id of created Gitea issue comment or -1 if cannot create comment
func (importer *Importer) importOwnershipIssueComment(issueID int64, change *trac.TicketChange, issueComment *gitea.IssueComment, userMap map[string]string) (int64, error) {
	issueComment.CommentType = gitea.AssigneeIssueCommentType

	var err error
	prevOwnerID := int64(0)
	prevOwnerName := change.OldValue
	if prevOwnerName != "" {
		prevOwnerID, err = importer.getUser(prevOwnerName, userMap)
		if err != nil {
			return -1, err
		}
		if prevOwnerID == -1 {
			return -1, nil // cannot map user onto Gitea
		}
	}

	assigneeID := int64(0)
	removedAssigneeID := int64(0)
	ownerName := change.NewValue
	if ownerName != "" {
		assigneeID, err = importer.getUser(ownerName, userMap)
		if err != nil {
			return -1, err
		}
		if assigneeID == -1 {
			return -1, nil // cannot map user onto Gitea
		}
	} else {
		removedAssigneeID = prevOwnerID
	}

	issueComment.AssigneeID = assigneeID
	issueComment.RemovedAssigneeID = removedAssigneeID
	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return -1, err
	}

	// add or remove assignee for issue
	if ownerName != "" {
		err = importer.giteaAccessor.AddIssueAssignee(issueID, assigneeID)
	} else {
		err = importer.giteaAccessor.RemoveIssueAssignee(issueID, removedAssigneeID)
	}
	if err != nil {
		return -1, err
	}

	return issueCommentID, nil
}
