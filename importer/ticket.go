// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
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

	reporterID, err := importer.getUser(ticket.Reporter, userMap)
	if err != nil {
		return -1, err
	}
	if reporterID == -1 {
		reporterID = importer.defaultAuthorID
	}

	convertedDescription := importer.markdownConverter.TicketConvert(ticket.TicketID, ticket.Description)
	issue := gitea.Issue{Index: ticket.TicketID, Summary: ticket.Summary, ReporterID: reporterID,
		Milestone: ticket.MilestoneName, OriginalAuthorID: 0, OriginalAuthorName: ticket.Owner,
		Closed: closed, Description: convertedDescription, Created: ticket.Created}
	issueID, err = importer.giteaAccessor.AddIssue(&issue)
	if err != nil {
		return -1, err
	}
	log.Info("created issue %d: %s", issueID, ticket.Summary)

	// if we have a Gitea user for the Trac ticket owner then assign the Gitea issue to that user
	if ticket.Owner != "" {
		ownerID, err := importer.getUser(ticket.Owner, userMap)
		if err != nil {
			return -1, err
		}
		if ownerID != -1 {
			err = importer.giteaAccessor.AddIssueAssignee(issueID, ownerID)
			if err != nil {
				return -1, err
			}
		}
	}

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
		lastUpdate, err = importer.importTicketComments(ticket.TicketID, issueID, lastUpdate, userMap)
		if err != nil {
			return err
		}

		count++
		if closed {
			closedCount++
		}

		err = importer.giteaAccessor.SetIssueUpdateTime(issueID, lastUpdate)

		return nil
	})
	if err != nil {
		return err
	}

	return importer.giteaAccessor.UpdateRepoIssueCount(count, closedCount)
}
