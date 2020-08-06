package gitea

import (
	"log"
)

func (accessor *Accessor) findRepoID(userName string, repoName string) int64 {
	row := accessor.db.QueryRow(`
		SELECT r.id FROM repository r, user u WHERE r.owner_id =
			u.id AND u.name = $1 AND r.name = $2
		`, userName, repoName)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Fatal("No Gitea repository " + repoName + " found for user " + userName)
	}

	return id
}

func (accessor *Accessor) UpdateRepoIssueCount(count int, closedCount int) {
	// Update issue count for repo
	if count > 0 {
		_, err := accessor.db.Exec(`
			UPDATE repository SET num_issues = num_issues+$1
				WHERE id = $2`,
			count, accessor.repoID)
		if err != nil {
			log.Fatal(err)
		}
	}
	if closedCount > 0 {
		_, err := accessor.db.Exec(`
			UPDATE repository
				SET num_closed_issues = num_closed_issues+$1
				WHERE id = $2`,
			closedCount, accessor.repoID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
