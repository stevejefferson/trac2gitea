// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/stevejefferson/trac2gitea/log"

// GetTickets retrieves all Trac tickets, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTickets(handlerFn func(
	ticketID int64, ticketType string, created int64,
	component string, severity string, priority string,
	owner string, reporter string, version string,
	milestone string, status string, resolution string,
	summary string, description string) error) error {
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
		log.Error("Problem retrieving trac tickets: %v\n", err)
		return err
	}

	for rows.Next() {
		var ticketID, created int64
		var component, ticketType, severity, priority, owner, reporter, version, milestone, status, resolution, summary, description string
		if err := rows.Scan(&ticketID, &ticketType, &created, &component, &severity, &priority, &owner, &reporter,
			&version, &milestone, &status, &resolution, &summary, &description); err != nil {
			log.Error("Problem extracting data on trac ticket: %v\n", err)
			return err
		}

		err = handlerFn(ticketID, ticketType, created, component, severity, priority, owner, reporter,
			version, milestone, status, resolution, summary, description)
		if err != nil {
			return err
		}
	}

	return nil
}
