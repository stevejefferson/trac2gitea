// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"

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
func (importer *Importer) importTicketComment(issueID int64, ticketID int64, time int64, author string, comment string) (int64, error) {
	authorID, _, err := importer.getUser(author)
	if err != nil {
		return -1, err
	}

	tracDetails := fmt.Sprintf("original comment by %s", author)
	markdownConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	convertedComment := markdownConverter.Convert(comment)
	fullComment := addTracContext(tracDetails, time, convertedComment)
	commentID, err := importer.giteaAccessor.GetCommentID(issueID, fullComment)
	if err != nil {
		return -1, err
	}

	truncatedComment := truncateString(comment, 16) // used for diagnostics
	if commentID != -1 {
		log.Debug("comment \"%s\" for issue %d already exists - skipping...", truncatedComment, issueID)
		return -1, nil
	}

	commentID, err = importer.giteaAccessor.AddComment(issueID, authorID, fullComment, time)
	if err != nil {
		return -1, err
	}

	log.Debug("issue %d: added comment \"%s\" (id %d)", issueID, truncatedComment, commentID)

	return commentID, nil
}

func (importer *Importer) importTicketComments(ticketID int64, issueID int64, lastUpdate int64) error {
	err := importer.tracAccessor.GetComments(ticketID, func(ticketID int64, time int64, author string, comment string) error {
		commentID, err := importer.importTicketComment(issueID, ticketID, time, author, comment)
		if err != nil {
			return err
		}

		if commentID != -1 && lastUpdate > time {
			lastUpdate = time
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Update issue modification time
	return importer.giteaAccessor.SetIssueUpdateTime(issueID, lastUpdate)
}
