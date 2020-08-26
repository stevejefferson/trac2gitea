// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetTypeNames retrieves all type names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTypeNames(handlerFn func(typeName string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT type FROM ticket`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac types")
		return err
	}

	for rows.Next() {
		var typ string
		if err := rows.Scan(&typ); err != nil {
			err = errors.Wrapf(err, "retrieving Trac type")
			return err
		}

		if err = handlerFn(typ); err != nil {
			return err
		}
	}

	return nil
}
