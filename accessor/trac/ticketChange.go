// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import (
	"database/sql"

	"github.com/pkg/errors"
)

// GetTicketChanges retrieves all changes on a given Trac ticket in ascending time order, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTicketChanges(ticketID int64, handlerFn func(change *TicketChange) error) error {
	rows, err := accessor.db.Query(`
		SELECT CAST(time*1e-6 AS int8), field, COALESCE(author, ''), COALESCE(newvalue, ''), COALESCE(oldvalue, '')
			FROM ticket_change
			WHERE ticket = $1
			AND field IN ('comment', 'owner') 
			AND trim(COALESCE(newvalue, ''), ' ') != ''
			ORDER BY time asc`, ticketID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac comments for ticket %d", ticketID)
		return err
	}

	for rows.Next() {
		var time int64
		var field, author, newValue, oldValue string
		if err := rows.Scan(&time, &field, &author, &newValue, &oldValue); err != nil {
			err = errors.Wrapf(err, "retrieving Trac change for ticket %d", ticketID)
			return err
		}

		change := TicketChange{TicketID: ticketID, Author: author, Time: time}
		switch field {
		case "comment":
			comment := TicketComment{Text: newValue}
			change.ChangeType = TicketCommentChange
			change.Comment = &comment
		case "owner":
			ownership := TicketOwnership{PrevOwner: oldValue, Owner: newValue}
			change.ChangeType = TicketOwnershipChange
			change.Ownership = &ownership
		}

		if err = handlerFn(&change); err != nil {
			return err
		}
	}

	return nil
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
