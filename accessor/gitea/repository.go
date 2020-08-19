package gitea

import (
	"database/sql"
	"fmt"

	"stevejefferson.co.uk/trac2gitea/log"
)

func (accessor *DefaultAccessor) getRepoID(userName string, repoName string) (int64, error) {
	var id int64 = -1
	err := accessor.db.QueryRow(`SELECT r.id FROM repository r, user u WHERE r.owner_id =
			u.id AND u.name = $1 AND r.name = $2`, userName, repoName).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return -1, err
	}

	return id, nil
}

// UpdateRepoIssueCount updates the count of total and closed issue for a our chosen Gitea repository.
func (accessor *DefaultAccessor) UpdateRepoIssueCount(count int, closedCount int) error {
	// Update issue count for repo
	if count > 0 {
		_, err := accessor.db.Exec(`
			UPDATE repository SET num_issues = num_issues+$1
				WHERE id = $2`,
			count, accessor.repoID)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	if closedCount > 0 {
		_, err := accessor.db.Exec(`
			UPDATE repository
				SET num_closed_issues = num_closed_issues+$1
				WHERE id = $2`,
			closedCount, accessor.repoID)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	log.Infof("Updated repository: total issues=%d, closed issues=%d\n", count, closedCount)
	return nil
}

// GetCommitURL retrieves the URL for viewing a given commit in the current repository
func (accessor *DefaultAccessor) GetCommitURL(commitID string) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/commit/%s", repoURL, commitID)
}

// GetSourceURL retrieves the URL for viewing the latest version of a source file on a given branch of the current repository
func (accessor *DefaultAccessor) GetSourceURL(branchPath string, filePath string) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/src/branch/%s/%s", repoURL, branchPath, filePath)
}
