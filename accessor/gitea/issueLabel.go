package gitea

import (
	"database/sql"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetIssueLabelID retrieves the id of the given Gitea issue and label - returns -1 if no such issue label.
func (accessor *DefaultAccessor) GetIssueLabelID(issueID int64, labelID int64) int64 {
	var issueLabelID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_label WHERE issue_id = $1 AND label_id = $2
		`, issueID, labelID).Scan(&issueLabelID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return issueLabelID
}

// AddIssueLabel adds an issue label to Gitea, returns issue label ID
func (accessor *DefaultAccessor) AddIssueLabel(issueID int64, labelID int64) int64 {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_label(issue_id, label_id) VALUES ( $1, $2 )`,
		issueID, labelID)
	if err != nil {
		log.Fatal(err)
	}

	var issueLabelID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueLabelID)
	if err != nil {
		log.Fatal(err)
	}

	return issueLabelID
}
