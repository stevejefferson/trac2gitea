// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
)

// AddIssueUser adds a user as being associated with a Gitea issue
func (accessor *DefaultAccessor) AddIssueUser(issueID int64, userID int64) error {
	// association may already exist so check before inserting...
	var issueUserID int64
	err := accessor.db.QueryRow(`SELECT id FROM issue_user WHERE issue_id = $1 AND uid = $2`, issueID, userID).Scan(&issueUserID)
	if err == nil {
		// issue <--> user association already exists - nothing to do here
		return nil
	} else if err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of association between issue %d and user %d", issueID, userID)
		return err
	}

	_, err = accessor.db.Exec(`
		INSERT INTO issue_user(issue_id, uid, is_read, is_mentioned) VALUES ($1, $2, 1, 0)`,
		issueID, userID)
	if err != nil {
		err = errors.Wrapf(err, "adding user %d as associated with issue id %d", userID, issueID)
		return err
	}

	return nil
}
