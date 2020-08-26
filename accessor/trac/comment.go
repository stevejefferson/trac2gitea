// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import (
	"database/sql"

	"github.com/pkg/errors"
)

// GetComments retrieves all comments on a given Trac ticket, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetComments(
	ticketID int64,
	handlerFn func(ticketID int64, time int64, author string, comment string) error) error {
	rows, err := accessor.db.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, COALESCE(newvalue, '') newval
			FROM ticket_change where ticket = $1 AND field = 'comment' AND trim(COALESCE(newvalue, ''), ' ') != ''
			ORDER BY time asc`, ticketID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac comments for ticket %d", ticketID)
		return err
	}

	for rows.Next() {
		var time int64
		var author, comment string
		if err := rows.Scan(&time, &author, &comment); err != nil {
			err = errors.Wrapf(err, "retrieving Trac comment for ticket %d", ticketID)
			return err
		}

		if err = handlerFn(ticketID, time, author, comment); err != nil {
			return err
		}
	}

	return nil
}

// GetCommentString retrieves a given comment string for a given Trac ticket
func (accessor *DefaultAccessor) GetCommentString(ticketID int64, commentNum int64) (string, error) {
	var commentStr string
	err := accessor.db.QueryRow(`
		SELECT COALESCE(newvalue, '') FROM ticket_change where ticket = $1 AND oldvalue = $2 AND field = 'comment'`,
		ticketID, commentNum).Scan(&commentStr)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving Trac comment number %d for ticket %d", commentNum, ticketID)
		return "", err
	}

	return commentStr, nil
}
