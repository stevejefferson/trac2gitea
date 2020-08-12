package wiki

import (
	"fmt"
	"time"

	"stevejefferson.co.uk/trac2gitea/log"

	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/giteawiki"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

// Importer imports Trac Wiki data into a Gitea wiki repository.
type Importer struct {
	tracAccessor           *trac.Accessor
	giteaAccessor          *gitea.Accessor
	wikiAccessor           *giteawiki.Accessor
	trac2MarkdownConverter *markdown.Converter
	defaultPageOwner       string
	defaultPageOwnerEMail  string
}

// CreateImporter creates a Trac wiki to Gitea wiki repository importer.
func CreateImporter(
	tAccessor *trac.Accessor,
	gAccessor *gitea.Accessor,
	wAccessor *giteawiki.Accessor,
	t2mConverter *markdown.Converter,
	dfltPageOwner string) *Importer {

	dfltPageOwnerID := gAccessor.GetUserID(dfltPageOwner)
	if dfltPageOwnerID == -1 {
		log.Fatalf("Cannot find default owner %s for wiki pages to be imported from Trac\n", dfltPageOwner)
	}
	dfltPageOwnerEMail := gAccessor.GetUserEMailAddress(dfltPageOwnerID)

	importer := Importer{
		wikiAccessor:           wAccessor,
		tracAccessor:           tAccessor,
		giteaAccessor:          gAccessor,
		trac2MarkdownConverter: t2mConverter,
		defaultPageOwner:       dfltPageOwner,
		defaultPageOwnerEMail:  dfltPageOwnerEMail}
	return &importer
}

// ImportWiki imports a Trac wiki into a Gitea wiki repository.
func (importer *Importer) ImportWiki() {
	importer.wikiAccessor.RepoClone()

	importer.tracAccessor.GetWikiPages(func(pageName string, pageText string, author string, comment string, version int64, updateTime int64) {
		// convert and write wiki page
		markdownText := importer.trac2MarkdownConverter.Convert(pageText)
		translatedPageName := importer.wikiAccessor.TranslatePageName(pageName)
		importer.wikiAccessor.WritePage(translatedPageName, markdownText)

		// translate Trac wiki page (version) author into a Gitea user
		giteaAuthor := importer.defaultPageOwner
		giteaAuthorEMail := importer.defaultPageOwnerEMail
		giteaAuthorID := importer.giteaAccessor.GetUserID(author)
		if giteaAuthorID != -1 {
			giteaAuthor = author
			giteaAuthorEMail = importer.giteaAccessor.GetUserEMailAddress(giteaAuthorID)
		}

		// commit version of wiki page to local repository
		updateTimeStr := time.Unix(updateTime, 0)
		comment = fmt.Sprintf("%s\n[Imported from trac: page %s (version %d) updated at %s by Trac user %s]\n",
			comment, translatedPageName, version, updateTimeStr, author)
		importer.wikiAccessor.RepoStageAndCommit(giteaAuthor, giteaAuthorEMail, comment)
	})

	importer.wikiAccessor.RepoComplete()
}
