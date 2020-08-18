package gitea

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3" // sqlite database driver
	"gopkg.in/src-d/go-git.v4"
)

// DefaultAccessor is the default implementation of the gitea Accessor interface, accessing Gitea directly via its database and filestore.
type DefaultAccessor struct {
	rootDir           string
	mainConfig        *ini.File
	customConfig      *ini.File
	db                *sql.DB
	userName          string
	repoName          string
	repoID            int64
	defaultAssigneeID int64
	defaultAuthorID   int64
	wikiRepoURL       string
	wikiRepoDir       string
	wikiRepo          *git.Repository
}

func fetchConfig(configPath string) *ini.File {
	_, err := os.Stat(configPath)
	if err != nil {
		return nil
	}

	config, err := ini.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

// CreateDefaultAccessor returns a new Gitea default accessor.
func CreateDefaultAccessor(
	giteaRootDir string,
	giteaUserName string,
	giteaRepoName string,
	giteaWikiRepoURL string,
	giteaWikiRepoDir string,
	defaultAssignee string,
	defaultAuthor string) *DefaultAccessor {
	stat, err := os.Stat(giteaRootDir)
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatalf("Gitea root path %s is not a directory\n", giteaRootDir)
	}

	giteaMainConfigPath := "/etc/gitea/conf/app.ini"
	giteaMainConfig := fetchConfig(giteaMainConfigPath)
	giteaCustomConfigPath := fmt.Sprintf("%s/custom/conf/app.ini", giteaRootDir)
	giteaCustomConfig := fetchConfig(giteaCustomConfigPath)
	if giteaMainConfig == nil && giteaCustomConfig == nil {
		log.Fatalf("cannot find Gitea config in %s or %s\n", giteaMainConfigPath, giteaCustomConfigPath)
	}

	giteaAccessor := DefaultAccessor{
		rootDir:           giteaRootDir,
		mainConfig:        giteaMainConfig,
		customConfig:      giteaCustomConfig,
		db:                nil,
		userName:          giteaUserName,
		repoName:          giteaRepoName,
		repoID:            0,
		defaultAssigneeID: 0,
		defaultAuthorID:   0,
		wikiRepoURL:       "",
		wikiRepoDir:       "",
		wikiRepo:          nil}

	// extract path to gitea DB - currently sqlite-specific...
	giteaDbPath := giteaAccessor.GetStringConfig("database", "PATH")
	giteaDb, err := sql.Open("sqlite3", giteaDbPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Using Gitea database %s\n", giteaDbPath)
	giteaAccessor.db = giteaDb

	giteaRepoID := giteaAccessor.getRepoID(giteaUserName, giteaRepoName)
	if giteaRepoID == -1 {
		log.Fatalf("Cannot find repository %s for user %s\n", giteaRepoName, giteaUserName)
	}
	giteaAccessor.repoID = giteaRepoID

	// work out user ids
	adminUserID := giteaAccessor.getAdminUserID()
	giteaDefaultAssigneeID := giteaAccessor.getAdminDefaultingUserID(defaultAssignee, adminUserID)
	giteaAccessor.defaultAssigneeID = giteaDefaultAssigneeID

	giteaDefaultAuthorID := giteaAccessor.getAdminDefaultingUserID(defaultAuthor, adminUserID)
	giteaAccessor.defaultAuthorID = giteaDefaultAuthorID

	// find directory into which to clone wiki
	wikiRepoName := giteaUserName + "/" + giteaRepoName + ".wiki"
	if giteaWikiRepoDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		giteaWikiRepoDir = filepath.Join(cwd, wikiRepoName)
	}
	giteaAccessor.wikiRepoDir = giteaWikiRepoDir

	// find URL from which clone wiki
	if giteaWikiRepoURL == "" {
		serverRootURL := giteaAccessor.GetStringConfig("server", "ROOT_URL")
		if !strings.HasSuffix(serverRootURL, "/") {
			serverRootURL = serverRootURL + "/"
		}
		giteaWikiRepoURL = serverRootURL + wikiRepoName + ".git"
	}
	giteaAccessor.wikiRepoURL = giteaWikiRepoURL

	return &giteaAccessor
}

func (accessor *DefaultAccessor) getUserRepoURL() string {
	baseURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", baseURL, accessor.userName, accessor.repoName)
}
