// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/stevejefferson/trac2gitea/importer"
	"github.com/stevejefferson/trac2gitea/markdown"

	"github.com/spf13/pflag"
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
)

var dbOnly bool
var wikiOnly bool
var wikiPush bool
var verbose bool
var wikiConvertPredefineds bool
var generateMaps bool
var tracRootDir string
var giteaRootDir string
var giteaUser string
var giteaRepo string
var userMapInputFile string
var userMapOutputFile string
var labelMapInputFile string
var labelMapOutputFile string
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

	generateMapsParam := pflag.Bool("generate-maps", false,
		"generate default user/label mappings into provided map files (note: no conversion will be performed in this case)")
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
			"Usage: %s [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo> [<user-map>] [<label-map>]\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	verbose = *verboseParam
	dbOnly = *dbOnlyParam
	wikiOnly = *wikiOnlyParam
	wikiPush = !*wikiNoPushParam
	generateMaps = *generateMapsParam

	if dbOnly && wikiOnly {
		log.Fatal("cannot generate only database AND only wiki!")
	}
	wikiConvertPredefineds = *wikiConvertPredefinedsParam
	giteaWikiRepoURL = *wikiURLParam
	giteaWikiRepoToken = *wikiTokenParam
	giteaWikiRepoDir = *wikiDirParam

	if (pflag.NArg() < 4) || (pflag.NArg() > 6) {
		pflag.Usage()
		os.Exit(1)
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUser = pflag.Arg(2)
	giteaRepo = pflag.Arg(3)
	if pflag.NArg() > 4 {
		userMapFile := pflag.Arg(4)
		if generateMaps {
			userMapOutputFile = userMapFile
		} else {
			userMapInputFile = userMapFile
		}
	}

	if pflag.NArg() > 5 {
		labelMapFile := pflag.Arg(5)
		if generateMaps {
			labelMapOutputFile = labelMapFile
		} else {
			labelMapInputFile = labelMapFile
		}
	}
}

func importData(dataImporter *importer.Importer, userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) error {
	var err error
	if err = dataImporter.ImportComponents(componentMap); err != nil {
		return err
	}
	if err = dataImporter.ImportPriorities(priorityMap); err != nil {
		return err
	}
	if err = dataImporter.ImportResolutions(resolutionMap); err != nil {
		return err
	}
	if err = dataImporter.ImportSeverities(severityMap); err != nil {
		return err
	}
	if err = dataImporter.ImportTypes(typeMap); err != nil {
		return err
	}
	if err = dataImporter.ImportVersions(versionMap); err != nil {
		return err
	}
	if err = dataImporter.ImportMilestones(); err != nil {
		return err
	}
	if err = dataImporter.ImportTickets(userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap); err != nil {
		return err
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
		log.Fatal("%+v", err)
	}
	giteaAccessor, err := gitea.CreateDefaultAccessor(
		giteaRootDir, giteaUser, giteaRepo, giteaWikiRepoURL, giteaWikiRepoToken, giteaWikiRepoDir)
	if err != nil {
		log.Fatal("%+v", err)
	}
	markdownConverter := markdown.CreateDefaultConverter(tracAccessor, giteaAccessor)

	dataImporter, err := importer.CreateImporter(tracAccessor, giteaAccessor, markdownConverter, giteaUser, wikiConvertPredefineds)
	if err != nil {
		log.Fatal("%+v", err)
	}

	userMap, err := readUserMap(userMapInputFile, dataImporter)
	if err != nil {
		log.Fatal("%+v", err)
	}

	componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap, err := readLabelMaps(labelMapInputFile, dataImporter)
	if err != nil {
		log.Fatal("%+v", err)
	}

	if generateMaps {
		if userMapOutputFile != "" {
			if err = writeUserMapToFile(userMapOutputFile, userMap); err != nil {
				log.Fatal("%+v", err)
			}
			log.Info("wrote user map to %s", userMapOutputFile)
		}
		if labelMapOutputFile != "" {
			if err = writeLabelMapsToFile(labelMapOutputFile, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap); err != nil {
				log.Fatal("%+v", err)
			}
			log.Info("wrote label map to %s", labelMapOutputFile)
		}
		return
	}

	if !wikiOnly {
		if err = importData(dataImporter, userMap, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap); err != nil {
			log.Fatal("%+v", err)
		}
	}

	if !dbOnly {
		if err = dataImporter.ImportWiki(userMap, wikiPush); err != nil {
			log.Fatal("%+v", err)
		}
	}
}
