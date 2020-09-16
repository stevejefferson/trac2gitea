// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importCommentIssueComment imports a Trac ticket comment into Gitea, returns id of created Gitea issue comment or -1 if cannot create comment
func (importer *Importer) importCommentIssueComment(issueID int64, change *trac.TicketChange, issueComment *gitea.IssueComment, userMap map[string]string) (int64, error) {
	issueComment.CommentType = gitea.CommentIssueCommentType
	issueComment.Text = importer.markdownConverter.TicketConvert(change.TicketID, change.Comment.Text)

	issueCommentID, err := importer.giteaAccessor.AddIssueComment(issueID, issueComment)
	if err != nil {
		return -1, err
	}

	return issueCommentID, nil
}
