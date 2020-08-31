// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package data

import (
	"fmt"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
	"github.com/stevejefferson/trac2gitea/markdown"
)

func truncateString(str string, maxlen int) string {
	strLen := len(str)
	if strLen > maxlen {
		return str[0:maxlen] + "..."
	}
	return str
}

// importTicketComment imports a single ticket comment from Trac to Gitea, returns ID of created comment or -1 if comment already exists
func (importer *Importer) importTicketComment(issueID int64, tracComment *trac.TicketComment, userMap map[string]string) (int64, error) {
	authorID, _, err := importer.getUser(tracComment.Author, userMap)
	if err != nil {
		return -1, err
	}

	tracDetails := fmt.Sprintf("original comment by %s", tracComment.Author)
	context := markdown.ConversionContext{TicketID: tracComment.TicketID, WikiPage: ""}
	convertedText := importer.markdownConverter.Convert(&context, tracComment.Text)
	fullText := addTracContext(tracDetails, tracComment.Time, convertedText)
	commentID, err := importer.giteaAccessor.GetIssueCommentID(issueID, fullText)
	if err != nil {
		return -1, err
	}

	truncatedText := truncateString(tracComment.Text, 16) // used for diagnostics
	if commentID != -1 {
		log.Debug("comment \"%s\" for issue %d already exists - skipping...", truncatedText, issueID)
		return -1, nil
	}

	giteaComment := gitea.IssueComment{IssueID: issueID, AuthorID: authorID, Text: fullText, Time: tracComment.Time}
	commentID, err = importer.giteaAccessor.AddIssueComment(&giteaComment)
	if err != nil {
		return -1, err
	}

	log.Debug("issue %d: added comment \"%s\" (id %d)", issueID, truncatedText, commentID)

	return commentID, nil
}

func (importer *Importer) importTicketComments(ticketID int64, issueID int64, lastUpdate int64, userMap map[string]string) (int64, error) {
	commentLastUpdate := lastUpdate
	err := importer.tracAccessor.GetTicketComments(ticketID, func(comment *trac.TicketComment) error {
		commentID, err := importer.importTicketComment(issueID, comment, userMap)
		if err != nil {
			return err
		}

		if commentID != -1 && commentLastUpdate < comment.Time {
			commentLastUpdate = comment.Time
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return commentLastUpdate, nil
}
