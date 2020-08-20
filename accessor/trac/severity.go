// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetSeverityNames retrieves all severity names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetSeverityNames(handlerFn func(severityName string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(severity,'') FROM ticket`)
	if err != nil {
		log.Error(err)
		return err
	}

	for rows.Next() {
		var severityName string
		if err := rows.Scan(&severityName); err != nil {
			log.Error(err)
			return err
		}

		err = handlerFn(severityName)
		if err != nil {
			return err
		}
	}

	return nil
}
