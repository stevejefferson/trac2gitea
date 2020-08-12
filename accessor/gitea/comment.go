package gitea

import (
	"fmt"

	"stevejefferson.co.uk/trac2gitea/log"
)

// AddComment adds a comment to Gitea
func (accessor *Accessor) AddComment(issueID int64, authorID int64, comment string, time int64) int64 {
	_, err := accessor.db.Exec(`
		INSERT INTO comment(
			type, issue_id, poster_id, content, created_unix, updated_unix)
			VALUES ( 0, $1, $2, $3, $4, $4 )`,
		issueID, authorID, comment, time)
	if err != nil {
		log.Fatal(err)
	}

	var commentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&commentID)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Issue %d: added comment (id %d)\n", issueID, commentID)

	return commentID
}

// GetCommentURL retrieves the URL for viewing a Gitea comment for a given issue.
func (accessor *Accessor) GetCommentURL(issueID int64, commentID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/issues/%d#issuecomment-%d", repoURL, issueID, commentID)
}
