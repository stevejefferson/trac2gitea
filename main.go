// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/import/issue"
	"github.com/stevejefferson/trac2gitea/import/wiki"
	"github.com/stevejefferson/trac2gitea/log"
)

var dbOnly bool
var wikiOnly bool
var wikiPush bool
var verbose bool
var wikiConvertPredefineds bool
var writeUserMap bool
var tracRootDir string
var giteaRootDir string
var giteaUser string
var giteaRepo string
var userMapFile string
var giteaWikiRepoURL string
var giteaWikiRepoToken string
var giteaWikiRepoDir string

func parseArgs() {
	wikiURLParam := pflag.String("wiki-url", "",
		"URL of wiki repository - defaults to <server-root-url>/<gitea-user>/<gitea-repo>.wiki.git")
	wikiTokenParam := pflag.String("wiki-token", "",
		"password/token for accessing wiki repository (ignored if wiki-url provided)")
	wikiDirParam := pflag.String("wiki-dir", "",
		"directory into which to checkout (clone) wiki repository - defaults to cwd")
	wikiConvertPredefinedsParam := pflag.Bool("wiki-convert-predefined", false,
		"convert Trac predefined wiki pages - by default we skip these")

	writeUserMapParam := pflag.Bool("write-user-map", false,
		"write default map of trac user to gitea user into the user map file (note: no conversion will be performed if this param is provided)")
	dbOnlyParam := pflag.Bool("db-only", false,
		"convert database only")
	wikiOnlyParam := pflag.Bool("wiki-only", false,
		"convert wiki only")
	wikiNoPushParam := pflag.Bool("no-wiki-push", false,
		"do not push wiki on completion")
	verboseParam := pflag.Bool("verbose", false,
		"verbose output")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo> [<user-map>]\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	verbose = *verboseParam
	dbOnly = *dbOnlyParam
	wikiOnly = *wikiOnlyParam
	wikiPush = !*wikiNoPushParam
	writeUserMap = *writeUserMapParam

	if dbOnly && wikiOnly {
		log.Fatal("cannot generate only database AND only wiki!")
	}
	wikiConvertPredefineds = *wikiConvertPredefinedsParam
	giteaWikiRepoURL = *wikiURLParam
	giteaWikiRepoToken = *wikiTokenParam
	giteaWikiRepoDir = *wikiDirParam

	if (pflag.NArg() < 4) || (pflag.NArg() > 5) {
		pflag.Usage()
		os.Exit(1)
	}
	if (pflag.NArg() == 4) && writeUserMap {
		log.Fatal("must provide user map file if writing user map")
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUser = pflag.Arg(2)
	giteaRepo = pflag.Arg(3)
	if pflag.NArg() == 5 {
		userMapFile = pflag.Arg(4)
	}
}

func readFromUserMap(mapFile string) (map[string]string, error) {
	fd, err := os.Open(mapFile)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	userMap := make(map[string]string)
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		userMapLine := scanner.Text()
		equalsPos := strings.LastIndex(userMapLine, "=")
		if equalsPos == -1 {
			return nil, fmt.Errorf("badly formatted user map file %s: found line %s", mapFile, userMapLine)
		}

		tracUserName := strings.Trim(userMapLine[0:equalsPos], " ")
		giteaUserName := strings.Trim(userMapLine[equalsPos+1:], " ")
		userMap[tracUserName] = giteaUserName
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return userMap, nil
}

func writeToUserMap(userMap map[string]string, mapFile string) error {
	fd, err := os.Create(mapFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	for tracUserName, giteaUserName := range userMap {
		if _, err := fd.WriteString(tracUserName + " = " + giteaUserName + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	parseArgs()

	var logLevel = log.INFO
	if verbose {
		logLevel = log.TRACE
	}
	log.SetLevel(logLevel)

	tracAccessor, err := trac.CreateDefaultAccessor(tracRootDir)
	if err != nil {
		log.Fatal("%v", err)
	}
	giteaAccessor, err := gitea.CreateDefaultAccessor(
		giteaRootDir, giteaUser, giteaRepo, giteaWikiRepoURL, giteaWikiRepoToken, giteaWikiRepoDir)
	if err != nil {
		log.Fatal("%v", err)
	}

	var userMap map[string]string
	if userMapFile != "" && !writeUserMap {
		userMap, err = readFromUserMap(userMapFile)
		if err != nil {
			log.Fatal("%v", err)
		}
	} else {
		userMap, err := tracAccessor.GetUserMap()
		if err != nil {
			log.Fatal("%v", err)
		}

		if err = giteaAccessor.GenerateDefaultUserMappings(userMap, giteaUser); err != nil {
			log.Fatal("%v", err)
		}
	}

	if writeUserMap {
		if err = writeToUserMap(userMap, userMapFile); err != nil {
			log.Fatal("%v", err)
		}

		log.Info("Trac to Gitea user mapping generated in %s - no conversion performed", userMapFile)
		return
	}

	if !wikiOnly {
		issueImporter, err := issue.CreateImporter(tracAccessor, giteaAccessor, userMap)
		if err != nil {
			log.Fatal("%v", err)
		}

		if err = issueImporter.ImportComponents(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportPriorities(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportSeverities(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportVersions(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportTypes(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportResolutions(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportMilestones(); err != nil {
			log.Fatal("%v", err)
		}
		if err = issueImporter.ImportTickets(); err != nil {
			log.Fatal("%v", err)
		}
	}

	if !dbOnly {
		wikiImporter, err := wiki.CreateImporter(tracAccessor, giteaAccessor, wikiConvertPredefineds, userMap)
		if err != nil {
			log.Fatal("%v", err)
		}

		if err = wikiImporter.ImportWiki(wikiPush); err != nil {
			log.Fatal("%v", err)
		}
	}
}
