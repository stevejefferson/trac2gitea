// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetTickets retrieves all Trac tickets, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTickets(handlerFn func(ticket *Ticket) error) error {
	rows, err := accessor.db.Query(`
		SELECT
			t.id,
			t.type,
			CAST(t.time*1e-6 AS int8),
			CAST(t.changetime*1e-6 AS int8),
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
		var ticketID, created, updated int64
		var summary, description, owner, reporter, milestoneName, componentName, priorityName, resolutionName, severityName, typeName, versionName, status string
		if err := rows.Scan(&ticketID, &typeName, &created, &updated, &componentName, &severityName, &priorityName, &owner, &reporter,
			&versionName, &milestoneName, &status, &resolutionName, &summary, &description); err != nil {
			err = errors.Wrapf(err, "retrieving Trac ticket")
			return err
		}

		ticket := Ticket{TicketID: ticketID, Summary: summary, Description: description, Owner: owner, Reporter: reporter,
			MilestoneName: milestoneName, ComponentName: componentName, PriorityName: priorityName, ResolutionName: resolutionName,
			SeverityName: severityName, TypeName: typeName, VersionName: versionName, Status: status, Created: created, Updated: updated}

		if err = handlerFn(&ticket); err != nil {
			return err
		}
	}

	return nil
}
