// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetVersionNames(handlerFn func(version string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(version,'') FROM ticket UNION SELECT COALESCE(name,'') FROM version`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}

		if err = handlerFn(version); err != nil {
			return err
		}
	}

	return nil
}
