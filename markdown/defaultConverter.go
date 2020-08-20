// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package markdown

import (
	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
)

// DefaultConverter is the default implementation of the Trac markdown to Gitea markdown converter.
// This is used in two circumstances:
// 1. for ticket comments - in which case ticketID != -1 and wikiAccessor == nil
// 2. for wiki imports - in which case ticketID == -1 and wikiAccessor != nil
type DefaultConverter struct {
	tracAccessor  trac.Accessor
	giteaAccessor gitea.Accessor
	wikiPage      string
	ticketID      int64
}

// CreateWikiDefaultConverter returns a Trac to Gitea markdown converter for converting wiki texts.
func CreateWikiDefaultConverter(tAccessor trac.Accessor, gAccessor gitea.Accessor, wikiPageName string) *DefaultConverter {
	converter := DefaultConverter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiPage: wikiPageName, ticketID: -1}
	return &converter
}

// CreateTicketDefaultConverter returns a Trac to Gitea markdown converter for converting trac ticket descriptions and comments.
func CreateTicketDefaultConverter(tAccessor trac.Accessor, gAccessor gitea.Accessor, tracTicketID int64) *DefaultConverter {
	converter := DefaultConverter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiPage: "", ticketID: tracTicketID}
	return &converter
}

func (converter *DefaultConverter) convertNonCodeBlockText(in string) string {
	out := in

	// do simple one-line constructs first
	out = converter.convertLinks(out)
	out = converter.convertAnchors(out)
	out = converter.convertEscapes(out)
	out = converter.convertLists(out)
	out = converter.convertDefinitionLists(out)
	out = converter.convertHeadings(out)

	// font styles require links to be disguised
	out = converter.disguiseLinks(out)
	out = converter.convertFontStyles(out)
	out = converter.undisguiseLinks(out)

	// now do potentially more complex constructs
	out = converter.convertBlockQuotes(out)
	out = converter.convertTables(out)

	// do paragraphs last because this results in the insertion of newlines
	// - this can upset other regexps which assume everything on a single line
	out = converter.convertParagraphs(out)

	return out
}

// Convert converts a string of Trac markdown to Gitea markdown
func (converter *DefaultConverter) Convert(in string) string {
	out := in

	// ensure we have Unix EOLs
	out = converter.convertEOL(out)

	// perform conversions on text not in a code block
	out = converter.convertNonCodeBlocks(out, converter.convertNonCodeBlockText)

	// finally, convert any code blocks
	// - this must be done after the non-code block conversions otherwise code blocks would get converted
	// and we wouldn't be able to recognise the code block boundaries
	out = converter.convertCodeBlocks(out)

	return out
}
