package markdown

import (
	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/giteawiki"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
)

// DefaultConverter is the default implementation of the Trac markdown to Gitea markdown converter.
// This is used in two circumstances:
// 1. for ticket comments - in which case ticketID != -1 and wikiAccessor == nil
// 2. for wiki imports - in which case ticketID == -1 and wikiAccessor != nil
type DefaultConverter struct {
	tracAccessor  trac.Accessor
	giteaAccessor gitea.Accessor
	wikiAccessor  giteawiki.Accessor
	ticketID      int64
}

// CreateWikiDefaultConverter returns a Trac to Gitea markdown converter for converting wiki texts.
func CreateWikiDefaultConverter(tAccessor trac.Accessor, gAccessor gitea.Accessor, wAccessor giteawiki.Accessor) *DefaultConverter {
	converter := DefaultConverter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiAccessor: wAccessor, ticketID: -1}
	return &converter
}

// CreateTicketDefaultConverter returns a Trac to Gitea markdown converter for converting trac ticket descriptions and comments.
func CreateTicketDefaultConverter(tAccessor trac.Accessor, gAccessor gitea.Accessor, tracTicketID int64) *DefaultConverter {
	converter := DefaultConverter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiAccessor: nil, ticketID: tracTicketID}
	return &converter
}

func (converter *DefaultConverter) convertNonCodeBlockText(in string) string {
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
func (converter *DefaultConverter) Convert(in string) string {
	out := in

	// ensure we have Unix EOLs
	out = converter.convertEOL(out)

	// convert any code blocks
	out = converter.convertCodeBlocks(out)

	// perform all other conversions only on text not in a code block
	out = converter.convertNonCodeBlocks(out, converter.convertNonCodeBlockText)

	return out
}
