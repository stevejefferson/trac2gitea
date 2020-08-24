// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/stevejefferson/trac2gitea/log"

// GetTypeNames retrieves all type names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetTypeNames(handlerFn func(typeName string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT type FROM ticket`)
	if err != nil {
		log.Error("Problem extracting names of trac types: %v\n", err)
		return err
	}

	for rows.Next() {
		var typ string
		if err := rows.Scan(&typ); err != nil {
			log.Error("Problem extracting name of trac type: %v\n", err)
			return err
		}

		err = handlerFn(typ)
		if err != nil {
			return err
		}
	}

	return nil
}
