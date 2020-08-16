package giteawiki

// Accessor is the interface for all of our interactions with our Gitea project Wiki.
type Accessor interface {
	/*
	 * Files
	 */
	// CopyFile copies an internal file into the Gitea Wiki, returning a URL through which the file can be viewed/
	CopyFile(externalFilePath string, giteaWikiRelPath string) string

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
