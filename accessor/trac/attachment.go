package trac

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path"

	"stevejefferson.co.uk/trac2gitea/log"
)

func encodeSha1(str string) string {
	// Encode string to sha1 hex value.
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetAttachmentPath retrieves the path to a named attachment to a Trac ticket.
func (accessor *Accessor) GetAttachmentPath(ticketID int64, name string) string {
	ticketDir := encodeSha1(fmt.Sprintf("%d", ticketID))
	ticketSub := ticketDir[0:3]

	pathFile := encodeSha1(name)
	pathExt := path.Ext(name)

	return accessor.GetFullPath("attachments", "ticket", ticketSub, ticketDir, pathFile+pathExt)
}

// GetAttachments retrieves all attachments for a given Trac ticket, passing data from each one to the provided "handler" function.
func (accessor *Accessor) GetAttachments(ticketID int64, handlerFn func(ticketID int64, time int64, size int64, author string, filename string, description string)) {
	rows, err := accessor.db.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, filename, description, size
			FROM attachment
			WHERE type = 'ticket' AND id = $1
			ORDER BY time asc`, ticketID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var time, size int64
		var author, filename, description string
		if err := rows.Scan(&time, &author, &filename, &description, &size); err != nil {
			log.Fatal(err)
		}

		handlerFn(ticketID, time, size, author, filename, description)
	}
}
