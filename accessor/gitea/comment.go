// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// AddIssueComment adds a comment on a Gitea issue, returns id of created comment
func (accessor *DefaultAccessor) AddIssueComment(issueID int64, comment *IssueComment) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO comment(
			type, issue_id, poster_id, original_author_id, original_author, content, created_unix, updated_unix)
			VALUES ( 0, $1, $2, $3, $4, $4 )`,
		issueID, comment.AuthorID, comment.OriginalAuthorID, comment.OriginalAuthorName, comment.Text, comment.Time)
	if err != nil {
		err = errors.Wrapf(err, "adding comment \"%s\" for issue %d", comment.Text, issueID)
		return -1, err
	}

	var commentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&commentID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new comment \"%s\" for issue %d", comment.Text, issueID)
		return -1, err
	}

	return commentID, nil
}

// GetTimedIssueCommentID retrives the ID of a comment created at a given time for a given issue or -1 if no such issue/comment
func (accessor *DefaultAccessor) GetTimedIssueCommentID(issueID int64, createdTime int64) (int64, error) {
	var commentID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM comment WHERE issue_id = $1 AND created_unix = $2
		`, issueID, createdTime).Scan(&commentID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of comment created at \"%s\" for issue %d", time.Unix(createdTime, 0), issueID)
		return -1, err
	}

	return commentID, nil
}

// GetIssueCommentURL retrieves the URL for viewing a Gitea comment for a given issue.
func (accessor *DefaultAccessor) GetIssueCommentURL(issueID int64, commentID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/issues/%d#issuecomment-%d", repoURL, issueID, commentID)
}
