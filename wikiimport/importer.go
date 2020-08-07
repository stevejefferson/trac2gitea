package wikiimport

import (
	"database/sql"
	"fmt"
	"log"

	"stevejefferson.co.uk/trac2gitea/gitea"
	"stevejefferson.co.uk/trac2gitea/markdown"
	"stevejefferson.co.uk/trac2gitea/trac"
	"stevejefferson.co.uk/trac2gitea/wiki"
)

// Importer imports Trac Wiki data into a Gitea wiki repository.
type Importer struct {
	tracAccessor           *trac.Accessor
	giteaAccessor          *gitea.Accessor
	wikiAccessor           *wiki.Accessor
	trac2MarkdownConverter *markdown.Converter
}

// CreateImporter creates a Trac wiki to Gitea wiki repository importer.
func CreateImporter(
	tAccessor *trac.Accessor,
	gAccessor *gitea.Accessor,
	wAccessor *wiki.Accessor,
	t2mConverter *markdown.Converter) *Importer {
	importer := Importer{wikiAccessor: wAccessor, tracAccessor: tAccessor, giteaAccessor: gAccessor, trac2MarkdownConverter: t2mConverter}
	return &importer
}

// ImportWiki imports a Trac wiki into a Gitea wiki repository.
func (importer *Importer) ImportWiki() {
	rows := importer.tracAccessor.Query(`
		SELECT w1.name, w1.text, w1.comment, w1.version, w1.time
			FROM wiki w1
			WHERE w1.version = (SELECT MAX(w2.version) FROM wiki w2 WHERE w1.name = w2.name)`)

	for rows.Next() {
		var name string
		var text string
		var commentStr sql.NullString
		var version int64
		var time int64
		if err := rows.Scan(&name, &text, &commentStr, &version, &time); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Converting Wiki page %s, version %d\n", name, version)
		comment := ""
		if !commentStr.Valid {
			comment = commentStr.String
		}
		markdownText := importer.trac2MarkdownConverter.Convert(text, "")
		importer.wikiAccessor.WritePageVersion(name, markdownText, version, comment, time)
	}
}
