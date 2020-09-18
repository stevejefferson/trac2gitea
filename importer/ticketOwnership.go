// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importOwnershipIssueComment imports a Trac ticket ownership change as a Gitea issue assignee change, returns id of created Gitea issue comment or gitea.NullID if cannot create comment
func (importer *Importer) importOwnershipIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string) (int64, error) {
	issueComment, err := importer.createIssueComment(issueID, change, userMap)
	if err != nil {
		return gitea.NullID, err
	}

	issueComment.CommentType = gitea.AssigneeIssueCommentType

	prevOwnerID := gitea.NullID
	prevOwnerName := change.OldValue
	if prevOwnerName != "" {
		prevOwnerID, err = importer.getUserID(prevOwnerName, userMap)
		if err != nil {
			return gitea.NullID, err
		}
		if prevOwnerID == gitea.NullID {
			return gitea.NullID, nil // cannot map user onto Gitea
		}
	}

	assigneeID := gitea.NullID
	removedAssigneeID := gitea.NullID
	ownerName := change.NewValue
	if ownerName != "" {
		assigneeID, err = importer.getUserID(ownerName, userMap)
		if err != nil {
			return gitea.NullID, err
		}
		if assigneeID == gitea.NullID {
			return gitea.NullID, nil // cannot map user onto Gitea
		}
	} else {
		removedAssigneeID = prevOwnerID
	}

	issueComment.AssigneeID = assigneeID
	issueComment.RemovedAssigneeID = removedAssigneeID
	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return gitea.NullID, err
	}

	return issueCommentID, nil
}
