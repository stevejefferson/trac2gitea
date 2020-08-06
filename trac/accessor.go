package trac

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3"
)

type Accessor struct {
	rootDir string
	db      *sql.DB
	config  *ini.File
}

func FindTrac(tracRootDir string) *Accessor {
	stat, err := os.Stat(tracRootDir)
	if err != nil {
		log.Fatal(err)
	}
	if stat.IsDir() != true {
		log.Fatal("Trac root directory is not a directory: ", tracRootDir)
	}

	tracIniPath := fmt.Sprintf("%s/conf/trac.ini", tracRootDir)
	stat, err = os.Stat(tracIniPath)
	if err != nil {
		log.Fatal(err)
	}

	tracConfig, err := ini.Load(tracIniPath)
	if err != nil {
		log.Fatal(err)
	}

	accessor := Accessor{db: nil, rootDir: tracRootDir, config: tracConfig}

	// extract path to trac DB - currently sqlite-specific...
	tracDatabaseString := accessor.GetStringConfig("trac", "database")
	tracDatabaseSegments := strings.SplitN(tracDatabaseString, ":", 2)
	tracDatabasePath := tracDatabaseSegments[1]
	if !filepath.IsAbs(tracDatabasePath) {
		tracDatabasePath = filepath.Join(tracRootDir, tracDatabasePath)
	}

	tracDb, err := sql.Open("sqlite3", tracDatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	accessor.db = tracDb

	return &accessor
}
