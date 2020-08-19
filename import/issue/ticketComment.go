package issue

import (
	"fmt"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

// importTicketComment imports a single ticket comment from Trac to Gitea, returns ID of created comment or -1 if comment already exists
func (importer *Importer) importTicketComment(issueID int64, ticketID int64, time int64, author, comment string) int64 {
	markdownConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	comment = markdownConverter.Convert(comment)

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

	truncatedComment := comment[0:15] + "..."
	if importer.giteaAccessor.GetCommentID(issueID, comment) != -1 {
		log.Debugf("Comment \"%s\" for issue %d already exists - skipping...\n", truncatedComment, issueID)
		return -1
	}

	commentID := importer.giteaAccessor.AddComment(issueID, authorID, comment, time)
	log.Debugf("Issue %d: added comment \"%s\" (id %d)\n", issueID, truncatedComment, commentID)

	return commentID
}

func (importer *Importer) importTicketComments(ticketID int64, issueID int64, lastUpdate int64) {
	importer.tracAccessor.GetComments(ticketID, func(ticketID int64, time int64, author string, comment string) {
		commentID := importer.importTicketComment(issueID, ticketID, time, author, comment)

		if commentID != -1 && lastUpdate > time {
			lastUpdate = time
		}
	})

	// Update issue modification time
	importer.giteaAccessor.SetIssueUpdateTime(issueID, lastUpdate)
}
