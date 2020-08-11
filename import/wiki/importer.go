package wiki

import (
	"fmt"
	"time"

	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/giteaWiki"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

// Importer imports Trac Wiki data into a Gitea wiki repository.
type Importer struct {
	tracAccessor           *trac.Accessor
	giteaAccessor          *gitea.Accessor
	wikiAccessor           *giteaWiki.Accessor
	trac2MarkdownConverter *markdown.Converter
}

// CreateImporter creates a Trac wiki to Gitea wiki repository importer.
func CreateImporter(
	tAccessor *trac.Accessor,
	gAccessor *gitea.Accessor,
	wAccessor *giteaWiki.Accessor,
	t2mConverter *markdown.Converter) *Importer {
	importer := Importer{wikiAccessor: wAccessor, tracAccessor: tAccessor, giteaAccessor: gAccessor, trac2MarkdownConverter: t2mConverter}
	return &importer
}

// ImportWiki imports a Trac wiki into a Gitea wiki repository.
func (importer *Importer) ImportWiki() {
	repoOwnerEmailAddress := importer.giteaAccessor.GetEMailAddress()
	importer.wikiAccessor.RepoClone()

	importer.tracAccessor.GetWikiPages(func(pageName string, pageText string, author string, comment string, version int64, updateTime int64) {
		markdownText := importer.trac2MarkdownConverter.Convert(pageText)
		translatedPageName := importer.wikiAccessor.TranslatePageName(pageName)
		importer.wikiAccessor.WritePage(translatedPageName, markdownText)

		updateTimeStr := time.Unix(updateTime, 0)
		comment = fmt.Sprintf("%s\n[Imported from trac: page %s (version %d) updated at %s]\n", comment, translatedPageName, version, updateTimeStr)
		importer.wikiAccessor.RepoStageAndCommit(author, repoOwnerEmailAddress, comment)
	})

	importer.wikiAccessor.RepoComplete()
}
