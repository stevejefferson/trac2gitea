// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"

	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
	"github.com/stevejefferson/trac2gitea/markdown"
)

// importTicket imports a Trac ticket as a Gitea issue, returning the id of the created issue or -1 if the issue was not created.
func (importer *Importer) importTicket(ticket *trac.Ticket, closed bool, userMap map[string]string) (int64, error) {
	issueID, err := importer.giteaAccessor.GetIssueID(ticket.TicketID)
	if err != nil {
		return -1, err
	}
	if issueID != -1 {
		// assume we have previously done this conversion
		log.Info("issue already exists for ticket %d - skipping...", ticket.TicketID)
		return -1, nil
	}

	reporterID, _, err := importer.getUser(ticket.Reporter, userMap)
	if err != nil {
		return -1, err
	}
	tracDetails := fmt.Sprintf("originally reported by %s", ticket.Reporter)

	var ownerID int64 = -1
	var ownerName = ""
	if ticket.Owner != "" {
		ownerID, ownerName, err = importer.getUser(ticket.Owner, userMap)
		if err != nil {
			return -1, err
		}
		tracDetails = tracDetails + fmt.Sprintf(", originally assigned to %s", ticket.Owner)
	}

	// Gitea comment consists of a header giving the original Trac context then the Trac description converted to markdown
	markdownConverter := markdown.CreateTicketDefaultConverter(importer.tracAccessor, importer.giteaAccessor, ticket.TicketID)
	convertedDescription := markdownConverter.Convert(ticket.Description)
	fullDescription := addTracContext(tracDetails, ticket.Created, convertedDescription)

	issueID, err = importer.giteaAccessor.AddIssue(ticket.TicketID, ticket.Summary, reporterID,
		ticket.MilestoneName, ownerID, ownerName, closed, fullDescription, ticket.Created)
	if err != nil {
		return -1, err
	}
	log.Info("created issue %d: %s", issueID, ticket.Summary)

	return issueID, nil
}

// ImportTickets imports Trac tickets as Gitea issues.
func (importer *Importer) ImportTickets(
	userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) error {
	count := 0
	closedCount := 0

	err := importer.tracAccessor.GetTickets(func(ticket *trac.Ticket) error {
		closed := ticket.Status == "closed"
		issueID, err := importer.importTicket(ticket, closed, userMap)
		if err != nil {
			return err
		}
		if issueID == -1 {
			return nil
		}

		_, err = importer.importTicketLabel(issueID, ticket.ComponentName, componentMap, componentLabelColor)
		if err != nil {
			return err
		}

		_, err = importer.importTicketLabel(issueID, ticket.PriorityName, priorityMap, priorityLabelColor)
		if err != nil {
			return err
		}

		_, err = importer.importTicketLabel(issueID, ticket.ResolutionName, resolutionMap, resolutionLabelColor)
		if err != nil {
			return err
		}

		_, err = importer.importTicketLabel(issueID, ticket.SeverityName, severityMap, severityLabelColor)
		if err != nil {
			return err
		}

		_, err = importer.importTicketLabel(issueID, ticket.TypeName, typeMap, typeLabelColor)
		if err != nil {
			return err
		}

		_, err = importer.importTicketLabel(issueID, ticket.VersionName, versionMap, versionLabelColor)
		if err != nil {
			return err
		}

		lastUpdate, err := importer.importTicketAttachments(ticket.TicketID, issueID, ticket.Created, userMap)
		if err != nil {
			return err
		}
		err = importer.importTicketComments(ticket.TicketID, issueID, lastUpdate, userMap)
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
