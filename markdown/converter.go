package markdown

import (
	"stevejefferson.co.uk/trac2gitea/gitea"
	"stevejefferson.co.uk/trac2gitea/trac"
	"stevejefferson.co.uk/trac2gitea/wiki"
)

// Converter of Trac markdown to Gitea markdown
type Converter struct {
	tracAccessor  *trac.Accessor
	giteaAccessor *gitea.Accessor
	wikiAccessor  *wiki.Accessor
}

// CreateConverter returns a general-purpose Trac to Gitea markdown converter.
func CreateConverter(tAccessor *trac.Accessor, gAccessor *gitea.Accessor, wAccessor *wiki.Accessor) *Converter {
	converter := Converter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiAccessor: wAccessor}
	return &converter
}

// Convert converts a string of Trac markdown to Gitea markdown.unassociated with any ticket
func (converter *Converter) Convert(in string) string {
	return converter.TicketConvert(in, 0)
}

// TicketConvert converts a string of Trac markdown associated with a given Trac ticket to Gitea markdown
func (converter *Converter) TicketConvert(in string, ticketID int64) string {
	out := in
	out = converter.convertEOL(out)
	out = converter.convertLink(out, ticketID)
	out = converter.disguiseLinks(out)
	out = converter.convertCodeBlock(out)
	out = converter.convertHeading(out)
	out = converter.convertFontStyle(out)
	out = converter.undisguiseLinks(out)
	return out
}
