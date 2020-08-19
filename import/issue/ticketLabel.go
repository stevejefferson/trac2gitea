package issue

import "stevejefferson.co.uk/trac2gitea/log"

// importTicketLabel imports a single issue label from Trac into Gitea, returns id of created issue label or -1 if issue label already exists
func (importer *Importer) importTicketLabel(issueID int64, tracLabel string, labelPrefix string, labelColor string) (int64, error) {
	if tracLabel == "" {
		return -1, nil
	}

	labelName := labelPrefix + tracLabel
	labelID, err := importer.giteaAccessor.GetLabelID(labelName)
	if err != nil {
		return -1, err
	}
	if labelID == -1 {
		log.Warnf("Cannot find label \"%s\" referenced by issue %d - creating it\n", labelName, issueID)
		labelID, err = importer.giteaAccessor.AddLabel(labelName, labelColor)
		if err != nil {
			return -1, err
		}
	}

	issueLabelID, err := importer.giteaAccessor.GetIssueLabelID(issueID, labelID)
	if err != nil {
		return -1, err
	}
	if issueLabelID != -1 {
		log.Debugf("Label %s already referenced by issue %d - skipping...\n", labelName, issueID)
		return -1, nil
	}

	issueLabelID, err = importer.giteaAccessor.AddIssueLabel(issueID, labelID)
	if err != nil {
		return -1, err
	}

	log.Debugf("Created issue label (id %d) for issue %d, label %d\n", issueLabelID, issueID, labelID)

	return issueLabelID, nil
}

func (importer *Importer) importTicketLabels(issueID int64, component string, severity string, priority string, version string, resolution string, typ string) error {
	_, err := importer.importTicketLabel(issueID, component, componentLabelPrefix, componentLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, severity, severityLabelPrefix, severityLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, priority, priorityLabelPrefix, priorityLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, version, versionLabelPrefix, versionLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, resolution, resolutionLabelPrefix, resolutionLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, typ, typeLabelPrefix, typeLabelColor)
	if err != nil {
		return err
	}

	return nil
}
