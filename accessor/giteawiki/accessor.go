package giteawiki

// Accessor is the interface for all of our interactions with our Gitea project Wiki.
type Accessor interface {
	/*
	 * Files
	 */
	// CopyFile copies an internal file into the Gitea Wiki
	CopyFile(externalFilePath string, giteaWikiRelPath string)

	/*
	 * Repository paths
	 */
	// GetAttachmentRelPath returns the location of an attachment to Trac a wiki page when stored in the Gitea wiki repository.
	// The returned path is relative to the root of the Gitea wiki repository.
	GetAttachmentRelPath(pageName string, filename string) string

	// GetHtdocRelPath returns the location of a given Trac 'htdocs' file when stored in the Gitea wiki repository.
	// The returned path is relative to the root of the Gitea wiki repository.
	GetHtdocRelPath(filename string) string

	// GetFileURL returns a URL for viewing a file stored in the Gitea wiki repository.
	GetFileURL(relpath string) string

	/*
	 * Wiki Repository
	 */
	// RepoClone clones our wiki repo to the provided directory.
	RepoClone()

	// RepoStageAndCommit stages any files added or updated since the last commit then commits them to our cloned wiki repo.
	RepoStageAndCommit(author string, authorEMail string, message string)

	// RepoComplete indicates that changes to the local wiki repository are complete.
	RepoComplete()

	/*
	 * Wiki Pages
	 */
	// WritePage writes (a version of) a wiki page to the checked-out wiki repository, returning the path to the written file.
	WritePage(pageName string, markdownText string) string

	// TranslatePageName translates a Trac wiki page name into a Gitea one
	TranslatePageName(pageName string) string
}
