package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/giteawiki"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
	"stevejefferson.co.uk/trac2gitea/import/issue"
	"stevejefferson.co.uk/trac2gitea/import/wiki"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

func setLogFormat() {
	log.SetFlags(log.Ldate)
	log.SetFlags(log.Ltime)
	log.SetFlags(log.Lshortfile)
}

var dbOnly bool
var wikiOnly bool
var tracRootDir string
var giteaRootDir string
var giteaUser string
var giteaRepo string
var giteaWikiRepoURL string
var giteaWikiRepoDir string
var giteaDefaultAssignee string
var giteaDefaultAuthor string
var giteaDefaultWikiAuthor string

func parseArgs() {
	defaultAssigneeParam := pflag.String("default-assignee", "",
		"`username` to assign tickets to when trac assignee is not found in Gitea - defaults to <gitea-user>")
	defaultAuthorParam := pflag.String("default-author", "",
		"`username` to attribute content to when trac author is not found in Gitea - defaults to <gitea-user>")
	defaultWikiAuthorParam := pflag.String("default-wiki-author", "",
		"`username` to attribute Wiki content to when trac author is not found in Gitea - defaults to <gitea-user>")

	wikiURLParam := pflag.String("wiki-url", "",
		"URL of wiki repository - defaults to <server-root-url>/<gitea-user>/<gitea-repo>.wiki.git")
	wikiDirParam := pflag.String("wiki-dir", "",
		"directory into which to checkout (clone) wiki repository - defaults to cwd")

	dbOnlyParam := pflag.Bool("db-only", false,
		"convert database only")
	wikiOnlyParam := pflag.Bool("wiki-only", false,
		"convert wiki only")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <trac-root> <gitea-root> <gitea-user> <gitea-repo>\n",
			os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	dbOnly = *dbOnlyParam
	wikiOnly = *wikiOnlyParam
	if dbOnly && wikiOnly {
		log.Fatal("Cannot generate only database AND only wiki!")
	}

	if pflag.NArg() < 4 {
		pflag.Usage()
		os.Exit(1)
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUser = pflag.Arg(2)
	giteaRepo = pflag.Arg(3)

	giteaDefaultAssignee = *defaultAssigneeParam
	if giteaDefaultAssignee == "" {
		giteaDefaultAssignee = giteaUser
	}
	giteaDefaultAuthor = *defaultAuthorParam
	if giteaDefaultAuthor == "" {
		giteaDefaultAuthor = giteaUser
	}
	giteaDefaultWikiAuthor = *defaultWikiAuthorParam
	if giteaDefaultWikiAuthor == "" {
		giteaDefaultWikiAuthor = giteaUser
	}
	giteaWikiRepoURL = *wikiURLParam
	giteaWikiRepoDir = *wikiDirParam
}

func main() {
	setLogFormat()

	parseArgs()

	tracAccessor := trac.CreateAccessor(tracRootDir)
	giteaAccessor := gitea.CreateAccessor(giteaRootDir, giteaUser, giteaRepo, giteaDefaultAssignee, giteaDefaultAuthor)

	if !wikiOnly {
		issueImporter := issue.CreateImporter(tracAccessor, giteaAccessor)

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
		if giteaWikiRepoURL == "" {
			giteaWikiRepoURL = giteaAccessor.GetWikiRepoURL()
		}

		if giteaWikiRepoDir == "" {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			wikiRepoName := giteaAccessor.GetWikiRepoName()
			giteaWikiRepoDir = filepath.Join(cwd, wikiRepoName)
		}

		wikiAccessor := giteawiki.CreateAccessor(giteaWikiRepoURL, giteaWikiRepoDir)
		wikiMarkdownConverter := markdown.CreateWikiConverter(tracAccessor, giteaAccessor, wikiAccessor)
		wikiImporter := wiki.CreateImporter(tracAccessor, giteaAccessor, wikiAccessor, wikiMarkdownConverter, giteaDefaultWikiAuthor)
		wikiImporter.ImportWiki()
	}
}
