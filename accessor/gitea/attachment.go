package gitea

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetAttachmentUUID returns the UUID for a named attachment of a given issue - returns empty string if cannot find issue/attachment.
func (accessor *Accessor) GetAttachmentUUID(issueID int64, name string) string {
	var uuid string = ""
	err := accessor.db.QueryRow(`
			select uuid from attachment where issue_id = $1 and name = $2
			`, issueID, name).Scan(&uuid)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return uuid
}

// AddAttachment adds a new attachment to a given issue with the provided data.
func (accessor *Accessor) AddAttachment(uuid string, issueID int64, commentID int64, attachmentName string, attachmentFile string, time int64) {
	_, err := accessor.db.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, uuid, issueID, commentID, attachmentName, time)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Issue:%d, comment:%d : added attachment %s\n", issueID, commentID, attachmentName)

	giteaAttachmentsRootDir := accessor.GetStringConfig("attachment", "PATH")
	if giteaAttachmentsRootDir == "" {
		giteaAttachmentsRootDir = filepath.Join(accessor.rootDir, "data", "attachments")
	}

	d1 := uuid[0:1]
	d2 := uuid[1:2]
	giteaAttachmentsPath := filepath.Join(giteaAttachmentsRootDir, d1, d2, uuid)
	accessor.copyFile(attachmentFile, giteaAttachmentsPath)
}

// GetAttachmentURL retrieves the URL for viewing a Gitea attachment
func (accessor *Accessor) GetAttachmentURL(uuid string) string {
	baseURL := accessor.getUserRepoURL()
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}
