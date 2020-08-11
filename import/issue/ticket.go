package issue

import (
	"database/sql"
	"fmt"
	"strings"

	"stevejefferson.co.uk/trac2gitea/markdown"
)

func (importer *Importer) importTicket(
	ticketID int64,
	created int64,
	owner string,
	reporter string,
	milestone string,
	closed bool,
	summary string,
	description string) int64 {
	t2mConverter := markdown.CreateTicketConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	description = t2mConverter.Convert(description)

	var header []string

	// find users first, and tweak description to add missing users
	reporterID := importer.giteaAccessor.GetUserID(reporter)
	if reporterID == -1 {
		header = append(header, fmt.Sprintf("    Originally reported by %s", reporter))
		reporterID = importer.giteaAccessor.DefaultAuthorID
	}
	var ownerID sql.NullString
	if owner != "" {
		tmp := importer.giteaAccessor.GetUserID(owner)
		if tmp == -1 {
			header = append(header, fmt.Sprintf("    Originally assigned to %s", owner))
			ownerID.String = fmt.Sprintf("%d", importer.giteaAccessor.DefaultAssigneeID)
			ownerID.Valid = true
		} else {
			ownerID.String = fmt.Sprintf("%d", tmp)
			ownerID.Valid = true
		}
	} else {
		ownerID.Valid = false
	}
	if len(header) > 0 {
		description = fmt.Sprintf("%s\n\n%s", strings.Join(header, "\n"), description)
	}

	issueID := importer.giteaAccessor.AddIssue(ticketID, summary, reporterID, milestone, ownerID, owner, closed, description, created)

	return issueID
}

// ImportTickets imports Trac tickets as Gitea issues.
func (importer *Importer) ImportTickets() {
	count := 0
	closedCount := 0

	importer.tracAccessor.GetTickets(func(
		ticketID int64, ticketType string, created int64,
		component string, severity string, priority string,
		owner string, reporter string, version string,
		milestone string, status string, resolution string,
		summary string, description string) {
		count++
		closed := status == "closed"
		if closed {
			closedCount++
		}

		issueID := importer.importTicket(ticketID, created, owner, reporter, milestone, closed, summary, description)
		importer.importTicketLabels(issueID, component, severity, priority, version, resolution, ticketType)
		lastUpdate := importer.importTicketAttachments(ticketID, issueID, created)
		importer.importTicketComments(ticketID, issueID, lastUpdate)
	})

	importer.giteaAccessor.UpdateRepoIssueCount(count, closedCount)

	// TODO: Update issue count for new labels
}
