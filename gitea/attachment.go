package gitea

import (
	"fmt"
	"log"
)

// GetAttachmentUUID returns the UUID for a named attachment of a given issue.
func (accessor *Accessor) GetAttachmentUUID(issueID int64, name string) string {
	var uuid string
	err := accessor.db.QueryRow(`
			select uuid from attachment where issue_id = $1 and name = $2
			`, issueID, name).Scan(&uuid)
	if err != nil {
		log.Fatal(err)
	}

	return uuid
}

// AddAttachment adds a new attachment to a given issue with the provided data.
func (accessor *Accessor) AddAttachment(uuid string, issueID int64, commentID int64, fname string, time int64) {
	_, err := accessor.db.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, uuid, issueID, commentID, fname, time)
	if err != nil {
		log.Fatal(err)
	}
}

// AttachmentURL retrieves URL for Gitea attachment
func (accessor *Accessor) AttachmentURL(uuid string) string {
	baseURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}

// AttachmentRelativePath retrieves relative path of attachment
func (accessor *Accessor) AttachmentRelativePath(uuid string) string {
	d1 := uuid[0:1]
	d2 := uuid[1:2]
	// TODO: seek for PATH under [attachment] in Gitea config
	subpath := "data/attachments"
	return fmt.Sprintf("%s/%s/%s/%s", subpath, d1, d2, uuid)
}
