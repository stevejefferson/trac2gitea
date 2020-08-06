package trac

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path"
)

func encodeSha1(str string) string {
	// Encode string to sha1 hex value.
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (accessor *Accessor) AttachmentPath(ticketID int64, fname string) string {
	ticketDir := encodeSha1(fmt.Sprintf("%d", ticketID))
	ticketSub := ticketDir[0:3]

	pathFile := encodeSha1(fname)
	pathExt := path.Ext(fname)

	return fmt.Sprintf("%s/attachments/ticket/%s/%s/%s%s", accessor.rootDir, ticketSub, ticketDir, pathFile, pathExt)
}
