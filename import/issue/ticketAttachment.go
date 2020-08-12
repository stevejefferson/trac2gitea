package issue

import (
	"fmt"
	"os"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
)

func (importer *Importer) importTicketAttachment(issueID int64, ticketID int64, time int64, size int64, author string, attachmentName string, desc string) string {
	comment := fmt.Sprintf("**Attachment** %s (%d bytes) added\n\n%s", attachmentName, size, desc)
	commentID := importer.importTicketComment(issueID, ticketID, time, author, comment)

	tracPath := importer.tracAccessor.GetAttachmentPath(ticketID, attachmentName)
	_, err := os.Stat(tracPath)
	if err != nil {
		log.Fatal(err)
	}
	elems := strings.Split(tracPath, "/")
	tracDir := elems[len(elems)-2]
	tracFile := elems[len(elems)-1]

	// 78ac is l33t for trac (horrible, I know)
	uuid := fmt.Sprintf("000078ac-%s-%s-%s-%s",
		tracDir[0:4], tracDir[4:8], tracDir[8:12],
		tracFile[0:12])

	// TODO: use a different uuid if file exists ?
	// TODO: avoid inserting record if uuid exist !
	importer.giteaAccessor.AddAttachment(uuid, issueID, commentID, attachmentName, tracPath, time)

	return uuid
}

func (importer *Importer) importTicketAttachments(ticketID int64, issueID int64, created int64) int64 {
	lastUpdate := created

	importer.tracAccessor.GetAttachments(ticketID, func(ticketID int64, time int64, size int64, author string, filename string, description string) {
		if lastUpdate > time {
			lastUpdate = time
		}
		importer.importTicketAttachment(issueID, ticketID, time, size, author, filename, description)
	})

	return lastUpdate
}
