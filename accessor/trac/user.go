// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import "github.com/pkg/errors"

// GetUserMap returns a blank user mapping mapping for every user name found in Trac database fields to be converted
func (accessor *DefaultAccessor) GetUserMap() (map[string]string, error) {
	rows, err := accessor.db.Query(`
		SELECT owner FROM ticket
		UNION SELECT author FROM attachment
		UNION SELECT author FROM ticket_change
		UNION SELECT author FROM wiki`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac users to be converted")
		return nil, err
	}

	userMap := make(map[string]string)
	for rows.Next() {
		var userName string
		if err = rows.Scan(&userName); err != nil {
			err = errors.Wrapf(err, "retrieving Trac user to be converted")
			return nil, err
		}

		userMap[userName] = ""
	}

	return userMap, nil
}
