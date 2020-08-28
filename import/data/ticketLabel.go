// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package data

import "github.com/stevejefferson/trac2gitea/log"

// importTicketLabel imports a single issue label from Trac into Gitea, returns id of created issue label or -1 if issue label already exists
func (importer *Importer) importTicketLabel(issueID int64, tracName string, labelMap map[string]string, labelColor string) (int64, error) {
	labelID, err := importer.importLabel(tracName, labelMap, labelColor)
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
		log.Debug("Trac label %s already referenced by issue %d - skipping...", tracName, issueID)
		return -1, nil
	}

	issueLabelID, err = importer.giteaAccessor.AddIssueLabel(issueID, labelID)
	if err != nil {
		return -1, err
	}

	log.Debug("created issue label (id %d) for issue %d, label %d", issueLabelID, issueID, labelID)

	return issueLabelID, nil
}
