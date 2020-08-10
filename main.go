package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"stevejefferson.co.uk/trac2gitea/gitea"
	"stevejefferson.co.uk/trac2gitea/issueimport"
	"stevejefferson.co.uk/trac2gitea/markdown"
	"stevejefferson.co.uk/trac2gitea/trac"
	"stevejefferson.co.uk/trac2gitea/wiki"
	"stevejefferson.co.uk/trac2gitea/wikiimport"
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
var giteaUserName string
var giteaRepoName string
var giteaWikiRepoDir string
var defaultAssignee string
var defaultAuthor string

func parseArgs() {
	defaultAssigneeParam := pflag.String("default-assignee", "",
		"`username` to assign tickets to when trac assignee is not found in Gitea")
	defaultAuthorParam := pflag.String("default-author", "",
		"`username` to attribute content to when trac author is not found in Gitea")
	dbOnlyParam := pflag.Bool("db-only", false,
		"convert database only")
	wikiOnlyParam := pflag.Bool("wiki-only", false,
		"convert wiki only")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [options] <trac_root> <gitea_root> <gitea_user> <gitea_repo_name> <gitea_wiki_repo_dir>\n",
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

	if pflag.NArg() < 5 {
		pflag.Usage()
		os.Exit(1)
	}

	tracRootDir = pflag.Arg(0)
	giteaRootDir = pflag.Arg(1)
	giteaUserName = pflag.Arg(2)
	giteaRepoName = pflag.Arg(3)
	giteaWikiRepoDir = pflag.Arg(4)

	defaultAssignee = *defaultAssigneeParam
	if defaultAssignee == "" {
		defaultAssignee = giteaUserName
	}
	defaultAuthor = *defaultAuthorParam
	if defaultAuthor == "" {
		defaultAuthor = giteaUserName
	}

}

func main() {
	setLogFormat()

	parseArgs()

	// low-level accessors
	tracAccessor := trac.CreateAccessor(tracRootDir)
	giteaAccessor := gitea.CreateAccessor(giteaRootDir, giteaUserName, giteaRepoName, defaultAssignee, defaultAuthor)
	wikiAccessor := wiki.CreateAccessor(giteaWikiRepoDir)

	// data converters
	trac2MarkdownConverter := markdown.CreateConverter(tracAccessor, giteaAccessor, wikiAccessor)

	// importers
	if !wikiOnly {
		issueImporter := issueimport.CreateImporter(tracAccessor, giteaAccessor, trac2MarkdownConverter)

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
		wikiImporter := wikiimport.CreateImporter(tracAccessor, giteaAccessor, wikiAccessor, trac2MarkdownConverter)
		wikiImporter.ImportWiki()
	}
}
