// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"

	"github.com/golang/mock/gomock"
)

var labelMap map[string]string

var (
	tracUnchangedLabel *trac.Label
	tracRenamedLabel   *trac.Label
	tracRemovedLabel   *trac.Label
	tracUnnamedLabel   *trac.Label

	giteaUnchangedLabel *gitea.Label
	giteaRenamedLabel   *gitea.Label
)

func createTracLabel(name string, description string) *trac.Label {
	return &trac.Label{
		Name:        name,
		Description: description,
	}
}

func createGiteaLabel(name string, description string) *gitea.Label {
	return &gitea.Label{
		Name:        name,
		Description: description,
		Color:       "",
	}
}

func setUpLabels(t *testing.T) {
	setUp(t)

	tracUnchangedLabel = createTracLabel("unchanged", "unchanged-description")
	tracRenamedLabel = createTracLabel("renamed", "renamed-description")
	tracRemovedLabel = createTracLabel("removed", "removed-description")
	tracUnnamedLabel = createTracLabel("", "")

	giteaUnchangedLabel = createGiteaLabel(tracUnchangedLabel.Name, tracUnchangedLabel.Description)
	giteaRenamedLabel = createGiteaLabel("not-"+tracRenamedLabel.Name, tracRenamedLabel.Description)

	labelMap = make(map[string]string)
	labelMap[tracUnchangedLabel.Name] = giteaUnchangedLabel.Name
	labelMap[tracRenamedLabel.Name] = giteaRenamedLabel.Name
	labelMap[tracRemovedLabel.Name] = ""
}

func expectToReturnTracComponents(t *testing.T, components ...*trac.Label) {
	mockTracAccessor.
		EXPECT().
		GetComponents(gomock.Any()).
		DoAndReturn(func(handlerFn func(label *trac.Label) error) error {
			for _, component := range components {
				handlerFn(component)
			}
			return nil
		})
}

func expectToReturnTracPriorities(t *testing.T, priorities ...*trac.Label) {
	mockTracAccessor.
		EXPECT().
		GetPriorities(gomock.Any()).
		DoAndReturn(func(handlerFn func(label *trac.Label) error) error {
			for _, priority := range priorities {
				handlerFn(priority)
			}
			return nil
		})
}

func expectToReturnTracResolutions(t *testing.T, resolutions ...*trac.Label) {
	mockTracAccessor.
		EXPECT().
		GetResolutions(gomock.Any()).
		DoAndReturn(func(handlerFn func(label *trac.Label) error) error {
			for _, resolution := range resolutions {
				handlerFn(resolution)
			}
			return nil
		})
}

func expectToReturnTracSeverities(t *testing.T, severities ...*trac.Label) {
	mockTracAccessor.
		EXPECT().
		GetSeverities(gomock.Any()).
		DoAndReturn(func(handlerFn func(label *trac.Label) error) error {
			for _, severity := range severities {
				handlerFn(severity)
			}
			return nil
		})
}

func expectToReturnTracTypes(t *testing.T, types ...*trac.Label) {
	mockTracAccessor.
		EXPECT().
		GetTypes(gomock.Any()).
		DoAndReturn(func(handlerFn func(label *trac.Label) error) error {
			for _, typ := range types {
				handlerFn(typ)
			}
			return nil
		})
}

func expectToReturnTracVersions(t *testing.T, versions ...*trac.Label) {
	mockTracAccessor.
		EXPECT().
		GetVersions(gomock.Any()).
		DoAndReturn(func(handlerFn func(label *trac.Label) error) error {
			for _, version := range versions {
				handlerFn(version)
			}
			return nil
		})
}

// gomock Matcher for Gitea label names
type giteaLabelNameMatcher struct{ name string }

func isGiteaLabel(labelName string) gomock.Matcher {
	return giteaLabelNameMatcher{name: labelName}
}

func (matcher giteaLabelNameMatcher) Matches(arg interface{}) bool {
	giteaLabel := arg.(*gitea.Label)
	result := giteaLabel.Name == matcher.name
	return result
}

func (matcher giteaLabelNameMatcher) String() string {
	return "is Gitea label " + matcher.name
}

func expectToAddGiteaLabels(t *testing.T, giteaLabels ...*gitea.Label) {
	giteaLabelID := int64(666)
	for _, giteaLabel := range giteaLabels {
		giteaLabelName := giteaLabel.Name
		giteaLabelDescription := giteaLabel.Description
		mockGiteaAccessor.
			EXPECT().
			AddLabel(isGiteaLabel(giteaLabelName)).
			DoAndReturn(func(label *gitea.Label) (int64, error) {
				assertEquals(t, giteaLabelName, label.Name)
				assertEquals(t, giteaLabelDescription, label.Description)
				giteaLabelID++
				return giteaLabelID, nil
			})
	}
}

func TestImportComponents(t *testing.T) {
	setUpLabels(t)
	defer tearDown(t)

	expectToReturnTracComponents(t, tracUnchangedLabel, tracRenamedLabel, tracRemovedLabel, tracUnnamedLabel)
	expectToAddGiteaLabels(t, giteaUnchangedLabel, giteaRenamedLabel)

	dataImporter.ImportComponents(labelMap)
}

func TestImportPriorities(t *testing.T) {
	setUpLabels(t)
	defer tearDown(t)

	expectToReturnTracPriorities(t, tracUnchangedLabel, tracRenamedLabel, tracRemovedLabel, tracUnnamedLabel)
	expectToAddGiteaLabels(t, giteaUnchangedLabel, giteaRenamedLabel)

	dataImporter.ImportPriorities(labelMap)
}

func TestImportResolutions(t *testing.T) {
	setUpLabels(t)
	defer tearDown(t)

	expectToReturnTracResolutions(t, tracUnchangedLabel, tracRenamedLabel, tracRemovedLabel, tracUnnamedLabel)
	expectToAddGiteaLabels(t, giteaUnchangedLabel, giteaRenamedLabel)

	dataImporter.ImportResolutions(labelMap)
}

func TestImportSeverities(t *testing.T) {
	setUpLabels(t)
	defer tearDown(t)

	expectToReturnTracSeverities(t, tracUnchangedLabel, tracRenamedLabel, tracRemovedLabel, tracUnnamedLabel)
	expectToAddGiteaLabels(t, giteaUnchangedLabel, giteaRenamedLabel)

	dataImporter.ImportSeverities(labelMap)
}

func TestImportTypes(t *testing.T) {
	setUpLabels(t)
	defer tearDown(t)

	expectToReturnTracTypes(t, tracUnchangedLabel, tracRenamedLabel, tracRemovedLabel, tracUnnamedLabel)
	expectToAddGiteaLabels(t, giteaUnchangedLabel, giteaRenamedLabel)

	dataImporter.ImportTypes(labelMap)
}

func TestImportVersions(t *testing.T) {
	setUpLabels(t)
	defer tearDown(t)

	expectToReturnTracVersions(t, tracUnchangedLabel, tracRenamedLabel, tracRemovedLabel, tracUnnamedLabel)
	expectToAddGiteaLabels(t, giteaUnchangedLabel, giteaRenamedLabel)

	dataImporter.ImportVersions(labelMap)
}
