// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetAttachmentUUID returns the UUID for a named attachment of a given issue - returns empty string if cannot find issue/attachment.
func (accessor *DefaultAccessor) GetAttachmentUUID(issueID int64, name string) (string, error) {
	var uuid string = ""
	err := accessor.db.QueryRow(`
			select uuid from attachment where issue_id = $1 and name = $2
			`, issueID, name).Scan(&uuid)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return "", err
	}

	return uuid, nil
}

// AddAttachment adds a new attachment to a given issue with the provided data - returns id of created attachment.
func (accessor *DefaultAccessor) AddAttachment(uuid string, issueID int64, commentID int64, attachmentName string, attachmentFile string, time int64) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, uuid, issueID, commentID, attachmentName, time)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	var attachmentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&attachmentID)
	if err != nil {
		log.Error(err)
		return -1, err
	}

	log.Debugf("Issue:%d, comment:%d : added attachment %s\n", issueID, commentID, attachmentName)

	giteaAttachmentsRootDir := accessor.GetStringConfig("attachment", "PATH")
	if giteaAttachmentsRootDir == "" {
		giteaAttachmentsRootDir = filepath.Join(accessor.rootDir, "data", "attachments")
	}

	d1 := uuid[0:1]
	d2 := uuid[1:2]
	giteaAttachmentsPath := filepath.Join(giteaAttachmentsRootDir, d1, d2, uuid)
	err = accessor.copyFile(attachmentFile, giteaAttachmentsPath)
	if err != nil {
		return -1, err
	}

	return attachmentID, nil
}

// GetAttachmentURL retrieves the URL for viewing a Gitea attachment
func (accessor *DefaultAccessor) GetAttachmentURL(uuid string) string {
	baseURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}
