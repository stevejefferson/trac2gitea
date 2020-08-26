// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
)

// GetLabelID retrieves the id of the given label, returns -1 if no such label
func (accessor *DefaultAccessor) GetLabelID(labelName string) (int64, error) {
	var labelID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM label WHERE repo_id = $1 AND name = $2
		`, accessor.repoID, labelName).Scan(&labelID)
	if err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return labelID, nil
}

// AddLabel adds a label to Gitea, returns label id
func (accessor *DefaultAccessor) AddLabel(labelName string, labelColor string) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO label(repo_id,name,color)
			SELECT $1,$2, $3 WHERE
			NOT EXISTS ( SELECT * FROM label WHERE repo_id = $1 AND name = $2 )`,
		accessor.repoID, labelName, labelColor)
	if err != nil {
		return -1, err
	}

	var labelID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&labelID)
	if err != nil {
		return -1, err
	}

	return labelID, nil
}
