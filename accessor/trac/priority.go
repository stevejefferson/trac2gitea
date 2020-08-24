// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/stevejefferson/trac2gitea/log"

// GetPriorityNames retrieves all priority names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetPriorityNames(handlerFn func(priorityName string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT priority FROM ticket`)
	if err != nil {
		log.Error("Problem extracting names of trac priorities: %v\n", err)
		return err
	}

	for rows.Next() {
		var priorityName string
		if err := rows.Scan(&priorityName); err != nil {
			log.Error("Problem extracting name of trac priority: %v\n", err)
			return err
		}

		err = handlerFn(priorityName)
		if err != nil {
			return err
		}
	}

	return nil
}
