// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"
	"strings"

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
func (importer *Importer) importTicketComment(issueID int64, ticketID int64, time int64, author, comment string) (int64, error) {
	markdownConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	comment = markdownConverter.Convert(comment)

	// find users first, and tweak description to add missing users
	var header []string
	authorID, err := importer.giteaAccessor.GetUserID(author)
	if err != nil {
		return -1, err
	}
	if authorID == -1 {
		header = append(header, fmt.Sprintf("    Original comment by %s", author))
		authorID = importer.giteaAccessor.GetDefaultAuthorID()
	}

	if len(header) > 0 {
		comment = fmt.Sprintf("%s\n\n%s", strings.Join(header, "\n"), comment)
	}

	truncatedComment := truncateString(comment, 16)
	commentID, err := importer.giteaAccessor.GetCommentID(issueID, comment)
	if err != nil {
		return -1, err
	}
	if commentID != -1 {
		log.Debug("Comment \"%s\" for issue %d already exists - skipping...\n", truncatedComment, issueID)
		return -1, nil
	}

	commentID, err = importer.giteaAccessor.AddComment(issueID, authorID, comment, time)
	if err != nil {
		return -1, err
	}

	log.Debug("Issue %d: added comment \"%s\" (id %d)\n", issueID, truncatedComment, commentID)

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
