package issueimport

import (
	"fmt"
	"log"
	"strings"
)

func (importer *Importer) importTicketComment(issueID int64, ticketID int64, time int64, author, comment string) int64 {
	comment = importer.trac2MarkdownConverter.TicketConvert(comment, ticketID)

	// find users first, and tweak description to add missing users
	var header []string
	authorID := importer.giteaAccessor.GetUserID(author)
	if authorID == -1 {
		header = append(header, fmt.Sprintf("    Original comment by %s", author))
		authorID = importer.giteaAccessor.DefaultAuthorID
	}
	if len(header) > 0 {
		comment = fmt.Sprintf("%s\n\n%s", strings.Join(header, "\n"), comment)
	}

	return importer.giteaAccessor.AddComment(issueID, authorID, comment, time)
}

func (importer *Importer) importTicketComments(ticketID int64, issueID int64, lastUpdate int64) {
	rows := importer.tracAccessor.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, COALESCE(newvalue, '') newval
			FROM ticket_change where ticket = $1 AND field = 'comment' AND trim(COALESCE(newvalue, ''), ' ') != ''
			ORDER BY time asc`, ticketID)

	for rows.Next() {
		var time int64
		var author, comment string
		if err := rows.Scan(&time, &author, &comment); err != nil {
			log.Fatal(err)
		}
		fmt.Println(" adding comment by", author)
		if lastUpdate > time {
			lastUpdate = time
		}
		importer.importTicketComment(issueID, ticketID, time, author, comment)
	}

	// Update issue modification time
	importer.giteaAccessor.SetIssueUpdateTime(issueID, lastUpdate)
}
