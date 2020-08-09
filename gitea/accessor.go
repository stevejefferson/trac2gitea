package gitea

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3" // sqlite database driver
)

// Accessor provides acess (retire and update) to Gitea (non-Wiki) data.
type Accessor struct {
	rootDir           string
	mainConfig        *ini.File
	customConfig      *ini.File
	db                *sql.DB
	userName          string
	repoName          string
	repoID            int64
	DefaultAssigneeID int64
	DefaultAuthorID   int64
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
		log.Fatal("Gitea root path is not a directory: ", giteaRootDir)
	}

	giteaMainConfigPath := "/etc/gitea/conf/app.ini"
	giteaMainConfig := fetchConfig(giteaMainConfigPath)
	giteaCustomConfigPath := fmt.Sprintf("%s/custom/conf/app.ini", giteaRootDir)
	giteaCustomConfig := fetchConfig(giteaCustomConfigPath)
	if giteaMainConfig == nil && giteaCustomConfig == nil {
		log.Fatal("cannot find Gitea config in " + giteaMainConfigPath + " or " + giteaCustomConfigPath)
	}

	giteaAccessor := Accessor{
		rootDir:           giteaRootDir,
		mainConfig:        giteaMainConfig,
		customConfig:      giteaCustomConfig,
		db:                nil,
		userName:          giteaUserName,
		repoName:          giteaRepoName,
		repoID:            0,
		DefaultAssigneeID: 0,
		DefaultAuthorID:   0}

	// extract path to gitea DB - currently sqlite-specific...
	giteaDbPath := giteaAccessor.GetStringConfig("database", "PATH")
	giteaDb, err := sql.Open("sqlite3", giteaDbPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using Gitea database %s\n", giteaDbPath)
	giteaAccessor.db = giteaDb

	giteaRepoID := giteaAccessor.findRepoID(giteaUserName, giteaRepoName)
	giteaAccessor.repoID = giteaRepoID

	adminUserID := giteaAccessor.findAdminUserID()
	giteaDefaultAssigneeID := giteaAccessor.findAdminDefaultingUserID(defaultAssignee, adminUserID)
	giteaAccessor.DefaultAssigneeID = giteaDefaultAssigneeID

	giteaDefaultAuthorID := giteaAccessor.findAdminDefaultingUserID(defaultAuthor, adminUserID)
	giteaAccessor.DefaultAuthorID = giteaDefaultAuthorID

	return &giteaAccessor
}

func (accessor *Accessor) getUserRepoURL() string {
	baseURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", baseURL, accessor.userName, accessor.repoName)
}
