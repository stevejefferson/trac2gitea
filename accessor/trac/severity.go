// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetSeverities retrieves all severities used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetSeverities(handlerFn func(severity *Label) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(severity,'') FROM ticket`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac severities")
		return err
	}

	for rows.Next() {
		var severityName string
		if err := rows.Scan(&severityName); err != nil {
			err = errors.Wrapf(err, "retrieving Trac severity")
			return err
		}

		severity := Label{Name: severityName, Description: ""}
		if err = handlerFn(&severity); err != nil {
			return err
		}
	}

	return nil
}
