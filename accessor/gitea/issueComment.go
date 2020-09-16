// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// GetIssueCommentIDsByTime retrieves the IDs of all comments created at a given time for a given issue
func (accessor *DefaultAccessor) GetIssueCommentIDsByTime(issueID int64, createdTime int64) ([]int64, error) {
	rows, err := accessor.db.Query(
		`SELECT id FROM comment WHERE issue_id = $1 AND created_unix = $2`, issueID, createdTime)
	if err != nil {
		err = errors.Wrapf(err, "retrieving ids of comments created at \"%s\" for issue %d", time.Unix(createdTime, 0), issueID)
		return []int64{}, err
	}

	var issueCommentIDs = []int64{}
	for rows.Next() {
		var issueCommentID int64 = -1
		if err := rows.Scan(&issueCommentID); err != nil {
			err = errors.Wrapf(err, "retrieving id of comment created at \"%s\" for issue %d", time.Unix(createdTime, 0), issueID)
			return []int64{}, err
		}

		issueCommentIDs = append(issueCommentIDs, issueCommentID)
	}

	return issueCommentIDs, nil
}

// updateIssueComment updates an existing issue comment
func (accessor *DefaultAccessor) updateIssueComment(issueCommentID int64, issueID int64, comment *IssueComment) error {
	_, err := accessor.db.Exec(`
		UPDATE comment SET
			type=?, issue_id=?, poster_id=?,
			original_author_id=?, original_author=?, 
			old_milestone_id=?, milestone_id=?,
			assignee_id=?, removed_assignee=?,
			old_title=?, new_title=?,
			content=?,
			created_unix=?, updated_unix=?
			WHERE id=?`,
		comment.CommentType, issueID, comment.AuthorID,
		comment.OriginalAuthorID, comment.OriginalAuthorName,
		comment.OldMilestoneID, comment.MilestoneID,
		comment.AssigneeID, comment.RemovedAssigneeID,
		comment.OldTitle, comment.Title,
		comment.Text,
		comment.Time, comment.Time,
		issueCommentID)
	if err != nil {
		err = errors.Wrapf(err, "updating comment on issue %d timed at %s", issueID, time.Unix(comment.Time, 0))
		return err
	}

	log.Debug("updated issue comment at %s for issue %d (id %d)", time.Unix(comment.Time, 0), issueID, issueCommentID)

	return nil
}

// insertIssueComment adds a new comment to a Gitea issue, returns id of created comment.
func (accessor *DefaultAccessor) insertIssueComment(issueID int64, comment *IssueComment) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO comment(
			type, issue_id, poster_id, 
			original_author_id, original_author, 
			old_milestone_id, milestone_id,
			assignee_id, removed_assignee,
			old_title, new_title,
			content, 
			created_unix, updated_unix)
			VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 )`,
		comment.CommentType, issueID, comment.AuthorID,
		comment.OriginalAuthorID, comment.OriginalAuthorName,
		comment.OldMilestoneID, comment.MilestoneID,
		comment.AssigneeID, comment.RemovedAssigneeID,
		comment.OldTitle, comment.Title,
		comment.Text,
		comment.Time, comment.Time)
	if err != nil {
		err = errors.Wrapf(err, "adding comment \"%s\" for issue %d", comment.Text, issueID)
		return -1, err
	}

	var issueCommentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueCommentID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new comment \"%s\" for issue %d", comment.Text, issueID)
		return -1, err
	}

	log.Debug("added issue comment at %s for issue %d (id %d)", time.Unix(comment.Time, 0), issueID, issueCommentID)

	return issueCommentID, nil
}

var prevIssueID = int64(0)
var prevCommentTime = int64(0)
var issueCommentIDIndex = 0
var issueCommentIDs = []int64{}

// AddIssueComment adds a comment on a Gitea issue, returns id of created comment
func (accessor *DefaultAccessor) AddIssueComment(issueID int64, comment *IssueComment) (int64, error) {
	var err error

	// HACK:
	// Timestamps associated with Gitea comments are not necessarily unique for comments originating from Trac
	// because Trac stores timestamps to a greater precision which we have to round to Gitea's precision.
	// Unfortunately timestamp is the best key we have for identifying whether a particular issue comment already exists
	// (and hence whether we need to insert or update it).
	// We get round this by observing that comments are always added consecutively for a given issue so we can
	// cache all comment IDs for our current issue and timestamp and extract the subsequent entries from that list.
	if issueID != prevIssueID || comment.Time != prevCommentTime {
		prevIssueID = issueID
		prevCommentTime = comment.Time
		issueCommentIDIndex = 0
		issueCommentIDs, err = accessor.GetIssueCommentIDsByTime(issueID, comment.Time)
		if err != nil {
			return -1, err
		}
	}

	if issueCommentIDIndex >= len(issueCommentIDs) {
		// should only happen where no issue comments for timestamp
		return accessor.insertIssueComment(issueID, comment)
	}

	issueCommentID := issueCommentIDs[issueCommentIDIndex]
	issueCommentIDIndex++

	if accessor.overwrite {
		err := accessor.updateIssueComment(issueCommentID, issueID, comment)
		if err != nil {
			return -1, err
		}
	} else {
		log.Debug("issue %d already has comment timed at %s - ignored", issueID, time.Unix(comment.Time, 0))
	}

	return issueCommentID, nil
}

// GetIssueCommentURL retrieves the URL for viewing a Gitea comment for a given issue.
func (accessor *DefaultAccessor) GetIssueCommentURL(issueID int64, commentID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/issues/%d#issuecomment-%d", repoURL, issueID, commentID)
}
