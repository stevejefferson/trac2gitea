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

// getIssueAttachmentIDandUUID retrieves the id and UUID of the given issue attachment, returns id of gitea.NullID if no such attachment
func (accessor *DefaultAccessor) getIssueAttachmentIDandUUID(issueID int64, fileName string) (int64, string, error) {
	var issueAttachmentID = NullID
	var issueAttachmentUUID string
	err := accessor.db.QueryRow(`
		SELECT id, uuid FROM attachment WHERE issue_id = $1 AND name = $2
		`, issueID, fileName).Scan(&issueAttachmentID, &issueAttachmentUUID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrapf(err, "retrieving id for attachment %s for issue %d", fileName, issueID)
		return NullID, "", err
	}

	return issueAttachmentID, issueAttachmentUUID, nil
}

// GetIssueAttachmentUUID returns the UUID for a named attachment of a given issue - returns empty string if cannot find issue/attachment.
func (accessor *DefaultAccessor) GetIssueAttachmentUUID(issueID int64, fileName string) (string, error) {
	_, uuid, err := accessor.getIssueAttachmentIDandUUID(issueID, fileName)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

// getAttachmentPath returns the path at which to store an attachment with a given UUID
func (accessor *DefaultAccessor) getAttachmentPath(UUID string) string {
	attachmentsRootDir := accessor.GetStringConfig("attachment", "PATH")
	if attachmentsRootDir == "" {
		attachmentsRootDir = filepath.Join(accessor.rootDir, "data", "attachments")
	}

	d1 := UUID[0:1]
	d2 := UUID[1:2]
	return filepath.Join(attachmentsRootDir, d1, d2, UUID)
}

// copyAttachment copies a given attachment file to the Gitea attachment with the given UUID
func (accessor *DefaultAccessor) copyAttachment(filePath string, UUID string) error {
	attachmentPath := accessor.getAttachmentPath(UUID)
	return copyFile(filePath, attachmentPath)
}

// deleteAttachment deletes the Gitea attachment with the given UUID
func (accessor *DefaultAccessor) deleteAttachment(UUID string) error {
	attachmentPath := accessor.getAttachmentPath(UUID)
	return deleteFile(attachmentPath)
}

// updateIssueAttachment updates an existing issue attachment
func (accessor *DefaultAccessor) updateIssueAttachment(issueAttachmentID int64, issueID int64, attachment *IssueAttachment, filePath string) error {
	_, err := accessor.db.Exec(`
		UPDATE attachment SET uuid=?, issue_id=?, comment_id=?, name=?, created_unix=? WHERE id=?`,
		attachment.UUID, issueID, attachment.CommentID, attachment.FileName, attachment.Time, issueAttachmentID)

	if err != nil {
		err = errors.Wrapf(err, "updating attachment %s for issue %d", attachment.FileName, issueID)
		return err
	}

	log.Debug("updated attachment %s for issue %d (id %d)", attachment.UUID, issueID, issueAttachmentID)

	return nil
}

// insertIssueAttachment creates a new attachment to a Gitea issue, returns id of created attachment
func (accessor *DefaultAccessor) insertIssueAttachment(issueID int64, attachment *IssueAttachment, filePath string) (int64, error) {
	_, err := accessor.db.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, attachment.UUID, issueID, attachment.CommentID, attachment.FileName, attachment.Time)
	if err != nil {
		err = errors.Wrapf(err, "adding attachment %s for issue %d", attachment.FileName, issueID)
		return NullID, err
	}

	var issueAttachmentID int64
	err = accessor.db.QueryRow(`SELECT last_insert_rowid()`).Scan(&issueAttachmentID)
	if err != nil {
		err = errors.Wrapf(err, "retrieving id of new attachment %s for issue %d", attachment.FileName, issueID)
		return NullID, err
	}

	log.Debug("added attachment %s for issue %d", attachment.FileName, issueID)

	return issueAttachmentID, nil
}

// AddIssueAttachment adds a new attachment to an issue using the provided file - returns id of created attachment
func (accessor *DefaultAccessor) AddIssueAttachment(issueID int64, attachment *IssueAttachment, filePath string) (int64, error) {
	issueAttachmentID, issueAttachmentUUID, err := accessor.getIssueAttachmentIDandUUID(issueID, attachment.FileName)
	if err != nil {
		return NullID, err
	}

	if issueAttachmentID == NullID {
		issueAttachmentID, err = accessor.insertIssueAttachment(issueID, attachment, filePath)
		if err != nil {
			return NullID, err
		}
	} else if accessor.overwrite {
		err = accessor.updateIssueAttachment(issueAttachmentID, issueID, attachment, filePath)
		if err != nil {
			return NullID, err
		}

		err = accessor.deleteAttachment(issueAttachmentUUID)
	} else {
		if attachment.UUID != issueAttachmentUUID {
			log.Warn("attachment %s already exists for issue %d but under UUID %s (expecting UUID %s)", attachment.FileName, issueID, issueAttachmentUUID, attachment.UUID)
		} else {
			log.Debug("issue %d already has attachment %s - ignored", issueID, attachment.FileName)
		}
		return issueAttachmentID, nil
	}

	err = accessor.copyAttachment(filePath, attachment.UUID)
	if err != nil {
		return NullID, err
	}

	return issueAttachmentID, nil
}

// GetIssueAttachmentURL retrieves the URL for viewing a Gitea attachment
func (accessor *DefaultAccessor) GetIssueAttachmentURL(issueID int64, uuid string) string {
	baseURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}
