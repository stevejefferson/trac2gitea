// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
)

// label colors - courtesy of `trac2gogs`
const (
	componentLabelColor  = "#fbca04"
	priorityLabelColor   = "#207de5"
	resolutionLabelColor = "#9e9e9e"
	severityLabelColor   = "#eb6420"
	typeLabelColor       = "#e11d21"
	versionLabelColor    = "#009800"
)

// defaultLabelMap retrieves the default mapping between the Trac items returned by the provided function and Gitea labels
func (importer *Importer) defaultLabelMap(tracMethod func(tAccessor trac.Accessor, handlerFn func(tracLabel *trac.Label) error) error) (map[string]string, error) {
	labelMap := make(map[string]string)

	err := tracMethod(importer.tracAccessor, func(tracLabel *trac.Label) error {
		// only interested in named trac items
		tracName := tracLabel.Name
		if tracName != "" {
			labelMap[tracName] = tracName
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return labelMap, nil
}

// DefaultComponentLabelMap retrieves the default mapping between Trac components and Gitea labels
func (importer *Importer) DefaultComponentLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetComponents)
}

// DefaultPriorityLabelMap retrieves the default mapping between Trac priorities and Gitea labels
func (importer *Importer) DefaultPriorityLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetPriorities)
}

// DefaultResolutionLabelMap retrieves the default mapping between Trac resolutions and Gitea labels
func (importer *Importer) DefaultResolutionLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetResolutions)
}

// DefaultSeverityLabelMap retrieves the default mapping between Trac severities and Gitea labels
func (importer *Importer) DefaultSeverityLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetSeverities)
}

// DefaultTypeLabelMap retrieves the default mapping between Trac types and Gitea labels
func (importer *Importer) DefaultTypeLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetTypes)
}

// DefaultVersionLabelMap retrieves the default mapping between Trac versions and Gitea labels
func (importer *Importer) DefaultVersionLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetVersions)
}

// getLabelID retrieves the Gitea label ID corresponding to a Trac label name
func (importer *Importer) getLabelID(tracName string, labelMap map[string]string) (int64, error) {
	giteaLabelName := labelMap[tracName]
	if giteaLabelName == "" {
		return gitea.NullID, nil
	}

	labelID, err := importer.giteaAccessor.GetLabelID(giteaLabelName)
	if err != nil {
		return gitea.NullID, err
	}

	log.Debug("mapped Trac label %s onto Gitea label %s", tracName, giteaLabelName)
	return labelID, nil
}

// importLabels imports a single trac label as a Gitea label - any created label will have the provided color.
// Returns ID of Gitea label.
func (importer *Importer) importLabel(tracLabel *trac.Label, labelMap map[string]string, labelColor string) (int64, error) {
	tracName := tracLabel.Name
	if tracName == "" {
		return gitea.NullID, nil // ignore unnamed trac items
	}

	giteaLabelName := labelMap[tracName]
	if giteaLabelName == "" {
		return gitea.NullID, nil // if no mapping provided, do not create a label
	}

	giteaLabel := gitea.Label{Name: giteaLabelName, Description: tracLabel.Description, Color: labelColor}
	labelID, err := importer.giteaAccessor.AddLabel(&giteaLabel)
	if err != nil {
		return gitea.NullID, err
	}

	return labelID, nil
}

// ImportComponents imports Trac components as Gitea labels.
func (importer *Importer) ImportComponents(componentNameMap map[string]string) error {
	return importer.tracAccessor.GetComponents(func(component *trac.Label) error {
		_, err := importer.importLabel(component, componentNameMap, componentLabelColor)
		return err
	})
}

// ImportPriorities imports Trac priorities as Gitea labels.
func (importer *Importer) ImportPriorities(priorityNameMap map[string]string) error {
	return importer.tracAccessor.GetPriorities(func(priority *trac.Label) error {
		_, err := importer.importLabel(priority, priorityNameMap, priorityLabelColor)
		return err
	})
}

// ImportResolutions imports Trac resolutions as Gitea labels.
func (importer *Importer) ImportResolutions(resolutionNameMap map[string]string) error {
	return importer.tracAccessor.GetResolutions(func(resolution *trac.Label) error {
		_, err := importer.importLabel(resolution, resolutionNameMap, resolutionLabelColor)
		return err
	})
}

// ImportSeverities imports Trac severities as Gitea labels.
func (importer *Importer) ImportSeverities(severityNameMap map[string]string) error {
	return importer.tracAccessor.GetSeverities(func(severity *trac.Label) error {
		_, err := importer.importLabel(severity, severityNameMap, severityLabelColor)
		return err
	})
}

// ImportTypes imports Trac types as Gitea labels.
func (importer *Importer) ImportTypes(typeNameMap map[string]string) error {
	return importer.tracAccessor.GetTypes(func(tracType *trac.Label) error {
		_, err := importer.importLabel(tracType, typeNameMap, typeLabelColor)
		return err
	})
}

// ImportVersions imports Trac versions as Gitea labels.
func (importer *Importer) ImportVersions(versionNameMap map[string]string) error {
	return importer.tracAccessor.GetVersions(func(version *trac.Label) error {
		_, err := importer.importLabel(version, versionNameMap, versionLabelColor)
		return err
	})
}
