// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"

	"stevejefferson.co.uk/trac2gitea/log"
)

// AddComment adds a comment to Gitea, returns id of created comment
func (accessor *DefaultAccessor) AddComment(issueID int64, authorID int64, comment string, time int64) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO comment(
			type, issue_id, poster_id, content, created_unix, updated_unix)
			VALUES ( 0, $1, $2, $3, $4, $4 )`,
		issueID, authorID, comment, time)
	if err != nil {
		log.Error("Cannot insert comment for issue %d: %v\n", issueID, err)
		return -1, err
	}

	var commentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&commentID)
	if err != nil {
		log.Error("Cannot find id of newly-inserted comment for issue %d: %v\n", issueID, err)
		return -1, err
	}

	return commentID, nil
}

// GetCommentID retrives the ID of a given comment for a given issue or -1 if no such issue/comment
func (accessor *DefaultAccessor) GetCommentID(issueID int64, commentStr string) (int64, error) {
	var commentID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM comment WHERE issue_id = $1 AND content = $2
		`, issueID, commentStr).Scan(&commentID)
	if err != nil && err != sql.ErrNoRows {
		log.Error("Cannot find comment for issue %d: %v\n", issueID, err)
		return -1, err
	}

	return commentID, nil
}

// GetCommentURL retrieves the URL for viewing a Gitea comment for a given issue.
func (accessor *DefaultAccessor) GetCommentURL(issueID int64, commentID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/issues/%d#issuecomment-%d", repoURL, issueID, commentID)
}
