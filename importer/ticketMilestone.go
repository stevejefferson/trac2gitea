// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importMilestoneIssueComment imports a Trac ticket milestone change into Gitea, returns id of created Gitea issue comment or -1 if cannot create comment
func (importer *Importer) importMilestoneIssueComment(issueID int64, change *trac.TicketChange, issueComment *gitea.IssueComment) (int64, error) {
	var err error
	var prevMilestoneID = int64(0)
	prevMilestone := change.OldValue
	if prevMilestone != "" {
		prevMilestoneID, err = importer.giteaAccessor.GetMilestoneID(prevMilestone)
		if err != nil {
			return -1, err
		}

		if prevMilestoneID == -1 {
			// unrecognised milestone - e.g. milestone deleted from Trac for which no definition remains
			// - best solution is just to pretend it did not exist
			prevMilestoneID = 0
		}
	}

	var milestoneID = int64(0)
	milestone := change.NewValue
	if milestone != "" {
		milestoneID, err = importer.giteaAccessor.GetMilestoneID(milestone)
		if err != nil {
			return -1, err
		}

		if milestoneID == -1 {
			// unrecognised milestone - e.g. milestone deleted from Trac for which no definition remains
			// - cannot cope with this so ignore entire milestone change
			return -1, nil
		}
	}

	issueComment.CommentType = gitea.MilestoneIssueCommentType
	issueComment.OldMilestoneID = prevMilestoneID
	issueComment.MilestoneID = milestoneID
	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return -1, err
	}
	return issueCommentID, nil
}
