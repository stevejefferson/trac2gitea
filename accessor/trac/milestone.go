// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetMilestones retrieves all Trac milestones, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetMilestones(
	handlerFn func(name string, description string, due int64, completed int64) error) error {
	// NOTE: trac timestamps are to the microseconds, we just need seconds
	rows, err := accessor.db.Query(`
		SELECT COALESCE(name,''), COALESCE(description,''), CAST(due*1e-6 AS int8), CAST(completed*1e-6 AS int8)
			FROM milestone UNION
			SELECT distinct(COALESCE(milestone,'')),'',0,0
				FROM ticket
				WHERE COALESCE(milestone,'') NOT IN ( select COALESCE(name,'') from milestone )`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac milestones")
		return err
	}

	for rows.Next() {
		var completed, due int64
		var name, description string
		if err := rows.Scan(&name, &description, &due, &completed); err != nil {
			err = errors.Wrapf(err, "retrieving Trac milestone")
			return err
		}

		if err = handlerFn(name, description, due, completed); err != nil {
			return err
		}
	}

	return nil
}
