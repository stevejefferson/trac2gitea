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

// CreateConverter returns a new Trac to Gitea markdown converter.
func CreateConverter(tAccessor *trac.Accessor, gAccessor *gitea.Accessor, wAccessor *wiki.Accessor) *Converter {
	converter := Converter{tracAccessor: tAccessor, giteaAccessor: gAccessor, wikiAccessor: wAccessor}
	return &converter
}

// Convert a string of Trac markdown to Gitea markdown.
// linkPrefix is applied to any Trac link - e.g. a linkPrefix of "ticket:1" applied to a Trac link "image.png" will result in a markdown link "ticket:1:image.png"
func (converter *Converter) Convert(tracText string, linkPrefix string) string {
	out := converter.convertEOL(tracText)
	out = converter.convertCodeBlock(out)
	out = converter.convertImageReference(out, linkPrefix)
	out = converter.convertHeading(out)
	return out
}
