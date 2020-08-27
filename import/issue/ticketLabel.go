// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import "github.com/stevejefferson/trac2gitea/log"

// importTicketLabel imports a single issue label from Trac into Gitea, returns id of created issue label or -1 if issue label already exists
func (importer *Importer) importTicketLabel(issueID int64, tracLabel string, labelMap map[string]string, labelColor string) (int64, error) {
	labelID, err := importer.importLabel(tracLabel, labelMap, labelColor)
	if err != nil {
		return -1, err
	}
	if labelID == -1 {
		return -1, nil
	}

	issueLabelID, err := importer.giteaAccessor.GetIssueLabelID(issueID, labelID)
	if err != nil {
		return -1, err
	}
	if issueLabelID != -1 {
		log.Debug("Trac label %s already referenced by issue %d - skipping...", tracLabel, issueID)
		return -1, nil
	}

	issueLabelID, err = importer.giteaAccessor.AddIssueLabel(issueID, labelID)
	if err != nil {
		return -1, err
	}

	log.Debug("created issue label (id %d) for issue %d, label %d", issueLabelID, issueID, labelID)

	return issueLabelID, nil
}

func (importer *Importer) importTicketLabels(
	issueID int64,
	component string, componentMap map[string]string,
	priority string, priorityMap map[string]string,
	resolution string, resolutionMap map[string]string,
	severity string, severityMap map[string]string,
	typ string, typeMap map[string]string,
	version string, versionMap map[string]string) error {
	_, err := importer.importTicketLabel(issueID, component, componentMap, componentLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, priority, priorityMap, priorityLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, resolution, resolutionMap, resolutionLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, severity, severityMap, severityLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, typ, typeMap, typeLabelColor)
	if err != nil {
		return err
	}

	_, err = importer.importTicketLabel(issueID, version, versionMap, versionLabelColor)
	if err != nil {
		return err
	}

	return nil
}
