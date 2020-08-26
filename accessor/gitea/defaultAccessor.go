// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stevejefferson/trac2gitea/log"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3" // sqlite database driver
	"gopkg.in/src-d/go-git.v4"
)

// DefaultAccessor is the default implementation of the gitea Accessor interface, accessing Gitea directly via its database and filestore.
type DefaultAccessor struct {
	rootDir       string
	mainConfig    *ini.File
	customConfig  *ini.File
	db            *sql.DB
	userName      string
	repoName      string
	repoID        int64
	wikiRepoURL   string
	wikiRepoToken string
	wikiRepoDir   string
	wikiRepo      *git.Repository
}

func fetchConfig(configPath string) (*ini.File, error) {
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, nil
	}

	config, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load config %s: %v", configPath, err)
	}

	return config, nil
}

// CreateDefaultAccessor returns a new Gitea default accessor.
func CreateDefaultAccessor(
	giteaRootDir string,
	giteaUserName string,
	giteaRepoName string,
	giteaWikiRepoURL string,
	giteaWikiRepoToken string,
	giteaWikiRepoDir string) (*DefaultAccessor, error) {
	stat, err := os.Stat(giteaRootDir)
	if err != nil {
		return nil, fmt.Errorf("cannot access Gitea root directory %s: %v", giteaRootDir, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("gitea root path %s is not a directory", giteaRootDir)
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
		return nil, fmt.Errorf("cannot find Gitea config in %s or %s", giteaMainConfigPath, giteaCustomConfigPath)
	}

	giteaAccessor := DefaultAccessor{
		rootDir:       giteaRootDir,
		mainConfig:    giteaMainConfig,
		customConfig:  giteaCustomConfig,
		db:            nil,
		userName:      giteaUserName,
		repoName:      giteaRepoName,
		repoID:        0,
		wikiRepoURL:   "",
		wikiRepoToken: "",
		wikiRepoDir:   "",
		wikiRepo:      nil}

	// extract path to gitea DB - currently sqlite-specific...
	giteaDbPath := giteaAccessor.GetStringConfig("database", "PATH")
	giteaDb, err := sql.Open("sqlite3", giteaDbPath)
	if err != nil {
		return nil, err
	}

	log.Info("using Gitea database %s", giteaDbPath)
	giteaAccessor.db = giteaDb

	giteaRepoID, err := giteaAccessor.getRepoID(giteaUserName, giteaRepoName)
	if err != nil {
		return nil, err
	}
	if giteaRepoID == -1 {
		return nil, fmt.Errorf("cannot find repository %s for user %s", giteaRepoName, giteaUserName)
	}
	giteaAccessor.repoID = giteaRepoID

	// find directory into which to clone wiki
	wikiRepoName := giteaRepoName + ".wiki"
	if giteaWikiRepoDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		giteaWikiRepoDir = filepath.Join(cwd, wikiRepoName)
	}
	_, err = os.Stat(giteaWikiRepoDir)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("wiki repository directory %s already exists", giteaWikiRepoDir)
	}
	giteaAccessor.wikiRepoDir = giteaWikiRepoDir

	// find URL from which clone wiki
	if giteaWikiRepoURL == "" {
		rootURL := giteaAccessor.GetStringConfig("server", "ROOT_URL")
		if giteaWikiRepoToken != "" {
			slashSlashPos := strings.Index(rootURL, "//")
			if slashSlashPos == -1 {
				return nil, fmt.Errorf("ROOT_URL %s malformed? expecting a '//'", rootURL)
			}

			// insert username and token into URL - 'http://example.com' should become 'http://<user>:<token>@example.com'
			rootURL = rootURL[0:slashSlashPos+2] + giteaUserName + ":" + giteaWikiRepoToken + "@" + rootURL[slashSlashPos+2:]

			giteaAccessor.wikiRepoToken = giteaWikiRepoToken
		}
		if rootURL[len(rootURL)-1:] != "/" {
			rootURL = rootURL + "/"
		}
		giteaWikiRepoURL = fmt.Sprintf("%s%s/%s.git", rootURL, giteaUserName, wikiRepoName)
	}
	log.Info("using Wiki repo URL %s", giteaWikiRepoURL)
	giteaAccessor.wikiRepoURL = giteaWikiRepoURL

	return &giteaAccessor, nil
}
