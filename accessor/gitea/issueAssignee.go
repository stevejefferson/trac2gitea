// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import "github.com/pkg/errors"

// AddIssueAssignee adds an assignee to a Gitea issue
func (accessor *DefaultAccessor) AddIssueAssignee(issueID int64, assigneeID int64) error {
	_, err := accessor.db.Exec(`
		INSERT INTO issue_assignees(issue_id, assignee_id) VALUES ($1, $2)`,
		issueID, assigneeID)
	if err != nil {
		err = errors.Wrapf(err, "adding user %d as assignee for issue id %d", assigneeID, issueID)
		return err
	}

	return nil
}
