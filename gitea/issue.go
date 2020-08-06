package gitea

import (
	"database/sql"
	"log"
)

func (accessor *Accessor) AddIssue(
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
		INSERT INTO issue('index', repo_id, name, poster_id, milestone_id, original_author_id, original_author, is_pull, is_closed, content, created_unix)
			SELECT $1, $2, $3, $4, (SELECT id FROM milestone WHERE repo_id = $2 AND name = $5), $6, $7, false, $8, $9, $10`,
		ticketID, accessor.repoID, summary, reporterID, milestone, ownerID, owner, closed, description, created)
	if err != nil {
		log.Fatal(err)
	}

	var gid int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&gid)
	if err != nil {
		log.Fatal(err)
	}

	return gid
}

func (accessor *Accessor) SetIssueUpdateTime(issueID int64, updateTime int64) {
	_, err := accessor.db.Exec(`UPDATE issue SET updated_unix = MAX(updated_unix,$1) WHERE id = $2`, updateTime, issueID)
	if err != nil {
		log.Fatal(err)
	}
}
