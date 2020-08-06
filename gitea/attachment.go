package gitea

import (
	"fmt"
	"log"
)

func (accessor *Accessor) AddAttachment(uuid string, issueID int64, commentID int64, fname string, time int64) {
	_, err := accessor.db.Exec(`
		INSERT INTO attachment(
			uuid, issue_id, comment_id, name, created_unix)
			VALUES ($1, $2, $3, $4, $5)`, uuid, issueID, commentID, fname, time)
	if err != nil {
		log.Fatal(err)
	}
}

func (accessor *Accessor) AttachmentURL(uuid string) string {
	baseURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/attachments/%s", baseURL, uuid)
}

func (accessor *Accessor) AttachmentRelativePath(uuid string) string {
	d1 := uuid[0:1]
	d2 := uuid[1:2]
	// TODO: seek for PATH under [attachment]
	//       in giteaRootPath/custom/conf/app.ini
	subpath := "data/attachments"
	return fmt.Sprintf("%s/%s/%s/%s", subpath, d1, d2, uuid)
}
