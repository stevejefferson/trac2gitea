package issue

import (
	"fmt"
	"strings"

	"stevejefferson.co.uk/trac2gitea/markdown"
)

func (importer *Importer) importTicketComment(issueID int64, ticketID int64, time int64, author, comment string) int64 {
	t2mConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	comment = t2mConverter.Convert(comment)

	// find users first, and tweak description to add missing users
	var header []string
	authorID := importer.giteaAccessor.GetUserID(author)
	if authorID == -1 {
		header = append(header, fmt.Sprintf("    Original comment by %s", author))
		authorID = importer.giteaAccessor.GetDefaultAuthorID()
	}
	if len(header) > 0 {
		comment = fmt.Sprintf("%s\n\n%s", strings.Join(header, "\n"), comment)
	}

	return importer.giteaAccessor.AddComment(issueID, authorID, comment, time)
}

func (importer *Importer) importTicketComments(ticketID int64, issueID int64, lastUpdate int64) {
	importer.tracAccessor.GetComments(ticketID, func(ticketID int64, time int64, author string, comment string) {
		if lastUpdate > time {
			lastUpdate = time
		}
		importer.importTicketComment(issueID, ticketID, time, author, comment)
	})

	// Update issue modification time
	importer.giteaAccessor.SetIssueUpdateTime(issueID, lastUpdate)
}
