package gitea

import (
	"database/sql"
	"fmt"
	"os"

	"stevejefferson.co.uk/trac2gitea/log"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3" // sqlite database driver
)

// Accessor provides access (retrieval and update) to Gitea (non-Wiki) data.
type Accessor struct {
	rootDir           string
	mainConfig        *ini.File
	customConfig      *ini.File
	db                *sql.DB
	userName          string
	repoName          string
	repoID            int64
	defaultAssigneeID int64
	defaultAuthorID   int64
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

// CreateAccessor returns a new Gitea accessor.
func CreateAccessor(giteaRootDir string, giteaUserName string, giteaRepoName string, defaultAssignee string, defaultAuthor string) *Accessor {
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

	giteaAccessor := Accessor{
		rootDir:           giteaRootDir,
		mainConfig:        giteaMainConfig,
		customConfig:      giteaCustomConfig,
		db:                nil,
		userName:          giteaUserName,
		repoName:          giteaRepoName,
		repoID:            0,
		defaultAssigneeID: 0,
		defaultAuthorID:   0}

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

	adminUserID := giteaAccessor.getAdminUserID()
	giteaDefaultAssigneeID := giteaAccessor.getAdminDefaultingUserID(defaultAssignee, adminUserID)
	giteaAccessor.defaultAssigneeID = giteaDefaultAssigneeID

	giteaDefaultAuthorID := giteaAccessor.getAdminDefaultingUserID(defaultAuthor, adminUserID)
	giteaAccessor.defaultAuthorID = giteaDefaultAuthorID

	return &giteaAccessor
}

func (accessor *Accessor) getUserRepoURL() string {
	baseURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", baseURL, accessor.userName, accessor.repoName)
}
