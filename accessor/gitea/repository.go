package gitea

import (
	"database/sql"
	"fmt"

	"stevejefferson.co.uk/trac2gitea/log"
)

func (accessor *Accessor) getRepoID(userName string, repoName string) int64 {
	row := accessor.db.QueryRow(`
		SELECT r.id FROM repository r, user u WHERE r.owner_id =
			u.id AND u.name = $1 AND r.name = $2
		`, userName, repoName)

	var id int64 = -1
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return id
}

// UpdateRepoIssueCount updates the count of total and closed issue for a our chosen Gitea repository.
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

	log.Infof("Updated repository: total issues=%d, closed issues=%d\n", count, closedCount)
}

// GetCommitURL retrieves the URL for viewing a given commit in the current repository
func (accessor *Accessor) GetCommitURL(commitID string) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/commit/%s", repoURL, commitID)
}

// GetSourceURL retrieves the URL for viewing the latest version of a source file on a given branch of the current repository
func (accessor *Accessor) GetSourceURL(branchPath string, filePath string) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/src/branch/%s/%s", repoURL, branchPath, filePath)
}
