// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// GetMilestoneID gets the ID of a named milestone - returns NullID if no such milestone
func (accessor *DefaultAccessor) GetMilestoneID(milestoneName string) (int64, error) {
	var milestoneID int64 = NullID
	err := accessor.db.QueryRow(`SELECT id FROM milestone WHERE name = $1 AND repo_id = $2`, milestoneName, accessor.repoID).Scan(&milestoneID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of milestone %s", milestoneName)
		return NullID, err
	}

	return milestoneID, nil
}

// updateMilestone updates an existing milestone
func (accessor *DefaultAccessor) updateMilestone(milestoneID int64, milestone *Milestone) error {
	_, err := accessor.db.Exec(`
		UPDATE milestone SET repo_id=?, name=?, content=?, is_closed=?, deadline_unix=?, closed_date_unix=? WHERE id=?`,
		accessor.repoID, milestone.Name, milestone.Description, milestone.Closed, milestone.DueTime, milestone.ClosedTime, milestoneID)
	if err != nil {
		err = errors.Wrapf(err, "updating milestone %s", milestone.Name)
		return err
	}

	log.Debug("updated milestone %s (id %d)", milestone.Name, milestoneID)

	return nil
}

// insertMilestone inserts a new milestone, returns milstone id.
func (accessor *DefaultAccessor) insertMilestone(milestone *Milestone) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO	milestone(repo_id, name, content, is_closed, deadline_unix, closed_date_unix) VALUES($1, $2, $3, $4, $5, $6)`,
		accessor.repoID, milestone.Name, milestone.Description, milestone.Closed, milestone.DueTime, milestone.ClosedTime)
	if err != nil {
		err = errors.Wrapf(err, "adding milestone %s", milestone.Name)
		return NullID, err
	}

	var milestoneID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&milestoneID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new milestone %s", milestone.Name)
		return NullID, err
	}

	log.Debug("added milestone %s (id %d)", milestone.Name, milestoneID)

	return milestoneID, nil
}

// AddMilestone adds a milestone to Gitea, returns id of created milestone
func (accessor *DefaultAccessor) AddMilestone(milestone *Milestone) (int64, error) {
	milestoneID, err := accessor.GetMilestoneID(milestone.Name)
	if err != nil {
		return NullID, err
	}

	if milestoneID == NullID {
		return accessor.insertMilestone(milestone)
	}

	if accessor.overwrite {
		err = accessor.updateMilestone(milestoneID, milestone)
		if err != nil {
			return NullID, err
		}
	} else {
		log.Debug("milestone %s already exists - ignored", milestone.Name)
	}

	return milestoneID, nil
}

// GetMilestoneURL gets the URL for accessing a given milestone
func (accessor *DefaultAccessor) GetMilestoneURL(milestoneID int64) string {
	repoURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/milestone/%d", repoURL, milestoneID)
}
