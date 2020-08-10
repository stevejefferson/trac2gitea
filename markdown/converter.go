package markdown

import (
	"stevejefferson.co.uk/trac2gitea/gitea"
	"stevejefferson.co.uk/trac2gitea/trac"
	"stevejefferson.co.uk/trac2gitea/wiki"
)

// Converter of Trac markdown to Gitea markdown
// This is used in two circumstances:
// 1. for ticket comments - in which case ticketID != -1 and wikiAccessor == nil
// 2. for wiki imports - in which case ticketID == -1 and wikiAccessor != nil
type Converter struct {
	tracAccessor  *trac.Accessor
	giteaAccessor *gitea.Accessor
	wikiAccessor  *wiki.Accessor
	ticketID      int64
}

// CreateWikiConverter returns a Trac to Gitea markdown converter for converting wiki texts.
func CreateWikiConverter(tAccessor *trac.Accessor, gAccessor *gitea.Accessor, wAccessor *wiki.Accessor) *Converter {
	converter := Converter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiAccessor: wAccessor, ticketID: -1}
	return &converter
}

// CreateTicketConverter returns a Trac to Gitea markdown converter for converting trac ticket descriptions and comments.
func CreateTicketConverter(tAccessor *trac.Accessor, gAccessor *gitea.Accessor, tracTicketID int64) *Converter {
	converter := Converter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiAccessor: nil, ticketID: tracTicketID}
	return &converter
}

func (converter *Converter) convertNonCodeBlockText(in string) string {
	out := in
	out = converter.convertLinks(out)
	out = converter.convertAnchors(out)
	out = converter.convertEscapes(out)
	out = converter.convertParagraphs(out)
	out = converter.convertLists(out)
	out = converter.convertDefinitionLists(out)
	out = converter.disguiseLinks(out)
	out = converter.convertBlockQuotes(out)
	out = converter.convertHeadings(out)
	out = converter.convertFontStyles(out)
	out = converter.convertTables(out)
	out = converter.undisguiseLinks(out)
	return out
}

// Convert converts a string of Trac markdown to Gitea markdown
func (converter *Converter) Convert(in string) string {
	out := in

	// ensure we have Unix EOLs
	out = converter.convertEOL(out)

	// convert any code blocks
	out = converter.convertCodeBlocks(out)

	// perform all other conversions only on text not in a code block
	out = converter.convertNonCodeBlocks(out, converter.convertNonCodeBlockText)

	return out
}
