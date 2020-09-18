// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"

	"github.com/stevejefferson/trac2gitea/log"

	"github.com/pkg/errors"
)

// GetLabelID retrieves the id of the given label, returns NullID if no such label
func (accessor *DefaultAccessor) GetLabelID(labelName string) (int64, error) {
	var labelID int64 = NullID
	err := accessor.db.QueryRow(`
		SELECT id FROM label WHERE repo_id = $1 AND name = $2
		`, accessor.repoID, labelName).Scan(&labelID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id of label %s", labelName)
		return NullID, err
	}

	return labelID, nil
}

// updateLabel updates an existing label
func (accessor *DefaultAccessor) updateLabel(labelID int64, label *Label) error {
	_, err := accessor.db.Exec(`UPDATE label SET repo_id=?, name=?, description=?, color=? WHERE id=?`,
		accessor.repoID, label.Name, label.Description, label.Color, labelID)
	if err != nil {
		err = errors.Wrapf(err, "updating label %s", label.Name)
		return err
	}

	log.Debug("updated label %s, color %s (id %d)", label.Name, label.Color, labelID)

	return nil
}

// insertLabel inserts a new label, returns label id.
func (accessor *DefaultAccessor) insertLabel(label *Label) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO label(repo_id, name, description, color) VALUES($1, $2, $3)`,
		accessor.repoID, label.Name, label.Description, label.Color)
	if err != nil {
		err = errors.Wrapf(err, "adding label %s", label.Name)
		return NullID, err
	}

	var labelID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&labelID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new label %s", label.Name)
		return NullID, err
	}

	log.Debug("added label %s, color %s (id %d)", label.Name, label.Color, labelID)

	return labelID, nil
}

// AddLabel adds a label to Gitea, returns label id.
func (accessor *DefaultAccessor) AddLabel(label *Label) (int64, error) {
	labelID, err := accessor.GetLabelID(label.Name)
	if err != nil {
		return NullID, err
	}

	if labelID == NullID {
		return accessor.insertLabel(label)
	}

	if accessor.overwrite {
		err = accessor.updateLabel(labelID, label)
		if err != nil {
			return NullID, err
		}
	} else {
		log.Debug("label %s already exists - ignored", label.Name)
	}

	return labelID, nil
}
