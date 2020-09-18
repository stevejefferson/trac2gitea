// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetPriorities retrieves all priorities used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetPriorities(handlerFn func(priority *Label) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT priority FROM ticket`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac priorities")
		return err
	}

	for rows.Next() {
		var priorityName string
		if err := rows.Scan(&priorityName); err != nil {
			err = errors.Wrapf(err, "retrieving Trac priority")
			return err
		}

		priority := Label{Name: priorityName, Description: ""}

		if err = handlerFn(&priority); err != nil {
			return err
		}
	}

	return nil
}
