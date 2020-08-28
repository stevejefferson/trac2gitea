// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/stevejefferson/trac2gitea/import/data"
	"github.com/stevejefferson/trac2gitea/import/issue"
)

const (
	componentTypeName  = "component"
	priorityTypeName   = "priority"
	resolutionTypeName = "resolution"
	severityTypeName   = "severity"
	typeTypeName       = "type"
	versionTypeName    = "version"
)

func readDefaultLabelMaps(importer *data.Importer) (componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string, err error) {
	componentMap, err = importer.DefaultComponentLabelMap()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	priorityMap, err = importer.DefaultPriorityLabelMap()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	resolutionMap, err = importer.DefaultResolutionLabelMap()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	severityMap, err = importer.DefaultSeverityLabelMap()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	typeMap, err = importer.DefaultTypeLabelMap()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	versionMap, err = importer.DefaultVersionLabelMap()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return
}

// readLabelMaps reads the label maps from the provided file, if no file provided, import default maps using the provided importer
func readLabelMaps(mapFile string, importer *issue.Importer) (componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string, err error) {
	if mapFile == "" {
		return readDefaultLabelMaps(importer)
	}

	fd, err := os.Open(mapFile)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}
	defer fd.Close()

	componentMap = make(map[string]string)
	priorityMap = make(map[string]string)
	resolutionMap = make(map[string]string)
	severityMap = make(map[string]string)
	typeMap = make(map[string]string)
	versionMap = make(map[string]string)

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		mapLine := scanner.Text()
		equalsPos := strings.LastIndex(mapLine, "=")
		if equalsPos == -1 {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("badly formatted label map file %s: expecting '=', found %s", mapFile, mapLine)
		}

		tracLabelAndType := strings.Trim(mapLine[0:equalsPos], " ")
		colonPos := strings.LastIndex(tracLabelAndType, ":")
		if equalsPos == -1 {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("badly formatted label map file %s: expecting ':', found %s", mapFile, mapLine)
		}
		labelType := strings.Trim(tracLabelAndType[0:colonPos], " ")
		tracLabel := strings.Trim(tracLabelAndType[colonPos+1:], " ")
		giteaLabel := strings.Trim(mapLine[equalsPos+1:], " ")

		switch labelType {
		case componentTypeName:
			componentMap[tracLabel] = giteaLabel
		case priorityTypeName:
			priorityMap[tracLabel] = giteaLabel
		case resolutionTypeName:
			resolutionMap[tracLabel] = giteaLabel
		case severityTypeName:
			severityMap[tracLabel] = giteaLabel
		case typeTypeName:
			typeMap[tracLabel] = giteaLabel
		case versionTypeName:
			versionMap[tracLabel] = giteaLabel
		default:
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("badly formatted label map file %s: expecting Trac label type before ':', found %s", mapFile, mapLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return
}

func writeLabelMapToFile(fd *os.File, labelType string, labelMap map[string]string) error {
	for tracLabel, giteaLabel := range labelMap {
		if _, err := fd.WriteString(labelType + ":" + tracLabel + " = " + giteaLabel + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func writeLabelMapsToFile(mapFile string, componentMap, priorityMap, resolutionMap, severityMap, typeMap, versionMap map[string]string) error {
	fd, err := os.Create(mapFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	writeLabelMapToFile(fd, componentTypeName, componentMap)
	writeLabelMapToFile(fd, priorityTypeName, priorityMap)
	writeLabelMapToFile(fd, resolutionTypeName, resolutionMap)
	writeLabelMapToFile(fd, severityTypeName, severityMap)
	writeLabelMapToFile(fd, typeTypeName, typeMap)
	writeLabelMapToFile(fd, versionTypeName, versionMap)

	return nil
}
