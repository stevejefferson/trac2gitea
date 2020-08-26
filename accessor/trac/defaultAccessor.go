// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/stevejefferson/trac2gitea/log"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3" // sqlite database driver
)

// DefaultAccessor is the default implementation of the trac Accessor interface, accessing Trac via its database and filestore.
type DefaultAccessor struct {
	rootDir string
	db      *sql.DB
	config  *ini.File
}

// CreateDefaultAccessor creates a new Trac accessor.
func CreateDefaultAccessor(tracRootDir string) (*DefaultAccessor, error) {
	stat, err := os.Stat(tracRootDir)
	if err != nil {
		err = errors.Wrapf(err, "looking for root directory %s of Trac instance", tracRootDir)
		return nil, err
	}
	if stat.IsDir() != true {
		return nil, fmt.Errorf("trac root %s is not a directory", tracRootDir)
	}

	tracIniPath := fmt.Sprintf("%s/conf/trac.ini", tracRootDir)
	stat, err = os.Stat(tracIniPath)
	if err != nil {
		err = errors.Wrapf(err, "looking for Trac config file %s", tracIniPath)
		return nil, err
	}

	tracConfig, err := ini.Load(tracIniPath)
	if err != nil {
		err = errors.Wrapf(err, "loading Trac config file %s", tracIniPath)
		return nil, err
	}

	accessor := DefaultAccessor{db: nil, rootDir: tracRootDir, config: tracConfig}

	// extract path to trac DB - currently sqlite-specific...
	tracDatabaseString := accessor.GetStringConfig("trac", "database")
	tracDatabaseSegments := strings.SplitN(tracDatabaseString, ":", 2)
	tracDatabasePath := tracDatabaseSegments[1]
	if !filepath.IsAbs(tracDatabasePath) {
		tracDatabasePath = filepath.Join(tracRootDir, tracDatabasePath)
	}

	log.Info("using trac database %s", tracDatabasePath)

	tracDb, err := sql.Open("sqlite3", tracDatabasePath)
	if err != nil {
		err = errors.Wrapf(err, "opening Trac sqlite database %s", tracDatabasePath)
		return nil, err
	}
	accessor.db = tracDb

	return &accessor, nil
}
