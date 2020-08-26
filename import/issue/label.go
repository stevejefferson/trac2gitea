// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
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
func (importer *Importer) importLabels(tracMethod func(tAccessor trac.Accessor, handlerFn func(name string) error) error, labelPrefix string, labelColor string) error {
	return tracMethod(importer.tracAccessor, func(name string) error {
		labelName := labelPrefix + name
		labelID, err := importer.giteaAccessor.GetLabelID(labelName)
		if err != nil {
			return err
		}
		if labelID != -1 {
			log.Debug("label %s already exists, skipping...", labelName)
			return nil
		}

		labelID, err = importer.giteaAccessor.AddLabel(labelName, labelColor)
		if err != nil {
			return err
		}

		log.Debug("created label (id %d), name %s, color %s", labelID, labelName, labelColor)
		return nil
	})
}

// ImportComponents imports Trac components as Gitea labels.
func (importer *Importer) ImportComponents() error {
	return importer.importLabels(trac.Accessor.GetComponentNames, componentLabelPrefix, componentLabelColor)
}

// ImportPriorities imports Trac priorities as Gitea labels.
func (importer *Importer) ImportPriorities() error {
	return importer.importLabels(trac.Accessor.GetPriorityNames, priorityLabelPrefix, priorityLabelColor)
}

// ImportSeverities imports Trac severities as Gitea labels.
func (importer *Importer) ImportSeverities() error {
	return importer.importLabels(trac.Accessor.GetSeverityNames, severityLabelPrefix, severityLabelColor)
}

// ImportVersions imports Trac versions as Gitea labels.
func (importer *Importer) ImportVersions() error {
	return importer.importLabels(trac.Accessor.GetVersionNames, versionLabelPrefix, versionLabelColor)
}

// ImportTypes imports Trac types as Gitea labels.
func (importer *Importer) ImportTypes() error {
	return importer.importLabels(trac.Accessor.GetTypeNames, typeLabelPrefix, typeLabelColor)
}

// ImportResolutions imports Trac resolutions as Gitea labels.
func (importer *Importer) ImportResolutions() error {
	return importer.importLabels(trac.Accessor.GetResolutionNames, resolutionLabelPrefix, resolutionLabelColor)
}
