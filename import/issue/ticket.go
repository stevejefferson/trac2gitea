// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"

	"github.com/stevejefferson/trac2gitea/log"
	"github.com/stevejefferson/trac2gitea/markdown"
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
		// assume we have previously done this conversion
		log.Info("Issue already exists for ticket %d - skipping...\n", ticketID)
		return -1, nil
	}

	reporterID, _, err := importer.getUser(reporter)
	if err != nil {
		return -1, err
	}
	tracDetails := fmt.Sprintf("originally reported by %s", reporter)

	var ownerID int64 = -1
	var ownerName = ""
	if owner != "" {
		ownerID, ownerName, err = importer.getUser(owner)
		if err != nil {
			return -1, err
		}
		tracDetails = tracDetails + fmt.Sprintf(", originally assigned to %s", owner)
	}

	// Gitea comment consists of a header giving the original Trac context then the Trac description converted to markdown
	markdownConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticketID)
	convertedDescription := markdownConverter.Convert(description)
	fullDescription := addTracContext(tracDetails, created, convertedDescription)

	issueID, err = importer.giteaAccessor.AddIssue(ticketID, summary, reporterID, milestone, ownerID, ownerName, closed, fullDescription, created)
	if err != nil {
		return -1, err
	}
	log.Info("Created issue %d: %s\n", issueID, summary)

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
