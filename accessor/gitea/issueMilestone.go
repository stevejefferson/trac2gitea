// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"github.com/pkg/errors"
)

// UpdateMilestoneIssueCounts updates issue counts for all milestones.
func (accessor *DefaultAccessor) UpdateMilestoneIssueCounts() error {
	_, err := accessor.db.Exec(`
		UPDATE milestone AS m SET 
			num_issues = (
				SELECT COUNT(i1.id)
				FROM issue i1
				WHERE m.id = i1.milestone_id
				GROUP BY i1.milestone_id),
			num_closed_issues = (
				SELECT COUNT(i2.id)
				FROM issue i2
				WHERE m.id = i2.milestone_id
				AND i2.is_closed = 1
				GROUP BY i2.milestone_id)`)
	if err != nil {
		err = errors.Wrapf(err, "updating number of issues for milestones")
		return err
	}

	return nil
}
