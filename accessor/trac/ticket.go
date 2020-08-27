// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetTickets retrieves all Trac tickets, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTickets(
	handlerFn func(ticketID int64, summary string, description string, owner string, reporter string, milestone string,
		component string, priority string, resolution string, severity string, typ string, version string,
		status string, created int64) error) error {
	rows, err := accessor.db.Query(`
		SELECT
			t.id,
			t.type,
			CAST(t.time*1e-6 AS int8),
			COALESCE(t.component, ''),
			COALESCE(t.severity,''),
			COALESCE(t.priority,''),
			COALESCE(t.owner,''),
			t.reporter,
			COALESCE(t.version,''),
			COALESCE(t.milestone,''),
			lower(COALESCE(t.status, '')),
			COALESCE(t.resolution,''),
			COALESCE(t.summary, ''),
			COALESCE(t.description, '')
		FROM ticket t ORDER BY id`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac tickets")
		return err
	}

	for rows.Next() {
		var ticketID, created int64
		var summary, description, owner, reporter, milestone, component, priority, resolution, severity, typ, version, status string
		if err := rows.Scan(&ticketID, &typ, &created, &component, &severity, &priority, &owner, &reporter,
			&version, &milestone, &status, &resolution, &summary, &description); err != nil {
			err = errors.Wrapf(err, "retrieving Trac ticket")
			return err
		}

		if err = handlerFn(ticketID, summary, description, owner, reporter, milestone,
			component, priority, resolution, severity, typ, version, status, created); err != nil {
			return err
		}
	}

	return nil
}
