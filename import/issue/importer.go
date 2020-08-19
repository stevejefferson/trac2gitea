package issue

import (
	"stevejefferson.co.uk/trac2gitea/accessor/gitea"
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
)

// Importer of issue data from Trac tickets.
type Importer struct {
	giteaAccessor gitea.Accessor
	tracAccessor  trac.Accessor
}

// CreateImporter returns a new Trac ticket to Gitea issue importer.
func CreateImporter(
	tAccessor trac.Accessor,
	gAccessor gitea.Accessor) (*Importer, error) {
	importer := Importer{tracAccessor: tAccessor, giteaAccessor: gAccessor}
	return &importer, nil
}
