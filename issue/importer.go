package issue

import (
	"stevejefferson.co.uk/trac2gitea/gitea"
	"stevejefferson.co.uk/trac2gitea/markdown"
	"stevejefferson.co.uk/trac2gitea/trac"
)

type Importer struct {
	giteaAccessor          *gitea.Accessor
	tracAccessor           *trac.Accessor
	trac2MarkdownConverter *markdown.Converter
}

func CreateImporter(
	tAccessor *trac.Accessor,
	gAccessor *gitea.Accessor,
	t2mConverter *markdown.Converter) *Importer {
	importer := Importer{tracAccessor: tAccessor, giteaAccessor: gAccessor, trac2MarkdownConverter: t2mConverter}
	return &importer
}
