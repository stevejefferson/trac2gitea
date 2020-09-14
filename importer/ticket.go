// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// importTicket imports a Trac ticket as a Gitea issue, returning the id of the created issue or -1 if the issue was not created.
func (importer *Importer) importTicket(ticket *trac.Ticket, closed bool, userMap map[string]string) (int64, error) {
	reporterID, err := importer.getUser(ticket.Reporter, userMap)
	if err != nil {
		return -1, err
	}
	if reporterID == -1 {
		reporterID = importer.defaultAuthorID
	}

	// record Trac owner as original author if it cannot be mapped onto a Gitea user
	ownerID := int64(-1)
	originalAuthorName := ticket.Owner
	if ticket.Owner != "" {
		ownerID, err = importer.getUser(ticket.Owner, userMap)
		if err != nil {
			return -1, err
		}
		if ownerID != -1 {
			originalAuthorName = ""
		}
	}

	convertedDescription := importer.markdownConverter.TicketConvert(ticket.TicketID, ticket.Description)
	issue := gitea.Issue{Index: ticket.TicketID, Summary: ticket.Summary, ReporterID: reporterID,
		Milestone: ticket.MilestoneName, OriginalAuthorID: 0, OriginalAuthorName: originalAuthorName,
		Closed: closed, Description: convertedDescription, Created: ticket.Created, Updated: ticket.Updated}
	issueID, err := importer.giteaAccessor.AddIssue(&issue)
	if err != nil {
		return -1, err
	}

	// if we have a Gitea user for the Trac ticket owner then assign the Gitea issue to that user
	if ownerID != -1 {
		err = importer.giteaAccessor.AddIssueAssignee(issueID, ownerID)
		if err != nil {
			return -1, err
		}
	}

	// add associations between our issue and its reporter and its owner (if a different user)
	err = importer.giteaAccessor.AddIssueUser(issueID, reporterID)
	if err != nil {
		return -1, err
	}
	if ownerID != -1 && ownerID != reporterID {
		err = importer.giteaAccessor.AddIssueUser(issueID, ownerID)
		if err != nil {
			return -1, err
		}
	}

	return issueID, nil
}

// ImportTickets imports Trac tickets as Gitea issues.
func (importer *Importer) ImportTickets(
	userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) error {
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
		lastUpdate, err = importer.importTicketChanges(ticket.TicketID, issueID, lastUpdate, userMap)
		if err != nil {
			return err
		}

		err = importer.giteaAccessor.SetIssueUpdateTime(issueID, lastUpdate)
		err = importer.giteaAccessor.UpdateIssueCommentCount(issueID)

		return nil
	})
	if err != nil {
		return err
	}

	return importer.giteaAccessor.UpdateRepoIssueCounts()
}
