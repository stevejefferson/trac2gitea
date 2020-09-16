// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importStatusChangeIssueComment imports a Trac ticket status change into Gitea, returns id of created Gitea issue comment or -1 if cannot create comment
func (importer *Importer) importStatusChangeIssueComment(issueID int64, change *trac.TicketChange, userMap map[string]string) (int64, error) {
	issueComment, err := importer.createIssueComment(issueID, change, userMap)
	if err != nil {
		return -1, err
	}

	var giteaCommentType gitea.IssueCommentType
	switch change.NewValue {
	case trac.TicketStatusClosed:
		giteaCommentType = gitea.CloseIssueCommentType
	case trac.TicketStatusReopened:
		giteaCommentType = gitea.ReopenIssueCommentType
	default:
		return -1, nil
	}

	issueComment.CommentType = giteaCommentType
	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return -1, err
	}
	return issueCommentID, nil
}
