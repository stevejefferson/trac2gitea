// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"time"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
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
	authorID, err := importer.getUser(tracComment.Author, userMap)
	if err != nil {
		return -1, err
	}
	if authorID == -1 {
		authorID = importer.defaultUserID
	}

	convertedText := importer.markdownConverter.TicketConvert(tracComment.TicketID, tracComment.Text)
	commentID, err := importer.giteaAccessor.GetTimedIssueCommentID(issueID, tracComment.Time)
	if err != nil {
		return -1, err
	}
	if commentID != -1 {
		log.Debug("comment for issue %d, created at %s already exists - skipping...", issueID, time.Unix(tracComment.Time, 0))
		return -1, nil
	}

	giteaComment := gitea.IssueComment{AuthorID: authorID, OriginalAuthorID: 0, OriginalAuthorName: tracComment.Author, Text: convertedText, Time: tracComment.Time}
	commentID, err = importer.giteaAccessor.AddIssueComment(issueID, &giteaComment)
	if err != nil {
		return -1, err
	}

	log.Debug("issue %d: added comment \"%s\" (id %d)", issueID, truncateString(convertedText, 20), commentID)

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
