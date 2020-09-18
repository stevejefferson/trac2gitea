// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetVersions retrieves all versions used in Trac, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetVersions(handlerFn func(version *Label) error) error {
	rows, err := accessor.db.Query(`
		SELECT DISTINCT COALESCE(version,'') name, '' desc FROM ticket
		UNION
		SELECT COALESCE(name,'') name, COALESCE(description,'') desc FROM version`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac versions")
		return err
	}

	for rows.Next() {
		var name, description string
		if err := rows.Scan(&name, &description); err != nil {
			err = errors.Wrapf(err, "retrieving Trac version")
			return err
		}

		version := Label{Name: name, Description: description}
		if err = handlerFn(&version); err != nil {
			return err
		}
	}

	return nil
}
