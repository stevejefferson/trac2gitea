// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

// GetPriorityNames retrieves all priority names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetPriorityNames(handlerFn func(priorityName string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT priority FROM ticket`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var priorityName string
		if err := rows.Scan(&priorityName); err != nil {
			return err
		}

		if err = handlerFn(priorityName); err != nil {
			return err
		}
	}

	return nil
}
