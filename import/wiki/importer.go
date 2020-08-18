package wiki

import (
	"fmt"
	"time"

	"stevejefferson.co.uk/trac2gitea/log"

	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
	"stevejefferson.co.uk/trac2gitea/markdown"
)

// Importer imports Trac Wiki data into a Gitea wiki repository.
type Importer struct {
	tracAccessor          trac.Accessor
	giteaAccessor         gitea.Accessor
	defaultPageOwner      string
	defaultPageOwnerEMail string
	convertPredefineds    bool
}

// CreateImporter creates a Trac wiki to Gitea wiki repository importer.
func CreateImporter(
	tAccessor trac.Accessor,
	gAccessor gitea.Accessor,
	dfltPageOwner string,
	convertPredefs bool) *Importer {

	dfltPageOwnerID := gAccessor.GetUserID(dfltPageOwner)
	if dfltPageOwnerID == -1 {
		log.Fatalf("Cannot find default owner %s for wiki pages to be imported from Trac\n", dfltPageOwner)
	}
	dfltPageOwnerEMail := gAccessor.GetUserEMailAddress(dfltPageOwnerID)

	importer := Importer{
		tracAccessor:          tAccessor,
		giteaAccessor:         gAccessor,
		defaultPageOwner:      dfltPageOwner,
		defaultPageOwnerEMail: dfltPageOwnerEMail,
		convertPredefineds:    convertPredefs}
	return &importer
}

// ImportWiki imports a Trac wiki into a Gitea wiki repository.
func (importer *Importer) ImportWiki() {
	importer.giteaAccessor.CloneWiki()

	importer.importWikiAttachments()
	importer.importWikiPages()

	importer.giteaAccessor.PushWiki()
}

func (importer *Importer) importWikiAttachments() {
	importer.tracAccessor.GetWikiAttachments(func(pageName string, filename string) {
		tracAttachmentPath := importer.tracAccessor.GetWikiAttachmentPath(pageName, filename)
		giteaAttachmentPath := importer.giteaAccessor.GetWikiAttachmentRelPath(pageName, filename)
		importer.giteaAccessor.CopyFileToWiki(tracAttachmentPath, giteaAttachmentPath)
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
			importer.tracAccessor, importer.giteaAccessor, pageName)

		// convert and write wiki page
		markdownText := tracToMarkdownConverter.Convert(pageText)
		translatedPageName := importer.giteaAccessor.TranslateWikiPageName(pageName)
		importer.giteaAccessor.WriteWikiPage(translatedPageName, markdownText)

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
		importer.giteaAccessor.CommitWiki(giteaAuthor, giteaAuthorEMail, comment)
		log.Infof("Wiki page %s: converted trac version %d\n", translatedPageName, version)
	})
}
