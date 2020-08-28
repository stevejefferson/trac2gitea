// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package trac

import (
	"database/sql"

	"github.com/pkg/errors"
)

// GetWikiPages retrieves all Trac wiki pages, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetWikiPages(handlerFn func(page *WikiPage) error) error {
	rows, err := accessor.db.Query(`SELECT name, text, author, comment, version, CAST(time*1e-6 AS int8) FROM wiki`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving Trac wiki pages")
		return err
	}

	for rows.Next() {
		var pageName string
		var pageText string
		var author string
		var commentStr sql.NullString
		var version int64
		var updateTime int64
		if err := rows.Scan(&pageName, &pageText, &author, &commentStr, &version, &updateTime); err != nil {
			err = errors.Wrapf(err, "retrieving Trac wiki page")
			return err
		}

		comment := ""
		if !commentStr.Valid {
			comment = commentStr.String
		}

		wikiPage := WikiPage{Name: pageName, Text: pageText, Author: author, Comment: comment, Version: version, UpdateTime: updateTime}

		if err = handlerFn(&wikiPage); err != nil {
			return err
		}
	}

	return nil
}

// GetWikiAttachmentPath retrieves the path to a named attachment to a Trac wiki page.
func (accessor *DefaultAccessor) GetWikiAttachmentPath(attachment *WikiAttachment) string {
	return accessor.getAttachmentPath(attachment.PageName, attachment.FileName, "wiki")
}

// GetWikiAttachments retrieves all Trac wiki page attachments, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetWikiAttachments(handlerFn func(attachment *WikiAttachment) error) error {
	rows, err := accessor.db.Query(`SELECT id, filename FROM attachment WHERE type = 'wiki'`)
	if err != nil {
		err = errors.Wrapf(err, "retrieving attachments to Trac wiki pages")
		return err
	}

	for rows.Next() {
		var pageName string
		var filename string
		if err := rows.Scan(&pageName, &filename); err != nil {
			err = errors.Wrapf(err, "retrieving attachment to Trac wiki page")
			return err
		}

		attachment := WikiAttachment{PageName: pageName, FileName: filename}

		if err = handlerFn(&attachment); err != nil {
			return err
		}
	}

	return nil
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
	//"WikiStart"	// keep WikiStart - the default contents are usually overwritten by projects to produce a "Home" page
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
