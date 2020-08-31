// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"
)

// GetIssueAttachmentUUID returns the UUID for a named attachment of a given issue - returns empty string if cannot find issue/attachment.
func (accessor *DefaultAccessor) GetIssueAttachmentUUID(issueID int64, fileName string) (string, error) {
	var uuid string = ""
	err := accessor.db.QueryRow(
		`select uuid from attachment where issue_id = $1 and name = $2`,
		issueID, fileName).Scan(&uuid)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving attachment %s for issue %d", fileName, issueID)
		return "", err
	}

	return uuid, nil
}

// AddIssueAttachment adds a new attachment to an issue using the provided file - returns id of created attachment
func (accessor *DefaultAccessor) AddIssueAttachment(issueID int64, fileName string, attachment *IssueAttachment) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, attachment.UUID, issueID, attachment.CommentID, fileName, attachment.Time)
	if err != nil {
		err = errors.Wrapf(err, "adding attachment %s for issue %d", fileName, issueID)
		return -1, err
	}

	var attachmentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&attachmentID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new attachment %s for issue %d", fileName, issueID)
		return -1, err
	}

	log.Debug("issue:%d, comment:%d : added attachment %s", issueID, attachment.CommentID, fileName)

	giteaAttachmentsRootDir := accessor.GetStringConfig("attachment", "PATH")
	if giteaAttachmentsRootDir == "" {
		giteaAttachmentsRootDir = filepath.Join(accessor.rootDir, "data", "attachments")
	}

	d1 := attachment.UUID[0:1]
	d2 := attachment.UUID[1:2]
	giteaAttachmentsPath := filepath.Join(giteaAttachmentsRootDir, d1, d2, attachment.UUID)
	err = accessor.copyFile(attachment.FilePath, giteaAttachmentsPath)
	if err != nil {
		return -1, err
	}

	return attachmentID, nil
}

// GetIssueAttachmentURL retrieves the URL for viewing a Gitea attachment
func (accessor *DefaultAccessor) GetIssueAttachmentURL(uuid string) string {
	baseURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}
