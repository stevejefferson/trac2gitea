package wikiimport

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

// ImportWiki imports a Trac wiki into a Gitea wiki repository, optionally pushing the resultant updates back to the remote.
func (importer *Importer) ImportWiki(pushRepo bool) {
	repoOwnerEmailAddress := importer.giteaAccessor.GetEMailAddress()
	importer.wikiAccessor.RepoClone()

	rows := importer.tracAccessor.Query(`
		SELECT w1.name, w1.text, w1.author, w1.comment, w1.version, CAST(w1.time*1e-6 AS int8)
			FROM wiki w1
			WHERE w1.version = (SELECT MAX(w2.version) FROM wiki w2 WHERE w1.name = w2.name)`)

	for rows.Next() {
		var name string
		var text string
		var author string
		var commentStr sql.NullString
		var version int64
		var updateTime int64
		if err := rows.Scan(&name, &text, &author, &commentStr, &version, &updateTime); err != nil {
			log.Fatal(err)
		}

		markdownText := importer.trac2MarkdownConverter.Convert(text)
		importer.wikiAccessor.WritePage(name, markdownText)

		comment := ""
		if !commentStr.Valid {
			comment = commentStr.String
		}

		updateTimeStr := time.Unix(updateTime, 0)
		comment = fmt.Sprintf("%s\n[Imported from trac: original page (version %d) updated at %s]\n", comment, version, updateTimeStr)
		importer.wikiAccessor.RepoStageAndCommit(author, repoOwnerEmailAddress, comment)
	}

	if pushRepo {
		importer.wikiAccessor.RepoPush()
	}
}
