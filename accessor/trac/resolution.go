// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetResolutions retrieves all resolutions used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetResolutions(handlerFn func(resolution *Label) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT resolution FROM ticket WHERE trim(resolution) != ''`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac resolutions")
		return err
	}

	for rows.Next() {
		var resolutionName string
		if err := rows.Scan(&resolutionName); err != nil {
			err = errors.Wrapf(err, "retrieving Trac resolution")
			return err
		}

		resolution := Label{Name: resolutionName, Description: ""}
		if err = handlerFn(&resolution); err != nil {
			return err
		}
	}

	return nil
}
