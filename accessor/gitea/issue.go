// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetIssueID retrieves the id of the Gitea issue corresponding to a given issue index - returns -1 if no such issue.
func (accessor *DefaultAccessor) GetIssueID(issueIndex int64) (int64, error) {
	var issueID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue WHERE repo_id = $1 AND "index" = $2
		`, accessor.repoID, issueIndex).Scan(&issueID)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return -1, err
	}

	return issueID, nil
}

// AddIssue adds a new issue to Gitea.
func (accessor *DefaultAccessor) AddIssue(
	issueIndex int64,
	summary string,
	reporterID int64,
	milestone string,
	ownerID sql.NullString,
	owner string,
	closed bool,
	description string,
	created int64) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO issue("index", repo_id, name, poster_id, milestone_id, original_author_id, original_author, is_pull, is_closed, content, created_unix)
			SELECT $1, $2, $3, $4, (SELECT id FROM milestone WHERE repo_id = $2 AND name = $5), $6, $7, false, $8, $9, $10`,
		issueIndex, accessor.repoID, summary, reporterID, milestone, ownerID, owner, closed, description, created)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	var issueID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueID)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	return issueID, nil
}

// SetIssueUpdateTime sets the update time on a given Gitea issue.
func (accessor *DefaultAccessor) SetIssueUpdateTime(issueID int64, updateTime int64) error {
	_, err := accessor.db.Exec(`UPDATE issue SET updated_unix = MAX(updated_unix,$1) WHERE id = $2`, updateTime, issueID)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// GetIssueURL retrieves a URL for viewing a given issue
func (accessor *DefaultAccessor) GetIssueURL(issueID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/issues/%d", repoURL, issueID)
}
