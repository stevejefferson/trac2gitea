// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// sqlForFieldList returns a bracketted SQL format list (as used in a 'IN' expression) containing the provided field names.
func sqlForFieldList(fields []TicketChangeType) string {
	fieldListSQL := ""
	for fieldIndex, field := range fields {
		if fieldIndex > 0 {
			fieldListSQL = fieldListSQL + `, `
		}
		fieldListSQL = fieldListSQL + `'` + string(field) + `'`
	}

	return "(" + fieldListSQL + ")"
}

// sqlForFirstChange returns the SQL for retrieving details of the first change to a ticket
func sqlForFirstChange() string {
	// note: this query requires additional protection against NULL for case where no change exists for ticket
	changeSQL := `
		SELECT 0 source,
			COALESCE(chg0.field, '') field,
			COALESCE(chg0.oldvalue, '') value,
			COALESCE(chg0.author, '') author,
			COALESCE(MIN(CAST(chg0.time*1e-6 AS int8)),0) time
		FROM ticket_change chg0
		WHERE chg0.ticket = $1
		`
	return changeSQL
}

// sqlForFirstChangeToEachField returns the SQL for retrieving details of the first change to each of a set of fields of a ticket
func sqlForFirstChangeToEachField(fields []TicketChangeType) string {
	changeSQL := `
		SELECT 1 source,
			chg1.field field,
			COALESCE(chg1.oldvalue, '') value,
			COALESCE(chg1.author, '') author,
			MIN(CAST(chg1.time*1e-6 AS int8)) time
		FROM ticket_change chg1
		WHERE chg1.ticket = $1
		AND chg1.field IN ` + sqlForFieldList(fields) + `
		GROUP BY chg1.field
		`

	return changeSQL
}

// sqlForTicketTableField returns the SQL for retrieving each of set of fields of a ticket in the same format as the the other queries above
func sqlForTicketTableField(fieldIndex int, field TicketChangeType) string {
	table := fmt.Sprintf("t%d", fieldIndex)
	strField := string(field)
	return `
		SELECT 2 source,
			'` + strField + `' field,
			COALESCE(` + table + `.` + strField + `, '') value,
			COALESCE(` + table + `.owner, '') author,
			CAST(` + table + `.time*1e-6 AS int8) time
		FROM ticket ` + table + `
		WHERE id=$1
		`
}

var initialTicketChangeFields = []TicketChangeType{
	TicketComponentChange, TicketMilestoneChange, TicketOwnerChange, TicketPriorityChange,
	TicketResolutionChange, TicketSeverityChange, TicketTypeChange, TicketVersionChange,
}

// getInitialTicketChanges generates a set of "synthetic" changes on a Trac ticket to model the assignments of its initial values
// and passes the resultant data to a "handler" function.
// This is necessary because Trac tickets can be assigned certain values (e.g. severity, type) on creation in contrast to Gitea where
// these assignments must occur as specific issue changes. The synthetic changes here are used to trigger those Gitea issue changes.
func (accessor *DefaultAccessor) getInitialTicketChanges(ticketID int64, handlerFn func(change *TicketChange) error) error {
	// The initial value for any field is found as follows:
	// 1. in the 'oldvalue' of the earliest recorded change for the field
	// 2. (if no recorded change for field) in the appropriate field of the main ticket table
	//
	// We retrieve this using a complex SQL query which returns the following columns:
	// - source: "source" of data - 0 (see below), 1 or 2 (see above)
	// - field: name of field
	// - value: initial value for field
	// - author: author for field
	// - time: timestamp for value (converted to Gitea precision)
	//
	// In addition to sources 1 and 2 above we also return a source 0 which contains the earliest of all changes to the ticket.
	// If this exists, its author and time values are used to sign any data returned from the main ticket table (source 2 above).
	// This is because the author and time in the main ticket table refer to the latest change to the ticket table not the earliest.
	initialValueSQL := `
		SELECT source, field, value, author, time FROM 
		(`
	initialValueSQL = initialValueSQL + sqlForFirstChange() + `UNION` + sqlForFirstChangeToEachField(initialTicketChangeFields)
	for fieldIndex, field := range initialTicketChangeFields {
		initialValueSQL = initialValueSQL + `UNION` + sqlForTicketTableField(fieldIndex, field)
	}
	initialValueSQL = initialValueSQL + `) ORDER BY source asc`

	rows, err := accessor.db.Query(initialValueSQL, ticketID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving initial changes for ticket %d", ticketID)
		return err
	}

	var firstChangeAuthor = ""
	var firstChangeTime = int64(0)
	var source1Fields = make(map[string]bool)
	for rows.Next() {
		var source int
		var time int64
		var field, author, value string
		if err := rows.Scan(&source, &field, &value, &author, &time); err != nil {
			err = errors.Wrapf(err, "retrieving initial change for ticket %d", ticketID)
			return err
		}

		switch source {
		case 0:
			// source 0 change: this just gives us the author and time of the first change
			// - where there a ticket has no changes, author will be "" and time 0 - the effect of this will be harmless
			firstChangeAuthor = author
			firstChangeTime = time
			continue
		case 1:
			// source 1 change: earliest recorded change for a given field
			// this type of change takes precedence over changes from the main ticket table (source 2)
			// - note down the field so that we can ignore any source 2 changes for it
			source1Fields[field] = true
		case 2:
			// source 2 change: data from main ticket table
			// - this is only important if we haven't found an explicit recorded change for our field (source 1)
			if source1Fields[field] {
				continue
			}

			// replace author and time with those from first change (if any)
			if firstChangeAuthor != "" {
				author = firstChangeAuthor
			}
			if firstChangeTime > 0 && firstChangeTime < time {
				time = firstChangeTime
			}
		}

		// if initial value of field is "" then we don't need to generate a change as this would be the default value under Gitea
		if value == "" {
			continue
		}

		// record a change from "" to the initial value
		change := TicketChange{
			TicketID:   ticketID,
			ChangeType: TicketChangeType(field),
			Author:     author,
			OldValue:   "",
			NewValue:   value,
			Time:       time,
		}
		if err = handlerFn(&change); err != nil {
			return err
		}
	}

	return nil
}

var recordedTicketChangeFields = []TicketChangeType{
	TicketComponentChange, TicketMilestoneChange, TicketOwnerChange, TicketPriorityChange, TicketResolutionChange,
	TicketSeverityChange, TicketStatusChange, TicketSummaryChange, TicketTypeChange, TicketVersionChange,
}

// getRecordedTicketChanges retrieves all changes on a given ticket recorded by Trac in ascending time order, passing data from each to a "handler" function.
func (accessor *DefaultAccessor) getRecordedTicketChanges(ticketID int64, handlerFn func(change *TicketChange) error) error {
	rows, err := accessor.db.Query(`
		SELECT field, COALESCE(author, ''), COALESCE(oldvalue, ''), COALESCE(newvalue, ''), CAST(time*1e-6 AS int8)
			FROM ticket_change
			WHERE ticket = $1
			AND (
				(field = '`+string(TicketCommentChange)+`' AND trim(COALESCE(newvalue, ''), ' ') != '')
				OR field IN `+sqlForFieldList(recordedTicketChangeFields)+`
			)
			ORDER BY time asc`,
		ticketID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac comments for ticket %d", ticketID)
		return err
	}

	for rows.Next() {
		var time int64
		var field, author, oldValue, newValue string
		if err := rows.Scan(&field, &author, &oldValue, &newValue, &time); err != nil {
			err = errors.Wrapf(err, "retrieving Trac change for ticket %d", ticketID)
			return err
		}

		change := TicketChange{
			TicketID:   ticketID,
			ChangeType: TicketChangeType(field),
			Author:     author,
			OldValue:   oldValue,
			NewValue:   newValue,
			Time:       time}

		if err = handlerFn(&change); err != nil {
			return err
		}
	}

	return nil
}

// GetTicketChanges retrieves all changes on a given Trac ticket in ascending time order, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTicketChanges(ticketID int64, handlerFn func(change *TicketChange) error) error {
	err := accessor.getInitialTicketChanges(ticketID, handlerFn)
	if err != nil {
		return err
	}
	return accessor.getRecordedTicketChanges(ticketID, handlerFn)
}

// GetTicketCommentTime retrieves the timestamp for a given comment for a given Trac ticket
func (accessor *DefaultAccessor) GetTicketCommentTime(ticketID int64, commentNum int64) (int64, error) {
	timestamp := int64(0)
	err := accessor.db.QueryRow(`
		SELECT time FROM ticket_change where ticket = $1 AND oldvalue = $2 AND field = 'comment'`,
		ticketID, commentNum).Scan(&timestamp)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving Trac comment number %d for ticket %d", commentNum, ticketID)
		return 0, err
	}

	return timestamp, nil
}
