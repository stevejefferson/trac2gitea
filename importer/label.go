// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/trac"
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
func (importer *Importer) defaultLabelMap(tracMethod func(tAccessor trac.Accessor, handlerFn func(tracName string) error) error) (map[string]string, error) {
	labelMap := make(map[string]string)

	err := tracMethod(importer.tracAccessor, func(tracName string) error {
		// only interested in named trac items
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
	return importer.defaultLabelMap(trac.Accessor.GetComponentNames)
}

// DefaultPriorityLabelMap retrieves the default mapping between Trac priorities and Gitea labels
func (importer *Importer) DefaultPriorityLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetPriorityNames)
}

// DefaultResolutionLabelMap retrieves the default mapping between Trac resolutions and Gitea labels
func (importer *Importer) DefaultResolutionLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetResolutionNames)
}

// DefaultSeverityLabelMap retrieves the default mapping between Trac severities and Gitea labels
func (importer *Importer) DefaultSeverityLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetSeverityNames)
}

// DefaultTypeLabelMap retrieves the default mapping between Trac types and Gitea labels
func (importer *Importer) DefaultTypeLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetTypeNames)
}

// DefaultVersionLabelMap retrieves the default mapping between Trac versions and Gitea labels
func (importer *Importer) DefaultVersionLabelMap() (map[string]string, error) {
	return importer.defaultLabelMap(trac.Accessor.GetVersionNames)
}

// importLabels imports a single trac item as a Gitea label Gitea labels - any created label will have the provided color.
// Returns ID of Gitea label.
func (importer *Importer) importLabel(tracName string, labelMap map[string]string, labelColor string) (int64, error) {
	if tracName == "" {
		return -1, nil // ignore unnamed trac items
	}

	labelName := labelMap[tracName]
	if labelName == "" {
		return -1, nil // if no mapping provided, do not create a label
	}

	labelID, err := importer.giteaAccessor.AddLabel(labelName, labelColor)
	if err != nil {
		return -1, err
	}

	return labelID, nil
}

// ImportComponents imports Trac components as Gitea labels.
func (importer *Importer) ImportComponents(componentNameMap map[string]string) error {
	return importer.tracAccessor.GetComponentNames(func(componentName string) error {
		_, err := importer.importLabel(componentName, componentNameMap, componentLabelColor)
		return err
	})
}

// ImportPriorities imports Trac priorities as Gitea labels.
func (importer *Importer) ImportPriorities(priorityNameMap map[string]string) error {
	return importer.tracAccessor.GetPriorityNames(func(priorityName string) error {
		_, err := importer.importLabel(priorityName, priorityNameMap, priorityLabelColor)
		return err
	})
}

// ImportResolutions imports Trac resolutions as Gitea labels.
func (importer *Importer) ImportResolutions(resolutionNameMap map[string]string) error {
	return importer.tracAccessor.GetResolutionNames(func(resolutionName string) error {
		_, err := importer.importLabel(resolutionName, resolutionNameMap, resolutionLabelColor)
		return err
	})
}

// ImportSeverities imports Trac severities as Gitea labels.
func (importer *Importer) ImportSeverities(severityNameMap map[string]string) error {
	return importer.tracAccessor.GetSeverityNames(func(severityName string) error {
		_, err := importer.importLabel(severityName, severityNameMap, severityLabelColor)
		return err
	})
}

// ImportTypes imports Trac types as Gitea labels.
func (importer *Importer) ImportTypes(typeNameMap map[string]string) error {
	return importer.tracAccessor.GetTypeNames(func(typeName string) error {
		_, err := importer.importLabel(typeName, typeNameMap, typeLabelColor)
		return err
	})
}

// ImportVersions imports Trac versions as Gitea labels.
func (importer *Importer) ImportVersions(versionNameMap map[string]string) error {
	return importer.tracAccessor.GetVersionNames(func(versionName string) error {
		_, err := importer.importLabel(versionName, versionNameMap, versionLabelColor)
		return err
	})
}
