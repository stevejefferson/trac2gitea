// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importMilestoneIssueComment imports a Trac ticket milestone change into Gitea, returns id of created Gitea issue comment or gitea.NullID if cannot create comment
func (importer *Importer) importMilestoneIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string) (int64, error) {
	issueComment, err := importer.createIssueComment(issueID, change, userMap)
	if err != nil {
		return gitea.NullID, err
	}

	var oldMilestoneID = gitea.NullID
	oldMilestone := change.OldValue
	if oldMilestone != "" {
		oldMilestoneID, err = importer.giteaAccessor.GetMilestoneID(oldMilestone)
		if err != nil {
			return gitea.NullID, err
		}
	}

	var milestoneID = gitea.NullID
	milestone := change.NewValue
	if milestone != "" {
		milestoneID, err = importer.giteaAccessor.GetMilestoneID(milestone)
		if err != nil {
			return gitea.NullID, err
		}

		if milestoneID == gitea.NullID {
			// unrecognised milestone - e.g. milestone deleted from Trac for which no definition remains
			// - cannot cope with this so ignore entire milestone change
			return gitea.NullID, nil
		}
	}

	issueComment.CommentType = gitea.MilestoneIssueCommentType
	issueComment.OldMilestoneID = oldMilestoneID
	issueComment.MilestoneID = milestoneID
	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return gitea.NullID, err
	}
	return issueCommentID, nil
}
