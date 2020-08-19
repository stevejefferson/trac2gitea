package gitea

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

func fetchConfig(configPath string) (*ini.File, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, nil
	}

	config, err := ini.Load(configPath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return config, nil
}

// CreateDefaultAccessor returns a new Gitea default accessor.
func CreateDefaultAccessor(
	giteaRootDir string,
	giteaUserName string,
	giteaRepoName string,
	giteaWikiRepoURL string,
	giteaWikiRepoDir string,
	defaultAssignee string,
	defaultAuthor string) (*DefaultAccessor, error) {
	stat, err := os.Stat(giteaRootDir)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !stat.IsDir() {
		err = errors.New("Gitea root path " + giteaRootDir + " is not a directory")
		log.Error(err)
		return nil, err
	}

	giteaMainConfigPath := "/etc/gitea/conf/app.ini"
	giteaMainConfig, err := fetchConfig(giteaMainConfigPath)
	if err != nil {
		return nil, err
	}
	giteaCustomConfigPath := fmt.Sprintf("%s/custom/conf/app.ini", giteaRootDir)
	if err != nil {
		return nil, err
	}
	giteaCustomConfig, err := fetchConfig(giteaCustomConfigPath)
	if giteaMainConfig == nil && giteaCustomConfig == nil {
		err = errors.New("Cannot find Gitea config in  " + giteaMainConfigPath + " or " + giteaCustomConfigPath)
		log.Error(err)
		return nil, err
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
		log.Error(err)
		return nil, err
	}

	log.Infof("Using Gitea database %s\n", giteaDbPath)
	giteaAccessor.db = giteaDb

	giteaRepoID, err := giteaAccessor.getRepoID(giteaUserName, giteaRepoName)
	if err != nil {
		return nil, err
	}
	if giteaRepoID == -1 {
		err = errors.New("Cannot find repository " + giteaRepoName + " for user " + giteaUserName)
		log.Error(err)
		return nil, err
	}
	giteaAccessor.repoID = giteaRepoID

	// work out user ids
	adminUserID, err := giteaAccessor.getAdminUserID()
	if err != nil {
		return nil, err
	}
	giteaDefaultAssigneeID, err := giteaAccessor.getAdminDefaultingUserID(defaultAssignee, adminUserID)
	if err != nil {
		return nil, err
	}
	giteaAccessor.defaultAssigneeID = giteaDefaultAssigneeID

	giteaDefaultAuthorID, err := giteaAccessor.getAdminDefaultingUserID(defaultAuthor, adminUserID)
	if err != nil {
		return nil, err
	}
	giteaAccessor.defaultAuthorID = giteaDefaultAuthorID

	// find directory into which to clone wiki
	wikiRepoName := giteaUserName + "/" + giteaRepoName + ".wiki"
	if giteaWikiRepoDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Error(err)
			return nil, err
		}

		giteaWikiRepoDir = filepath.Join(cwd, wikiRepoName)
	}
	_, err = os.Stat(giteaWikiRepoDir)
	if !os.IsNotExist(err) {
		err = errors.New("wiki repository directory " + giteaWikiRepoDir + " already exists!")
		log.Error(err)
		return nil, err
	}
	giteaAccessor.wikiRepoDir = giteaWikiRepoDir

	// find URL from which clone wiki
	if giteaWikiRepoURL == "" {
		giteaWikiRepoURL = giteaAccessor.getUserRepoURL() + ".git"
	}
	giteaAccessor.wikiRepoURL = giteaWikiRepoURL

	return &giteaAccessor, nil
}

func (accessor *DefaultAccessor) getUserRepoURL() string {
	baseURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", baseURL, accessor.userName, accessor.repoName)
}
