// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func (accessor *DefaultAccessor) getRepoID(userName string, repoName string) (int64, error) {
	var id int64 = -1
	err := accessor.db.QueryRow(`SELECT r.id FROM repository r, user u WHERE r.owner_id =
			u.id AND u.name = $1 AND r.name = $2`, userName, repoName).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of repository %s for user %s", repoName, userName)
		return -1, err
	}

	return id, nil
}

// UpdateRepoIssueCounts updates issue counts for a our chosen Gitea repository.
func (accessor *DefaultAccessor) UpdateRepoIssueCounts() error {
	_, err := accessor.db.Exec(`
		UPDATE repository SET 
			num_issues = (SELECT COUNT(id) FROM issue),
			num_closed_issues = (SELECT COUNT(id) FROM issue WHERE is_closed=1)
			WHERE id = $2`, accessor.repoID)
	if err != nil {
		err = errors.Wrapf(err, "updating number of issues for repository %d", accessor.repoID)
		return err
	}

	return nil
}

// UpdateRepoMilestoneCounts updates milestone counts for a our chosen Gitea repository.
func (accessor *DefaultAccessor) UpdateRepoMilestoneCounts() error {
	_, err := accessor.db.Exec(`
		UPDATE repository SET 
			num_milestones = (SELECT COUNT(id) FROM milestone),
			num_closed_milestones = (SELECT COUNT(id) FROM milestone WHERE is_closed=1)
			WHERE id = $2`, accessor.repoID)
	if err != nil {
		err = errors.Wrapf(err, "updating number of milestones for repository %d", accessor.repoID)
		return err
	}

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
