// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
)

const (
	tracItemNoNameChange = "no-change"
	tracItemRenamed      = "renamed"
	tracItemRemoved      = "removed"
	tracItemUnnamed      = ""

	labelName1 = tracItemNoNameChange
	labelName2 = "was-" + tracItemRenamed + "-now-something-else"
	labelName3 = ""

	labelID1 int64 = 1234
	labelID2 int64 = 2345
)

var labelMap map[string]string

func setUpLabels(t *testing.T) {
	setUp(t)

	// set label map to contain each type of trac item - unnamed is missing because unnamed trac items are ignored
	labelMap = make(map[string]string)
	labelMap[tracItemNoNameChange] = labelName1
	labelMap[tracItemRenamed] = labelName2
	labelMap[tracItemRemoved] = labelName3
}

func setUpComponents(t *testing.T) {
	setUpLabels(t)

	// expect trac accessor to return each of our trac items as component names
	mockTracAccessor.
		EXPECT().
		GetComponentNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(name string) error) error {
			handlerFn(tracItemNoNameChange)
			handlerFn(tracItemRenamed)
			handlerFn(tracItemRemoved)
			handlerFn(tracItemUnnamed)
			return nil
		})
}

func expectToAddLabels(t *testing.T) {
	// expect call to lookup ids of each of our (non-unnamed, non-removed) labels, return -1 as they don't exist
	mockGiteaAccessor.EXPECT().GetLabelID(labelName1).Return(int64(-1), nil)
	mockGiteaAccessor.EXPECT().GetLabelID(labelName2).Return(int64(-1), nil)

	// expect to add new labels on the basis of them not existing above
	mockGiteaAccessor.EXPECT().AddLabel(labelName1, gomock.Any()).Return(labelID1, nil)
	mockGiteaAccessor.EXPECT().AddLabel(labelName2, gomock.Any()).Return(labelID2, nil)
}

func expectToNotAddLabels(t *testing.T) {
	// expect call to lookup ids of each of our (non-unnamed, non-removed) labels, return ids because they exist
	mockGiteaAccessor.EXPECT().GetLabelID(labelName1).Return(labelID1, nil)
	mockGiteaAccessor.EXPECT().GetLabelID(labelName2).Return(labelID2, nil)
	// do not expect to add new labels...
}

func TestImportComponentsWhereNoLabelsExist(t *testing.T) {
	setUpComponents(t)
	defer tearDown(t)

	expectToAddLabels(t)

	dataImporter.ImportComponents(labelMap)
}

func TestImportComponentsWhereLabelsExist(t *testing.T) {
	setUpComponents(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	dataImporter.ImportComponents(labelMap)
}

func setUpPriorities(t *testing.T) {
	setUpLabels(t)

	// expect trac accessor to return each of our trac items as priority names
	mockTracAccessor.
		EXPECT().
		GetPriorityNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(name string) error) error {
			handlerFn(tracItemNoNameChange)
			handlerFn(tracItemRenamed)
			handlerFn(tracItemRemoved)
			handlerFn(tracItemUnnamed)
			return nil
		})
}

func TestImportPrioritiesWhereNoLabelsExist(t *testing.T) {
	setUpPriorities(t)
	defer tearDown(t)

	expectToAddLabels(t)

	dataImporter.ImportPriorities(labelMap)
}

func TestImportPrioritiesWhereLabelsExist(t *testing.T) {
	setUpPriorities(t)
	defer tearDown(t)

	expectLookupOfDefaultUser(t)
	expectToNotAddLabels(t)

	dataImporter.ImportPriorities(labelMap)
}

func setUpResolutions(t *testing.T) {
	setUpLabels(t)

	// expect trac accessor to return each of our trac items as resolution names
	mockTracAccessor.
		EXPECT().
		GetResolutionNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(name string) error) error {
			handlerFn(tracItemNoNameChange)
			handlerFn(tracItemRenamed)
			handlerFn(tracItemRemoved)
			handlerFn(tracItemUnnamed)
			return nil
		})
}

func TestImportResolutionsWhereNoLabelsExist(t *testing.T) {
	setUpResolutions(t)
	defer tearDown(t)

	expectToAddLabels(t)

	dataImporter.ImportResolutions(labelMap)
}

func TestImportResolutionsWhereLabelsExist(t *testing.T) {
	setUpResolutions(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	dataImporter.ImportResolutions(labelMap)
}

func setUpSeverities(t *testing.T) {
	setUpLabels(t)

	// expect trac accessor to return each of our trac items as severity names
	mockTracAccessor.
		EXPECT().
		GetSeverityNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(name string) error) error {
			handlerFn(tracItemNoNameChange)
			handlerFn(tracItemRenamed)
			handlerFn(tracItemRemoved)
			handlerFn(tracItemUnnamed)
			return nil
		})
}

func TestImportSeveritiesWhereNoLabelsExist(t *testing.T) {
	setUpSeverities(t)
	defer tearDown(t)

	expectToAddLabels(t)

	dataImporter.ImportSeverities(labelMap)
}

func TestImportSeveritiesWhereLabelsExist(t *testing.T) {
	setUpSeverities(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	dataImporter.ImportSeverities(labelMap)
}

func setUpTypes(t *testing.T) {
	setUpLabels(t)

	// expect trac accessor to return each of our trac items as type names
	mockTracAccessor.
		EXPECT().
		GetTypeNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(name string) error) error {
			handlerFn(tracItemNoNameChange)
			handlerFn(tracItemRenamed)
			handlerFn(tracItemRemoved)
			handlerFn(tracItemUnnamed)
			return nil
		})
}

func TestImportTypesWhereNoLabelsExist(t *testing.T) {
	setUpTypes(t)
	defer tearDown(t)

	expectToAddLabels(t)

	dataImporter.ImportTypes(labelMap)
}

func TestImportTypesWhereLabelsExist(t *testing.T) {
	setUpTypes(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	dataImporter.ImportTypes(labelMap)
}

func setUpVersions(t *testing.T) {
	setUpLabels(t)

	// expect trac accessor to return each of our trac items as version names
	mockTracAccessor.
		EXPECT().
		GetVersionNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(name string) error) error {
			handlerFn(tracItemNoNameChange)
			handlerFn(tracItemRenamed)
			handlerFn(tracItemRemoved)
			handlerFn(tracItemUnnamed)
			return nil
		})
}

func TestImportVersionsWhereNoLabelsExist(t *testing.T) {
	setUpVersions(t)
	defer tearDown(t)

	expectToAddLabels(t)

	dataImporter.ImportVersions(labelMap)
}

func TestImportVersionsWhereLabelsExist(t *testing.T) {
	setUpVersions(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	dataImporter.ImportVersions(labelMap)
}
