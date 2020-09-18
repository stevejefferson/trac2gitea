// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetComponents retrieves all Trac components, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetComponents(handlerFn func(component *Label) error) error {
	rows, err := accessor.db.Query(`SELECT name, COALESCE(description,'') FROM component`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac components")
		return err
	}

	for rows.Next() {
		var name, description string
		if err := rows.Scan(&name, &description); err != nil {
			err = errors.Wrapf(err, "retrieving Trac component")
			return err
		}

		component := Label{Name: name, Description: description}
		if err = handlerFn(&component); err != nil {
			return err
		}
	}

	return nil
}
