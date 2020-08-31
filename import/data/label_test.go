// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package data_test

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
		}).
		AnyTimes()
}

func expectToAddLabels(t *testing.T) {
	// expect call to lookup ids of each of our (non-unnamed, non-removed) labels, return -1 as they don't exist
	mockGiteaAccessor.EXPECT().GetLabelID(labelName1).Return(int64(-1), nil).AnyTimes()
	mockGiteaAccessor.EXPECT().GetLabelID(labelName2).Return(int64(-1), nil).AnyTimes()

	// expect to add new labels on the basis of them not existing above
	mockGiteaAccessor.EXPECT().AddLabel(labelName1, gomock.Any()).Return(labelID1, nil).AnyTimes()
	mockGiteaAccessor.EXPECT().AddLabel(labelName2, gomock.Any()).Return(labelID2, nil).AnyTimes()
}

func expectToNotAddLabels(t *testing.T) {
	// expect call to lookup ids of each of our (non-unnamed, non-removed) labels, return ids because they exist
	mockGiteaAccessor.EXPECT().GetLabelID(labelName1).Return(labelID1, nil).AnyTimes()
	mockGiteaAccessor.EXPECT().GetLabelID(labelName2).Return(labelID2, nil).AnyTimes()

	// do not expect to add new labels...
}

func TestComponentsWhereNoLabelsExist(t *testing.T) {
	setUpComponents(t)
	defer tearDown(t)

	expectToAddLabels(t)

	importer.ImportComponents(labelMap)
}

func TestComponentsWhereLabelsExist(t *testing.T) {
	setUpComponents(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	importer.ImportComponents(labelMap)
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
		}).
		AnyTimes()
}

func TestPrioritiesWhereNoLabelsExist(t *testing.T) {
	setUpPriorities(t)
	defer tearDown(t)

	expectToAddLabels(t)

	importer.ImportPriorities(labelMap)
}

func TestPrioritiesWhereLabelsExist(t *testing.T) {
	setUpPriorities(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	importer.ImportPriorities(labelMap)
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
		}).
		AnyTimes()
}

func TestResolutionsWhereNoLabelsExist(t *testing.T) {
	setUpResolutions(t)
	defer tearDown(t)

	expectToAddLabels(t)

	importer.ImportResolutions(labelMap)
}

func TestResolutionsWhereLabelsExist(t *testing.T) {
	setUpResolutions(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	importer.ImportResolutions(labelMap)
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
		}).
		AnyTimes()
}

func TestSeveritiesWhereNoLabelsExist(t *testing.T) {
	setUpSeverities(t)
	defer tearDown(t)

	expectToAddLabels(t)

	importer.ImportSeverities(labelMap)
}

func TestSeveritiesWhereLabelsExist(t *testing.T) {
	setUpSeverities(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	importer.ImportSeverities(labelMap)
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
		}).
		AnyTimes()
}

func TestTypesWhereNoLabelsExist(t *testing.T) {
	setUpTypes(t)
	defer tearDown(t)

	expectToAddLabels(t)

	importer.ImportTypes(labelMap)
}

func TestTypesWhereLabelsExist(t *testing.T) {
	setUpTypes(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	importer.ImportTypes(labelMap)
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
		}).
		AnyTimes()
}

func TestVersionsWhereNoLabelsExist(t *testing.T) {
	setUpVersions(t)
	defer tearDown(t)

	expectToAddLabels(t)

	importer.ImportVersions(labelMap)
}

func TestVersionsWhereLabelsExist(t *testing.T) {
	setUpVersions(t)
	defer tearDown(t)

	expectToNotAddLabels(t)

	importer.ImportVersions(labelMap)
}
