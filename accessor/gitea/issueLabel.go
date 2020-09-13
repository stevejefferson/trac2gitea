// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// getIssueLabelID retrieves the id of the given Gitea issue label - returns -1 if no such issue label.
func (accessor *DefaultAccessor) getIssueLabelID(issueID int64, labelID int64) (int64, error) {
	var issueLabelID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_label WHERE issue_id=$1 AND label_id=$2
		`, issueID, labelID).Scan(&issueLabelID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of issue label for issue %d, label %d", issueID, labelID)
		return -1, err
	}

	return issueLabelID, nil
}

// insertIssueLabel adds a new label to a Gitea issue, returns id of created issue label.
func (accessor *DefaultAccessor) insertIssueLabel(issueID int64, labelID int64) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_label(issue_id, label_id) VALUES ( $1, $2 )`,
		issueID, labelID)
	if err != nil {
		err = errors.Wrapf(err, "adding issue label for issue %d, label %d", issueID, labelID)
		return -1, err
	}

	var issueLabelID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueLabelID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new issue label for issue %d, label %d", issueID, labelID)
		return -1, err
	}

	log.Debug("added label %d for issue %d (id %d)", labelID, issueID, issueLabelID)

	return issueLabelID, nil
}

// AddIssueLabel adds an issue label to Gitea, returns issue label ID
func (accessor *DefaultAccessor) AddIssueLabel(issueID int64, labelID int64) (int64, error) {
	issueLabelID, err := accessor.getIssueLabelID(issueID, labelID)
	if err != nil {
		return -1, err
	}

	if issueLabelID == -1 {
		return accessor.insertIssueLabel(issueID, labelID)
	}

	// association between issue_id and label_id already exists - nothing to do
	return issueLabelID, nil
}
