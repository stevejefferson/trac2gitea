// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// getIssueAssigneeID retrieves the id of the given issue/assignee association, returns -1 if no such association
func (accessor *DefaultAccessor) getIssueAssigneeID(issueID int64, assigneeID int64) (int64, error) {
	var issueAssigneeID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM issue_assignees WHERE issue_id = $1 AND assignee_id = $2
		`, issueID, assigneeID).Scan(&issueAssigneeID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id for issue %d/assignee %d", issueID, assigneeID)
		return -1, err
	}

	return issueAssigneeID, nil
}

// updateIssueAssignee updates an existing issue assignee
func (accessor *DefaultAccessor) updateIssueAssignee(issueAssigneeID int64, issueID int64, assigneeID int64) error {
	_, err := accessor.db.Exec(`UPDATE issue_assignees SET issue_id=?, assignee_id=?, WHERE id= ?`, issueID, assigneeID, issueAssigneeID)
	if err != nil {
		err = errors.Wrapf(err, "updating issue %d/assignee %d", issueID, assigneeID)
		return err
	}

	log.Debug("updated assignee %d for issue %d (id %d)", assigneeID, issueID, issueAssigneeID)

	return nil
}

// insertIssueAssignee adds a new assignee to a Gitea issue
func (accessor *DefaultAccessor) insertIssueAssignee(issueID int64, assigneeID int64) error {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_assignees(issue_id, assignee_id) VALUES ($1, $2)`,
		issueID, assigneeID)
	if err != nil {
		err = errors.Wrapf(err, "adding user %d as assignee for issue id %d", assigneeID, issueID)
		return err
	}

	log.Debug("added assignee %d for issue %d", assigneeID, issueID)

	return nil
}

// AddIssueAssignee adds an assignee to a Gitea issue
func (accessor *DefaultAccessor) AddIssueAssignee(issueID int64, assigneeID int64) error {
	issueAssigneeID, err := accessor.getIssueAssigneeID(issueID, assigneeID)
	if err != nil {
		return err
	}

	if issueAssigneeID == -1 {
		return accessor.insertIssueAssignee(issueID, assigneeID)
	}

	if accessor.overwrite {
		err = accessor.updateIssueAssignee(issueAssigneeID, issueID, assigneeID)
		if err != nil {
			return err
		}
	} else {
		log.Debug("issue %d already has assignee %d - ignored", issueID, assigneeID)
	}

	return nil
}
