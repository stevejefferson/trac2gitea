package gitea

import "strings"

// GetWikiRepoName retrieves the name of the wiki repo associated with the current project.
func (accessor *Accessor) GetWikiRepoName() string {
	return accessor.userName + "/" + accessor.repoName + ".wiki"
}

// GetWikiRepoURL retrieves the URL of the wiki repo associated with the current project.
func (accessor *Accessor) GetWikiRepoURL() string {
	serverRootURL := accessor.GetStringConfig("server", "ROOT_URL")
	if !strings.HasSuffix(serverRootURL, "/") {
		serverRootURL = serverRootURL + "/"
	}
	return serverRootURL + accessor.GetWikiRepoName() + ".git"
}
