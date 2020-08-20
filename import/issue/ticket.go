// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package issue

import (
	"database/sql"
	"fmt"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

// importTicket imports a Trac ticket as a Gitea issue, returning the id of the created issue or -1 if the issue was not created.
func (importer *Importer) importTicket(
	ticketID int64,
	created int64,
	owner string,
	reporter string,
	milestone string,
	closed bool,
	summary string,
	description string) (int64, error) {
	issueID, err := importer.giteaAccessor.GetIssueID(ticketID)
	if err != nil {
		return -1, err
	}
	if issueID != -1 {
		log.Infof("Issue already exists for ticket %d - skipping...\n", ticketID)
		return -1, nil
	}

	markdownConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	description = markdownConverter.Convert(description)

	var header []string

	// find users first, and tweak description to add missing users
	reporterID, err := importer.giteaAccessor.GetUserID(reporter)
	if err != nil {
		return -1, err
	}

	if reporterID == -1 {
		header = append(header, fmt.Sprintf("    Originally reported by %s", reporter))
		reporterID = importer.giteaAccessor.GetDefaultAuthorID()
	}
	var ownerID sql.NullString
	if owner != "" {
		tmp, err := importer.giteaAccessor.GetUserID(owner)
		if err != nil {
			return -1, err
		}

		if tmp == -1 {
			header = append(header, fmt.Sprintf("    Originally assigned to %s", owner))
			ownerID.String = fmt.Sprintf("%d", importer.giteaAccessor.GetDefaultAssigneeID())
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

	issueID, err = importer.giteaAccessor.AddIssue(ticketID, summary, reporterID, milestone, ownerID, owner, closed, description, created)
	if err != nil {
		return -1, err
	}
	log.Infof("Created issue %d: %s\n", issueID, summary)

	return issueID, nil
}

// ImportTickets imports Trac tickets as Gitea issues.
func (importer *Importer) ImportTickets() error {
	count := 0
	closedCount := 0

	err := importer.tracAccessor.GetTickets(func(
		ticketID int64, ticketType string, created int64,
		component string, severity string, priority string,
		owner string, reporter string, version string,
		milestone string, status string, resolution string,
		summary string, description string) error {
		closed := status == "closed"
		issueID, err := importer.importTicket(ticketID, created, owner, reporter, milestone, closed, summary, description)
		if err != nil {
			return err
		}
		if issueID == -1 {
			return nil
		}

		err = importer.importTicketLabels(issueID, component, severity, priority, version, resolution, ticketType)
		if err != nil {
			return err
		}

		lastUpdate, err := importer.importTicketAttachments(ticketID, issueID, created)
		if err != nil {
			return err
		}
		err = importer.importTicketComments(ticketID, issueID, lastUpdate)
		if err != nil {
			return err
		}

		count++
		if closed {
			closedCount++
		}

		return nil
	})
	if err != nil {
		return err
	}

	return importer.giteaAccessor.UpdateRepoIssueCount(count, closedCount)
}
