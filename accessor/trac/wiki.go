package trac

import (
	"database/sql"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetWikiPages retrieves all Trac wiki pages, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetWikiPages(handlerFn func(pageName string, pageText string, author string, comment string, version int64, updateTime int64)) {
	rows, err := accessor.db.Query(`SELECT name, text, author, comment, version, CAST(time*1e-6 AS int8) FROM wiki`)
	// SELECT w1.name, w1.text, w1.author, w1.comment, w1.version, CAST(w1.time*1e-6 AS int8)
	// 	FROM wiki w1
	// 	WHERE w1.version = (SELECT MAX(w2.version) FROM wiki w2 WHERE w1.name = w2.name)`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var pageName string
		var pageText string
		var author string
		var commentStr sql.NullString
		var version int64
		var updateTime int64
		if err := rows.Scan(&pageName, &pageText, &author, &commentStr, &version, &updateTime); err != nil {
			log.Fatal(err)
		}

		comment := ""
		if !commentStr.Valid {
			comment = commentStr.String
		}

		handlerFn(pageName, pageText, author, comment, version, updateTime)
	}
}

var prefinedTracPages = []string{
	"CamelCase",
	"InterMapTxt",
	"InterTrac",
	"InterWiki",
	"PageTemplates",
	"RecentChanges",
	"SandBox",
	"TicketQuery",
	"TitleIndex",
	"TracAccessibility",
	"TracAdmin",
	"TracBackup",
	"TracBatchModify",
	"TracBrowser",
	"TracCgi",
	"TracChangeLog",
	"TracChangeset",
	"TracEnvironment",
	"TracFastCgi",
	"TracFineGrainedPermissions",
	"TracGuide",
	"TracImport",
	"TracIni",
	"TracInstall",
	"TracInterfaceCustomization",
	"TracLinks",
	"TracLogging",
	"TracModPython",
	"TracModWSGI",
	"TracNavigation",
	"TracNotification",
	"TracPermissions",
	"TracPlugins",
	"TracQuery",
	"TracReports",
	"TracRepositoryAdmin",
	"TracRevisionLog",
	"TracRoadmap",
	"TracRss",
	"TracSearch",
	"TracStandalone", "TracUnicode",
	"TracSupport",
	"TracSyntaxColoring",
	"TracTickets",
	"TracTicketsCustomFields",
	"TracTimeline",
	"TracUpgrade",
	"TracWiki",
	"TracWorkflow",
	"WikiDeletePage",
	"WikiFormatting",
	"WikiHtml",
	"WikiMacros",
	"WikiNewPage",
	"WikiPageNames",
	"WikiProcessors",
	"WikiRestructuredText",
	"WikiRestructuredTextLinks",
	//"WikiStart"	// keep WikiStart - the default contents are usually overwritten by projects to produce a "home" page
}

// IsPredefinedPage returns true if the provided page name is one of Trac's predefined ones - by default we ignore these
func (accessor *DefaultAccessor) IsPredefinedPage(pageName string) bool {
	for _, predefinedTracPage := range prefinedTracPages {
		if pageName == predefinedTracPage {
			return true
		}
	}

	return false
}
