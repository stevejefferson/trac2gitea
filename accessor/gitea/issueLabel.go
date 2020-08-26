// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
)

// GetIssueLabelID retrieves the id of the given Gitea issue and label - returns -1 if no such issue label.
func (accessor *DefaultAccessor) GetIssueLabelID(issueID int64, labelID int64) (int64, error) {
	var issueLabelID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_label WHERE issue_id = $1 AND label_id = $2
		`, issueID, labelID).Scan(&issueLabelID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of issue label for issue %d, label %d", issueID, labelID)
		return -1, err
	}

	return issueLabelID, nil
}

// AddIssueLabel adds an issue label to Gitea, returns issue label ID
func (accessor *DefaultAccessor) AddIssueLabel(issueID int64, labelID int64) (int64, error) {
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

	return issueLabelID, nil
}
