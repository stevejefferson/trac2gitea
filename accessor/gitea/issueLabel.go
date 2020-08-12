package gitea

import (
	"stevejefferson.co.uk/trac2gitea/log"
)

// AddIssueLabel adds an issue label to Gitea.
func (accessor *Accessor) AddIssueLabel(issueID int64, label string) {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_label(issue_id, label_id)
			SELECT $1, (SELECT id FROM label where repo_id = $2 and name = $3)`,
		issueID, accessor.repoID, label)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Issue %d: created label %s\n", issueID, label)
}
