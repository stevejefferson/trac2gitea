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
var generatedUserMapFile string
var userMapFile string
var tracRootDir string
var giteaRootDir string
var giteaUser string
var giteaRepo string
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

	generateUserMapParam := pflag.String("generate-user-map", "",
		"generate a or update map of trac users to gitea users into the specified file (note: no conversion will be performed if this param is provided)")
	userMapParam := pflag.String("user-map", "",
		"user provided mapping of trac users to gitea users for conversion")
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
			"Usage: %s [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo>\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	verbose = *verboseParam
	dbOnly = *dbOnlyParam
	wikiOnly = *wikiOnlyParam
	wikiPush = !*wikiNoPushParam

	generatedUserMapFile = *generateUserMapParam
	userMapFile = *userMapParam

	if dbOnly && wikiOnly {
		log.Fatal("Cannot generate only database AND only wiki!")
	}
	wikiConvertPredefineds = *wikiConvertPredefinedsParam

	if pflag.NArg() < 4 {
		pflag.Usage()
		os.Exit(1)
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUser = pflag.Arg(2)
	giteaRepo = pflag.Arg(3)

	giteaWikiRepoURL = *wikiURLParam
	giteaWikiRepoToken = *wikiTokenParam
	giteaWikiRepoDir = *wikiDirParam
}

func readUserMap(mapFile string) (map[string]string, error) {
	fd, err := os.Open(mapFile)
	if err != nil {
		log.Error("Cannot create user map file %s: %v\n", mapFile, err)
		return nil, err
	}
	defer fd.Close()

	userMap := make(map[string]string)
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		userMapLine := scanner.Text()
		equalsPos := strings.LastIndex(userMapLine, "=")
		if equalsPos == -1 {
			err = fmt.Errorf("badly formatted user map file %s: found line %s", mapFile, userMapLine)
			log.Error("%v\n", err)
			return nil, err
		}

		tracUserName := userMapLine[0:equalsPos]
		giteaUserName := userMapLine[equalsPos+1:]
		userMap[tracUserName] = giteaUserName
	}

	if err := scanner.Err(); err != nil {
		log.Error("Problem reading user map file %s: %v\n", mapFile, err)
		return nil, err
	}

	return userMap, nil
}

func writeUserMap(userMap map[string]string, mapFile string) error {
	fd, err := os.Create(mapFile)
	if err != nil {
		log.Error("Cannot create user map file %s: %v\n", mapFile, err)
		return err
	}
	defer fd.Close()

	for tracUserName, giteaUserName := range userMap {
		_, err := fd.WriteString(tracUserName + "=" + giteaUserName + "\n")
		if err != nil {
			log.Error("Cannot write user mapping %s=%s to map file %s: %v\n", tracUserName, giteaUserName, mapFile, err)
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
		log.Fatal("%v\n", err)
	}
	giteaAccessor, err := gitea.CreateDefaultAccessor(
		giteaRootDir, giteaUser, giteaRepo, giteaWikiRepoURL, giteaWikiRepoToken, giteaWikiRepoDir)
	if err != nil {
		log.Fatal("%v\n", err)
	}

	var userMap map[string]string
	if userMapFile != "" {
		userMap, err = readUserMap(userMapFile)
	} else {
		userMap, err = tracAccessor.GetUserMap()
		giteaAccessor.GenerateDefaultUserMappings(userMap, giteaUser)
	}
	if err != nil {
		log.Fatal("%v\n", err)
	}

	if generatedUserMapFile != "" {
		writeUserMap(userMap, generatedUserMapFile)
		log.Info("Trac -> Gitea user mapping generated in %s\n", generatedUserMapFile)
		return
	}

	if !wikiOnly {
		issueImporter, err := issue.CreateImporter(tracAccessor, giteaAccessor, userMap)
		if err != nil {
			log.Fatal("%v\n", err)
		}

		issueImporter.ImportComponents()
		issueImporter.ImportPriorities()
		issueImporter.ImportSeverities()
		issueImporter.ImportVersions()
		issueImporter.ImportTypes()
		issueImporter.ImportResolutions()
		issueImporter.ImportMilestones()
		issueImporter.ImportTickets()
	}

	if !dbOnly {
		wikiImporter, err := wiki.CreateImporter(tracAccessor, giteaAccessor, wikiConvertPredefineds, userMap)
		if err != nil {
			log.Fatal("%v\n", err)
		}

		wikiImporter.ImportWiki(wikiPush)
	}
}
