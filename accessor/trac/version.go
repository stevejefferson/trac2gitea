// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/stevejefferson/trac2gitea/log"

// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetVersionNames(handlerFn func(version string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(version,'') FROM ticket UNION SELECT COALESCE(name,'') FROM version`)
	if err != nil {
		log.Error("Problem extracting names of trac versions: %v\n", err)
		return err
	}

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			log.Error("Problem extracting name of trac version: %v\n", err)
			return err
		}

		err = handlerFn(version)
		if err != nil {
			return err
		}
	}

	return nil
}
