package issue

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func (importer *Importer) importTicketAttachment(issueID int64, ticketID int64, time int64, size int64, author string, fname string, desc string) string {
	comment := fmt.Sprintf("**Attachment** %s (%d bytes) added\n\n%s",
		fname, size, desc)
	commentID := importer.importTicketComment(issueID, ticketID, time, author, comment)

	tracPath := importer.tracAccessor.AttachmentPath(ticketID, fname)
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
	importer.giteaAccessor.AddAttachment(uuid, issueID, commentID, fname, time)

	giteaRelPath := importer.giteaAccessor.AttachmentRelativePath(uuid)
	importer.giteaAccessor.CopyFile(tracPath, giteaRelPath)

	return uuid
}

func (importer *Importer) importTicketAttachments(id int64, issueID int64, created int64) int64 {
	rows := importer.tracAccessor.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, filename, description, size
			FROM attachment
  			WHERE type = 'ticket' AND id = $1
			ORDER BY time asc`, id)

	lastUpdate := created
	for rows.Next() {
		var time, size int64
		var author, fname, desc string
		if err := rows.Scan(&time, &author, &fname, &desc, &size); err != nil {
			log.Fatal(err)
		}

		fmt.Println(" adding attachment by", author)
		if lastUpdate > time {
			lastUpdate = time
		}
		importer.importTicketAttachment(issueID, id, time, size, author, fname, desc)
	}

	return lastUpdate
}
