// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetTypes retrieves all types used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTypes(handlerFn func(tracType *Label) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT type FROM ticket`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac types")
		return err
	}

	for rows.Next() {
		var typeName string
		if err := rows.Scan(&typeName); err != nil {
			err = errors.Wrapf(err, "retrieving Trac type")
			return err
		}

		tracType := Label{Name: typeName, Description: ""}
		if err = handlerFn(&tracType); err != nil {
			return err
		}
	}

	return nil
}
