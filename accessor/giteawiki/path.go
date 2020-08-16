package giteawiki

import "path/filepath"

// GetAttachmentRelPath returns the location of an attachment to Trac a wiki page when stored in the Gitea wiki repository.
// The returned path is relative to the root of the Gitea wiki repository.
func (accessor *DefaultAccessor) GetAttachmentRelPath(pageName string, filename string) string {
	return filepath.Join("attachments", pageName, filename)
}

// GetHtdocRelPath returns the location of a given Trac 'htdocs' file when stored in the Gitea wiki repository.
// The returned path is relative to the root of the Gitea wiki repository.
func (accessor *DefaultAccessor) GetHtdocRelPath(filename string) string {
	return filepath.Join("htdocs", filename)
}

// GetFileURL returns a URL for viewing a file stored in the Gitea wiki repository.
func (accessor *DefaultAccessor) GetFileURL(relpath string) string {
	return "../raw/" + relpath
}
