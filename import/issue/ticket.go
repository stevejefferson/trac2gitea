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
	description string,
	userMap map[string]string) (int64, error) {
	issueID, err := importer.giteaAccessor.GetIssueID(ticketID)
	if err != nil {
		return -1, err
	}
	if issueID != -1 {
		// assume we have previously done this conversion
		log.Info("issue already exists for ticket %d - skipping...", ticketID)
		return -1, nil
	}

	reporterID, _, err := importer.getUser(reporter, userMap)
	if err != nil {
		return -1, err
	}
	tracDetails := fmt.Sprintf("originally reported by %s", reporter)

	var ownerID int64 = -1
	var ownerName = ""
	if owner != "" {
		ownerID, ownerName, err = importer.getUser(owner, userMap)
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
	log.Info("created issue %d: %s", issueID, summary)

	return issueID, nil
}

// ImportTickets imports Trac tickets as Gitea issues.
func (importer *Importer) ImportTickets(
	userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) error {
	count := 0
	closedCount := 0

	err := importer.tracAccessor.GetTickets(func(
		ticketID int64, summary string, description string, owner string, reporter string, milestone string,
		component string, priority string, resolution string, severity string, typ string, version string,
		status string, created int64) error {
		closed := status == "closed"
		issueID, err := importer.importTicket(ticketID, created, owner, reporter, milestone, closed, summary, description, userMap)
		if err != nil {
			return err
		}
		if issueID == -1 {
			return nil
		}

		err = importer.importTicketLabels(issueID,
			component, componentMap, priority, priorityMap, resolution, resolutionMap, severity, severityMap, typ, typeMap, version, versionMap)
		if err != nil {
			return err
		}

		lastUpdate, err := importer.importTicketAttachments(ticketID, issueID, created, userMap)
		if err != nil {
			return err
		}
		err = importer.importTicketComments(ticketID, issueID, lastUpdate, userMap)
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
