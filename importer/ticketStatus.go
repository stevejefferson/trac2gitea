// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importStatusChangeIssueComment imports a Trac ticket status change into Gitea, returns id of created Gitea issue comment or -1 if cannot create comment
func (importer *Importer) importStatusChangeIssueComment(issueID int64, change *trac.TicketChange, issueComment *gitea.IssueComment) (int64, error) {
	// the only Trac status change that interests us is closing a ticket
	if change.NewValue != string(trac.TicketStatusClosed) {
		return -1, nil
	}

	issueComment.CommentType = gitea.CloseIssueCommentType
	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return -1, err
	}
	return issueCommentID, nil
}
