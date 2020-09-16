// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/stevejefferson/trac2gitea/log"

	"github.com/pkg/errors"
)

// GetLabelID retrieves the id of the given label, returns -1 if no such label
func (accessor *DefaultAccessor) GetLabelID(labelName string) (int64, error) {
	var labelID int64 = -1
	err := accessor.db.QueryRow(`
		SELECT id FROM label WHERE repo_id = $1 AND name = $2
		`, accessor.repoID, labelName).Scan(&labelID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of label %s", labelName)
		return -1, err
	}

	return labelID, nil
}

// updateLabel updates an existing label
func (accessor *DefaultAccessor) updateLabel(labelID int64, labelName string, labelColor string) error {
	_, err := accessor.db.Exec(`UPDATE label SET repo_id=?, name=?, color=? WHERE id=?`,
		accessor.repoID, labelName, labelColor, labelID)
	if err != nil {
		err = errors.Wrapf(err, "updating label %s", labelName)
		return err
	}

	log.Debug("updated label %s, color %s (id %d)", labelName, labelColor, labelID)

	return nil
}

// insertLabel inserts a new label, returns label id.
func (accessor *DefaultAccessor) insertLabel(labelName string, labelColor string) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO label(repo_id, name, color) VALUES($1, $2, $3)`,
		accessor.repoID, labelName, labelColor)
	if err != nil {
		err = errors.Wrapf(err, "adding label %s", labelName)
		return -1, err
	}

	var labelID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&labelID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new label %s", labelName)
		return -1, err
	}

	log.Debug("added label %s, color %s (id %d)", labelName, labelColor, labelID)

	return labelID, nil
}

// AddLabel adds a label to Gitea, returns label id.
func (accessor *DefaultAccessor) AddLabel(labelName string, labelColor string) (int64, error) {
	labelID, err := accessor.GetLabelID(labelName)
	if err != nil {
		return -1, err
	}

	if labelID == -1 {
		return accessor.insertLabel(labelName, labelColor)
	}

	if accessor.overwrite {
		err = accessor.updateLabel(labelID, labelName, labelColor)
		if err != nil {
			return -1, err
		}
	} else {
		log.Debug("label %s already exists - ignored", labelName)
	}

	return labelID, nil
}
