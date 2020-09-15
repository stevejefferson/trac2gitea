// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"strings"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importTicketAttachment imports a single ticket attachment from Trac into Gitea, returns UUID if newly-created attachment or "" if attachment already existed
func (importer *Importer) importTicketAttachment(issueID int64, tracAttachment *trac.TicketAttachment, userMap map[string]string) (string, error) {
	commentText := fmt.Sprintf("**Attachment** %s (%d bytes) added\n\n%s", tracAttachment.FileName, tracAttachment.Size, tracAttachment.Description)

	tracComment := trac.TicketComment{Text: commentText}
	tracChange := trac.TicketChange{
		TicketID:   tracAttachment.TicketID,
		Author:     tracAttachment.Author,
		Time:       tracAttachment.Time,
		ChangeType: trac.TicketCommentChange,
		Comment:    &tracComment}
	commentID, err := importer.importTicketChange(issueID, &tracChange, userMap)
	if err != nil {
		return "", err
	}

	tracPath := importer.tracAccessor.GetTicketAttachmentPath(tracAttachment)
	elems := strings.Split(tracPath, "/")
	tracDir := elems[len(elems)-2]
	tracFile := elems[len(elems)-1]

	if len(tracDir) < 12 {
		return "", fmt.Errorf("Malformed trac attachment path \"%s\" - \"%s\" should be trac UUID (at least 12 chars)", tracPath, tracDir)
	}

	// use '78ac' to identify Trac UUIDs (from trac2gogs)
	uuid := fmt.Sprintf("000078ac-%s-%s-%s-%s",
		tracDir[0:4], tracDir[4:8], tracDir[8:12],
		tracFile[0:12])

	giteaAttachment := gitea.IssueAttachment{UUID: uuid, CommentID: commentID, FileName: tracAttachment.FileName, Time: tracAttachment.Time}
	_, err = importer.giteaAccessor.AddIssueAttachment(issueID, &giteaAttachment, tracPath)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (importer *Importer) importTicketAttachments(ticketID int64, issueID int64, lastUpdate int64, userMap map[string]string) (int64, error) {
	attachmentLastUpdate := lastUpdate

	err := importer.tracAccessor.GetTicketAttachments(ticketID, func(attachment *trac.TicketAttachment) error {
		uuid, err := importer.importTicketAttachment(issueID, attachment, userMap)
		if err != nil {
			return err
		}

		if uuid != "" && attachmentLastUpdate < attachment.Time {
			attachmentLastUpdate = attachment.Time
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return attachmentLastUpdate, nil
}
