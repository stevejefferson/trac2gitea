// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"
	"strings"

	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
)

// importTicketAttachment imports a single ticket attachment from Trac into Gitea, returns UUID if newly-created attachment or "" if attachment already existed
func (importer *Importer) importTicketAttachment(issueID int64, attachment *trac.TicketAttachment, userMap map[string]string) (string, error) {
	commentText := fmt.Sprintf("**Attachment** %s (%d bytes) added\n\n%s", attachment.FileName, attachment.Size, attachment.Description)

	tracComment := trac.TicketComment{TicketID: attachment.TicketID, Time: attachment.Time, Author: attachment.Author, Text: commentText}
	commentID, err := importer.importTicketComment(issueID, &tracComment, userMap)
	if err != nil {
		return "", err
	}

	tracPath := importer.tracAccessor.GetTicketAttachmentPath(attachment)
	elems := strings.Split(tracPath, "/")
	tracDir := elems[len(elems)-2]
	tracFile := elems[len(elems)-1]

	// use '78ac' to identify Trac UUIDs (from trac2gogs)
	uuid := fmt.Sprintf("000078ac-%s-%s-%s-%s",
		tracDir[0:4], tracDir[4:8], tracDir[8:12],
		tracFile[0:12])

	existingUUID, err := importer.giteaAccessor.GetAttachmentUUID(issueID, attachment.FileName)
	if err != nil {
		return "", err
	}

	if existingUUID != "" {
		if existingUUID == uuid {
			log.Debug("attachment %s, (uuid=\"%s\") already exists for issue %d - skipping...", attachment.FileName, uuid, issueID)
		} else {
			log.Warn("attachment %s already exists for issue %d but under uuid \"%s\" (expecting \"%s\") - skipping...",
				attachment.FileName, issueID, existingUUID, uuid)
		}
		return "", nil
	}

	_, err = importer.giteaAccessor.AddAttachment(uuid, issueID, commentID, attachment.FileName, tracPath, attachment.Time)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (importer *Importer) importTicketAttachments(ticketID int64, issueID int64, created int64, userMap map[string]string) (int64, error) {
	lastUpdate := created

	err := importer.tracAccessor.GetTicketAttachments(ticketID, func(attachment *trac.TicketAttachment) error {
		uuid, err := importer.importTicketAttachment(issueID, attachment, userMap)
		if err != nil {
			return err
		}

		if uuid != "" && lastUpdate < attachment.Time {
			lastUpdate = attachment.Time
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return lastUpdate, nil
}
