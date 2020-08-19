package issue

import (
	"stevejefferson.co.uk/trac2gitea/accessor/trac"
	"stevejefferson.co.uk/trac2gitea/log"
)

// label colors and prefixes - courtesy of `trac2gogs`
const (
	componentLabelPrefix  = "Component/"
	componentLabelColor   = "#fbca04"
	priorityLabelPrefix   = "Priority/"
	priorityLabelColor    = "#207de5"
	severityLabelPrefix   = "Severity/"
	severityLabelColor    = "#eb6420"
	versionLabelPrefix    = "Version/"
	versionLabelColor     = "#009800"
	typeLabelPrefix       = "Type/"
	typeLabelColor        = "#e11d21"
	resolutionLabelPrefix = "Resolution/"
	resolutionLabelColor  = "#9e9e9e"
)

// importLabels imports all labels returned by the provided Trac accessor method as Gitea labels with the provided prefix and color.
func (importer *Importer) importLabels(tracMethod func(tAccessor trac.Accessor, handlerFn func(name string)), labelPrefix string, labelColor string) {
	tracMethod(importer.tracAccessor, func(name string) {
		labelName := labelPrefix + name
		if importer.giteaAccessor.GetLabelID(labelName) != -1 {
			log.Debugf("label %s already exists, skipping...\n", labelName)
			return
		}

		labelID := importer.giteaAccessor.AddLabel(labelName, labelColor)
		log.Debugf("Created label (id %d), name %s, color %s\n", labelID, labelName, labelColor)
	})
}

// ImportComponents imports Trac components as Gitea labels.
func (importer *Importer) ImportComponents() {
	importer.importLabels(trac.Accessor.GetComponentNames, componentLabelPrefix, componentLabelColor)
}

// ImportPriorities imports Trac priorities as Gitea labels.
func (importer *Importer) ImportPriorities() {
	importer.importLabels(trac.Accessor.GetPriorityNames, priorityLabelPrefix, priorityLabelColor)
}

// ImportSeverities imports Trac severities as Gitea labels.
func (importer *Importer) ImportSeverities() {
	importer.importLabels(trac.Accessor.GetSeverityNames, severityLabelPrefix, severityLabelColor)
}

// ImportVersions imports Trac versions as Gitea labels.
func (importer *Importer) ImportVersions() {
	importer.importLabels(trac.Accessor.GetVersionNames, versionLabelPrefix, versionLabelColor)
}

// ImportTypes imports Trac types as Gitea labels.
func (importer *Importer) ImportTypes() {
	importer.importLabels(trac.Accessor.GetTypeNames, typeLabelPrefix, typeLabelColor)
}

// ImportResolutions imports Trac resolutions as Gitea labels.
func (importer *Importer) ImportResolutions() {
	importer.importLabels(trac.Accessor.GetResolutionNames, resolutionLabelPrefix, resolutionLabelColor)
}
