// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// getIssueLabelID retrieves the id of the given Gitea issue label - returns NullID if no such issue label.
func (accessor *DefaultAccessor) getIssueLabelID(issueID int64, labelID int64) (int64, error) {
	var issueLabelID int64 = NullID
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_label WHERE issue_id=$1 AND label_id=$2
		`, issueID, labelID).Scan(&issueLabelID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of issue label for issue %d, label %d", issueID, labelID)
		return NullID, err
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
		return NullID, err
	}

	var issueLabelID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueLabelID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new issue label for issue %d, label %d", issueID, labelID)
		return NullID, err
	}

	log.Debug("added label %d for issue %d (id %d)", labelID, issueID, issueLabelID)

	return issueLabelID, nil
}

// AddIssueLabel adds an issue label to Gitea, returns issue label ID
func (accessor *DefaultAccessor) AddIssueLabel(issueID int64, labelID int64) (int64, error) {
	issueLabelID, err := accessor.getIssueLabelID(issueID, labelID)
	if err != nil {
		return NullID, err
	}

	if issueLabelID == NullID {
		return accessor.insertIssueLabel(issueID, labelID)
	}

	// association between issue_id and label_id already exists - nothing to do
	return issueLabelID, nil
}

// UpdateLabelIssueCounts updates issue counts for all labels.
func (accessor *DefaultAccessor) UpdateLabelIssueCounts() error {
	_, err := accessor.db.Exec(`
		UPDATE label AS l SET 
			num_issues = (
				SELECT COUNT(il1.issue_id)
				FROM issue_label il1
				WHERE l.id = il1.label_id
				GROUP BY il1.label_id),
			num_closed_issues = (
				SELECT COUNT(il2.issue_id)
				FROM issue_label il2, issue i
				WHERE l.id = il2.label_id
				AND il2.issue_id = i.id
				AND i.is_closed = 1
				GROUP BY il2.label_id)`)
	if err != nil {
		err = errors.Wrapf(err, "updating number of issues for labels")
		return err
	}

	return nil
}
