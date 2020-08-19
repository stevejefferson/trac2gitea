package issue

import (
	"fmt"
	"os"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
)

// importTicketAttachment imports a single ticket attachment from Trac into Gitea, returns UUID if newly-created attachment or "" if attachment already existed
func (importer *Importer) importTicketAttachment(issueID int64, ticketID int64, time int64, size int64, author string, attachmentName string, desc string) string {
	comment := fmt.Sprintf("**Attachment** %s (%d bytes) added\n\n%s", attachmentName, size, desc)
	commentID := importer.importTicketComment(issueID, ticketID, time, author, comment)

	tracPath := importer.tracAccessor.GetTicketAttachmentPath(ticketID, attachmentName)
	_, err := os.Stat(tracPath)
	if err != nil {
		log.Fatal(err)
	}
	elems := strings.Split(tracPath, "/")
	tracDir := elems[len(elems)-2]
	tracFile := elems[len(elems)-1]

	// use '78ac' to identify Trac UUIDs (from trac2gogs)
	uuid := fmt.Sprintf("000078ac-%s-%s-%s-%s",
		tracDir[0:4], tracDir[4:8], tracDir[8:12],
		tracFile[0:12])

	existingUUID := importer.giteaAccessor.GetAttachmentUUID(issueID, attachmentName)
	if existingUUID != "" {
		if existingUUID == uuid {
			log.Debugf("Attachment %s, (uuid=\"%s\") already exists for issue %d - skipping...\n", attachmentName, uuid, issueID)
		} else {
			log.Warnf("Attachment %s already exists for issue %d but under uuid \"%s\" (expecting \"%s\") - skipping...\n", attachmentName, issueID, existingUUID, uuid)
		}
		return ""
	}

	importer.giteaAccessor.AddAttachment(uuid, issueID, commentID, attachmentName, tracPath, time)

	return uuid
}

func (importer *Importer) importTicketAttachments(ticketID int64, issueID int64, created int64) int64 {
	lastUpdate := created

	importer.tracAccessor.GetAttachments(ticketID, func(ticketID int64, time int64, size int64, author string, filename string, description string) {
		uuid := importer.importTicketAttachment(issueID, ticketID, time, size, author, filename, description)
		if uuid != "" && lastUpdate > time {
			lastUpdate = time
		}
	})

	return lastUpdate
}
