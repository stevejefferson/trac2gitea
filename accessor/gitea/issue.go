package gitea

import (
	"database/sql"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetIssueID retrieves the id of the Gitea issue corresponding to a given Trac ticket - returns -1 if no such issue.
func (accessor *DefaultAccessor) GetIssueID(ticketID int64) int64 {
	var issueID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue WHERE repo_id = $1 AND "index" = $2
		`, accessor.repoID, ticketID).Scan(&issueID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return issueID
}

// AddIssue adds a new issue to Gitea.
func (accessor *DefaultAccessor) AddIssue(
	ticketID int64,
	summary string,
	reporterID int64,
	milestone string,
	ownerID sql.NullString,
	owner string,
	closed bool,
	description string,
	created int64) int64 {
	_, err := accessor.db.Exec(`
		INSERT INTO issue("index", repo_id, name, poster_id, milestone_id, original_author_id, original_author, is_pull, is_closed, content, created_unix)
			SELECT $1, $2, $3, $4, (SELECT id FROM milestone WHERE repo_id = $2 AND name = $5), $6, $7, false, $8, $9, $10`,
		ticketID, accessor.repoID, summary, reporterID, milestone, ownerID, owner, closed, description, created)
	if err != nil {
		log.Fatal(err)
	}

	var issueID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueID)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Created issue %d: %s\n", issueID, summary)

	return issueID
}

// SetIssueUpdateTime sets the update time on a given Gitea issue.
func (accessor *DefaultAccessor) SetIssueUpdateTime(issueID int64, updateTime int64) {
	_, err := accessor.db.Exec(`UPDATE issue SET updated_unix = MAX(updated_unix,$1) WHERE id = $2`, updateTime, issueID)
	if err != nil {
		log.Fatal(err)
	}
}
