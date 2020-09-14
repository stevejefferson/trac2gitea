// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// GetIssueID retrieves the id of the Gitea issue corresponding to a given issue index - returns -1 if no such issue.
func (accessor *DefaultAccessor) GetIssueID(issueIndex int64) (int64, error) {
	var issueID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue WHERE repo_id = $1 AND "index" = $2
		`, accessor.repoID, issueIndex).Scan(&issueID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving issue with index %d", issueIndex)
		return -1, err
	}

	return issueID, nil
}

func toNullInt64(value int64) sql.NullInt64 {
	var nullValue sql.NullInt64
	nullValue.Valid = (value != -1)
	nullValue.Int64 = value
	return nullValue
}

// updateIssue updates an existing issue in Gitea
func (accessor *DefaultAccessor) updateIssue(issueID int64, issue *Issue) error {
	nullOwnerID := toNullInt64(issue.OriginalAuthorID)
	milestoneID, err := accessor.GetMilestoneID(issue.Milestone)
	if err != nil {
		return err
	}

	_, err = accessor.db.Exec(`
		UPDATE issue SET "index"=?, repo_id=?, name=?, poster_id=?,
			milestone_id=?, original_author_id=?, original_author=?, 
			is_pull=0, is_closed=?, content=?, created_unix=?, updated_unix=?
			WHERE id=?`,
		issue.Index, accessor.repoID, issue.Summary, issue.ReporterID,
		milestoneID, nullOwnerID, issue.OriginalAuthorName,
		issue.Closed, issue.Description, issue.Created, issue.Updated,
		issueID)
	if err != nil {
		err = errors.Wrapf(err, "updating issue with index %d", issue.Index)
		return err
	}

	log.Info("updated issue %d: %s", issue.Index, issue.Summary)

	return nil
}

// insertIssue adds a new issue to Gitea, returns id of added issue.
func (accessor *DefaultAccessor) insertIssue(issue *Issue) (int64, error) {
	nullOwnerID := toNullInt64(issue.OriginalAuthorID)
	milestoneID, err := accessor.GetMilestoneID(issue.Milestone)
	if err != nil {
		return -1, err
	}

	_, err = accessor.db.Exec(`
		INSERT INTO issue("index", repo_id, name, poster_id, milestone_id, original_author_id, original_author, is_pull, is_closed, content, created_unix)
			SELECT $1, $2, $3, $4, $5, $6, $7, 0, $8, $9, $10`,
		issue.Index, accessor.repoID, issue.Summary, issue.ReporterID, milestoneID, nullOwnerID, issue.OriginalAuthorName, issue.Closed, issue.Description, issue.Created)
	if err != nil {
		err = errors.Wrapf(err, "adding issue with index %d", issue.Index)
		return -1, err
	}

	var issueID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new issue with index %d", issue.Index)
		return -1, err
	}

	log.Info("created issue %d: %s", issue.Index, issue.Summary)

	return issueID, nil
}

// AddIssue adds a new issue to Gitea.
func (accessor *DefaultAccessor) AddIssue(issue *Issue) (int64, error) {
	issueID, err := accessor.GetIssueID(issue.Index)
	if err != nil {
		return -1, err
	}

	if issueID == -1 {
		return accessor.insertIssue(issue)
	}

	if accessor.overwrite {
		err = accessor.updateIssue(issueID, issue)
		if err != nil {
			return -1, err
		}
	} else {
		log.Info("issue %d already exists - ignored", issue.Index)
	}

	return issueID, nil
}

// SetIssueUpdateTime sets the update time on a given Gitea issue.
func (accessor *DefaultAccessor) SetIssueUpdateTime(issueID int64, updateTime int64) error {
	_, err := accessor.db.Exec(`UPDATE issue SET updated_unix = MAX(updated_unix,$1) WHERE id = $2`, updateTime, issueID)
	if err != nil {
		err = errors.Wrapf(err, "setting updated time for issue %d", issueID)
		return err
	}

	return nil
}

// GetIssueURL retrieves a URL for viewing a given issue
func (accessor *DefaultAccessor) GetIssueURL(issueID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/issues/%d", repoURL, issueID)
}

// UpdateIssueCommentCount updates the count of comments a given issue
func (accessor *DefaultAccessor) UpdateIssueCommentCount(issueID int64) error {
	_, err := accessor.db.Exec(`
	UPDATE issue SET 
		num_comments = (SELECT COUNT(id) FROM comment)
		WHERE id = $1`, issueID)
	if err != nil {
		err = errors.Wrapf(err, "updating number of comments for issue %d", issueID)
		return err
	}

	return nil
}
