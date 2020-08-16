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
	tracAccessor          trac.Accessor
	giteaAccessor         gitea.Accessor
	wikiAccessor          giteawiki.Accessor
	defaultPageOwner      string
	defaultPageOwnerEMail string
	convertPredefineds    bool
}

// CreateImporter creates a Trac wiki to Gitea wiki repository importer.
func CreateImporter(
	tAccessor trac.Accessor,
	gAccessor gitea.Accessor,
	wAccessor giteawiki.Accessor,
	dfltPageOwner string,
	convertPredefs bool) *Importer {

	dfltPageOwnerID := gAccessor.GetUserID(dfltPageOwner)
	if dfltPageOwnerID == -1 {
		log.Fatalf("Cannot find default owner %s for wiki pages to be imported from Trac\n", dfltPageOwner)
	}
	dfltPageOwnerEMail := gAccessor.GetUserEMailAddress(dfltPageOwnerID)

	importer := Importer{
		wikiAccessor:          wAccessor,
		tracAccessor:          tAccessor,
		giteaAccessor:         gAccessor,
		defaultPageOwner:      dfltPageOwner,
		defaultPageOwnerEMail: dfltPageOwnerEMail,
		convertPredefineds:    convertPredefs}
	return &importer
}

// ImportWiki imports a Trac wiki into a Gitea wiki repository.
func (importer *Importer) ImportWiki() {
	importer.wikiAccessor.RepoClone()

	importer.importWikiAttachments()
	importer.importWikiPages()

	importer.wikiAccessor.RepoComplete()
}

func (importer *Importer) importWikiAttachments() {
	importer.tracAccessor.GetWikiAttachments(func(pageName string, filename string) {
		tracAttachmentPath := importer.tracAccessor.GetWikiAttachmentPath(pageName, filename)
		giteaAttachmentPath := importer.wikiAccessor.GetAttachmentRelPath(pageName, filename)
		importer.wikiAccessor.CopyFile(tracAttachmentPath, giteaAttachmentPath)
	})
}

func (importer *Importer) importWikiPages() {
	importer.tracAccessor.GetWikiPages(func(pageName string, pageText string, author string, comment string, version int64, updateTime int64) {
		// skip predefined pages
		if !importer.convertPredefineds && importer.tracAccessor.IsPredefinedPage(pageName) {
			log.Debugf("Skipping predefined Trac page %s\n", pageName)
			return
		}

		tracToMarkdownConverter := markdown.CreateWikiDefaultConverter(
			importer.tracAccessor, importer.giteaAccessor, importer.wikiAccessor, pageName)

		// convert and write wiki page
		markdownText := tracToMarkdownConverter.Convert(pageText)
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
		log.Infof("Wiki page %s: wrote version %d to repository\n", translatedPageName, version)
	})

	importer.wikiAccessor.RepoComplete()
}
