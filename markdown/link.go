package markdown

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// @param prefix is the prefix to a link, like "ticket:5" for a link like "ticket:5:inside_the_nest.png"
func (converter *Converter) resolveTracLink(link, prefix string) string {
	// 'http...' links are left as-is
	if strings.HasPrefix(link, "http") {
		return link
	}

	// 'htdocs:...' links refer to trac htdocs directory
	if strings.HasPrefix(link, "htdocs:") {
		htdocsPath := strings.Replace(link, "htdocs:", "", -1)
		tracHtdocsPath := filepath.Join(converter.tracAccessor.RootDir, "htdocs", htdocsPath)
		wikiHtdocsRelPath := filepath.Join("htdocs", htdocsPath)
		converter.wikiAccessor.CopyFile(tracHtdocsPath, wikiHtdocsRelPath)
		return htdocsPath
	}

	// 'ticket:...' links refer to attachments
	if strings.HasPrefix(prefix, "ticket:") {
		ticketIDStr := strings.Replace(prefix, "ticket:", "", -1)
		ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		// Find issue id
		issueID := converter.giteaAccessor.GetIssueID(ticketID)
		uuid := converter.giteaAccessor.GetAttachmentUUID(issueID, link)
		return converter.giteaAccessor.AttachmentURL(uuid)
	}

	// TODO 'wiki:...' links
	fmt.Fprintf(os.Stderr, "WARNING: cannot resolve trac link %s with prefix '%s'", link, prefix)
	return link
}
