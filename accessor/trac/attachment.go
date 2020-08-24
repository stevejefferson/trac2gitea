// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path"

	"github.com/stevejefferson/trac2gitea/log"
)

func encodeSha1(str string) string {
	// Encode string to sha1 hex value.
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (accessor *DefaultAccessor) getAttachmentPath(idStr string, attachmentName string, attachmentType string) string {
	idHash := encodeSha1(idStr)
	idHashSub := idHash[0:3]

	pathFile := encodeSha1(attachmentName)
	pathExt := path.Ext(attachmentName)

	return accessor.GetFullPath("files", "attachments", attachmentType, idHashSub, idHash, pathFile+pathExt)
}

// GetTicketAttachmentPath retrieves the path to a named attachment to a Trac ticket.
func (accessor *DefaultAccessor) GetTicketAttachmentPath(ticketID int64, attachmentName string) string {
	ticketIDStr := fmt.Sprintf("%d", ticketID)
	return accessor.getAttachmentPath(ticketIDStr, attachmentName, "ticket")
}

// GetWikiAttachmentPath retrieves the path to a named attachment to a Trac wiki page.
func (accessor *DefaultAccessor) GetWikiAttachmentPath(wikiPage string, attachmentName string) string {
	return accessor.getAttachmentPath(wikiPage, attachmentName, "wiki")
}

// GetAttachments retrieves all attachments for a given Trac ticket, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetAttachments(
	ticketID int64,
	handlerFn func(ticketID int64, time int64, size int64, author string, filename string, description string) error) error {
	rows, err := accessor.db.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, filename, description, size
			FROM attachment
			WHERE type = 'ticket' AND id = $1
			ORDER BY time asc`, ticketID)
	if err != nil {
		log.Error("Problem retrieving attachment for trac ticket %d: %v\n", ticketID, err)
		return err
	}

	for rows.Next() {
		var time, size int64
		var author, filename, description string
		if err := rows.Scan(&time, &author, &filename, &description, &size); err != nil {
			log.Error("Problem extracting attachment data for trac ticket %d: %v\n", ticketID, err)
			return err
		}

		err = handlerFn(ticketID, time, size, author, filename, description)
		if err != nil {
			return err
		}
	}

	return nil
}
