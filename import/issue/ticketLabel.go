package issue

import "stevejefferson.co.uk/trac2gitea/log"

// importTicketLabel imports a single issue label from Trac into Gitea, returns id of created issue label or -1 if issue label already exists
func (importer *Importer) importTicketLabel(issueID int64, tracLabel string, labelPrefix string, labelColor string) int64 {
	if tracLabel == "" {
		return -1
	}

	labelName := labelPrefix + tracLabel
	labelID := importer.giteaAccessor.GetLabelID(labelName)
	if labelID == -1 {
		log.Warnf("Cannot find label \"%s\" referenced by issue %d - creating it\n", labelName, issueID)
		labelID = importer.giteaAccessor.AddLabel(labelName, labelColor)
	}

	if importer.giteaAccessor.GetIssueLabelID(issueID, labelID) != -1 {
		log.Debugf("Label %s already referenced by issue %d - skipping...\n", labelName, issueID)
		return -1
	}

	issueLabelID := importer.giteaAccessor.AddIssueLabel(issueID, labelID)
	log.Debugf("Created issue label (id %d) for issue %d, label %d\n", issueLabelID, issueID, labelID)

	return issueLabelID
}

func (importer *Importer) importTicketLabels(issueID int64, component string, severity string, priority string, version string, resolution string, typ string) {
	importer.importTicketLabel(issueID, component, componentLabelPrefix, componentLabelColor)
	importer.importTicketLabel(issueID, severity, severityLabelPrefix, severityLabelColor)
	importer.importTicketLabel(issueID, priority, priorityLabelPrefix, priorityLabelColor)
	importer.importTicketLabel(issueID, version, versionLabelPrefix, versionLabelColor)
	importer.importTicketLabel(issueID, resolution, resolutionLabelPrefix, resolutionLabelColor)
	importer.importTicketLabel(issueID, typ, typeLabelPrefix, typeLabelColor)
}
