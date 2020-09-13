// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// getIssueUserID retrieves the id of the given issue/user association, returns -1 if no such association
func (accessor *DefaultAccessor) getIssueUserID(issueID int64, userID int64) (int64, error) {
	var issueUserID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_user WHERE issue_id = $1 AND uid = $2
		`, issueID, userID).Scan(&issueUserID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id for issue %d/user %d", issueID, userID)
		return -1, err
	}

	return issueUserID, nil
}

// updateIssueUser updates an existing issue user
func (accessor *DefaultAccessor) updateIssueUser(issueUserID int64, issueID int64, userID int64) error {
	_, err := accessor.db.Exec(`UPDATE issue_user SET issue_id=?, uid=? WHERE id=?`,
		issueID, userID, issueUserID)
	if err != nil {
		err = errors.Wrapf(err, "updating issue %d/user %d", issueID, userID)
		return err
	}

	log.Debug("updated user %d for issue %d (id %d)", userID, issueID, issueUserID)

	return nil
}

// insertIssueUser associates a new user with a Gitea issue
func (accessor *DefaultAccessor) insertIssueUser(issueID int64, userID int64) error {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_user(issue_id, uid, is_read, is_mentioned) VALUES ($1, $2, 1, 0)`,
		issueID, userID)
	if err != nil {
		err = errors.Wrapf(err, "adding user %d to issue id %d", userID, issueID)
		return err
	}

	log.Debug("added user %d for issue %d", userID, issueID)

	return nil
}

// AddIssueUser adds a user as being associated with a Gitea issue
func (accessor *DefaultAccessor) AddIssueUser(issueID int64, userID int64) error {
	issueUserID, err := accessor.getIssueUserID(issueID, userID)
	if err != nil {
		return err
	}

	if issueUserID == -1 {
		return accessor.insertIssueUser(issueID, userID)
	}

	if accessor.overwrite {
		err = accessor.updateIssueUser(issueUserID, issueID, userID)
		if err != nil {
			return err
		}
	} else {
		log.Debug("issue %d already associated with user %d - ignored", issueID, userID)
	}

	return nil
}
