// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
)

// importTicketLabel imports a single issue label from Trac into Gitea, returns id of created issue label or gitea.NullID if issue label already exists
func (importer *Importer) importTicketLabel(issueID int64, tracName string, labelMap map[string]string) (int64, error) {
	labelID, err := importer.getLabelID(tracName, labelMap)
	if err != nil {
		return gitea.NullID, err
	}
	if labelID == gitea.NullID {
		return gitea.NullID, nil
	}

	issueLabelID, err := importer.giteaAccessor.AddIssueLabel(issueID, labelID)
	if err != nil {
		return gitea.NullID, err
	}

	log.Debug("created issue label (id %d) for issue %d, label %d", issueLabelID, issueID, labelID)

	return issueLabelID, nil
}

// addLabelChangeIssueComment adds a single label change issue comment into Gitea, returns id of created Gitea issue comment or gitea.NullID if cannot create comment
func (importer *Importer) addLabelChangeIssueComment(issueID int64, change *trac.TicketChange, labelName string, isAdd bool, userMap map[string]string, labelMap map[string]string) (int64, error) {
	var issueCommentID int64

	labelID, err := importer.getLabelID(labelName, labelMap)
	if err != nil {
		return gitea.NullID, err
	}

	if labelID != gitea.NullID {
		issueComment, err := importer.createIssueComment(issueID, change, userMap)
		if err != nil {
			return gitea.NullID, err
		}
		issueComment.CommentType = gitea.LabelIssueCommentType
		issueComment.LabelID = labelID
		if isAdd {
			issueComment.Text = "1"
		}
		issueCommentID, err = importer.giteaAccessor.AddIssueComment(issueID, issueComment)
		if err != nil {
			return gitea.NullID, err
		}
	}

	return issueCommentID, nil
}

// importLabelChangeIssueComment imports a Trac ticket label change into Gitea, returns id of created Gitea issue comment or NullID if cannot create comment
func (importer *Importer) importLabelChangeIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string, labelMap map[string]string) (int64, error) {
	var err error
	var issueCommentID int64

	// if we have previous label name and it maps on to Gitea then this is a label removal event and we must generate an issue comment for that
	prevLabelName := change.OldValue
	if prevLabelName != "" {
		issueCommentID, err = importer.addLabelChangeIssueComment(issueID, change, prevLabelName, false, userMap, labelMap)
		if err != nil {
			return gitea.NullID, err
		}
	}

	// if we have a new label name and it maps onto Gitea then this is a label addition event and we generate an issue comment for that
	labelName := change.NewValue
	if labelName != "" {
		issueCommentID, err = importer.addLabelChangeIssueComment(issueID, change, labelName, true, userMap, labelMap)
		if err != nil {
			return gitea.NullID, err
		}
	}

	return issueCommentID, nil
}
