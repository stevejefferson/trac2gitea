// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

const (
	unnamedMilestoneName     = ""
	completedMilestoneName   = "completed"
	uncompletedMilestoneName = "uncompleted"

	unnamedMilestoneDescription     = "n/a"
	completedMilestoneDescription   = "this is a completed milestone"
	uncompletedMilestoneDescription = "this milestone has not been completed"

	unnamedMilestoneDueTime     = int64(12345)
	completedMilestoneDueTime   = int64(23456)
	uncompletedMilestoneDueTime = int64(34567)

	unnamedMilestoneCompletedTime     = int64(54321)
	completedMilestoneCompletedTime   = int64(65432)
	uncompletedMilestoneCompletedTime = int64(0)

	completedMilestoneID   = int64(111)
	uncompletedMilestoneID = int64(222)
)

var (
	tracUnnamedMilestone     trac.Milestone
	tracCompletedMilestone   trac.Milestone
	tracUncompletedMilestone trac.Milestone
)

func setUpMilestones(t *testing.T) {
	setUp(t)

	tracUnnamedMilestone = trac.Milestone{
		Name:        unnamedMilestoneName,
		Description: unnamedMilestoneDescription,
		Due:         unnamedMilestoneDueTime,
		Completed:   unnamedMilestoneCompletedTime}

	tracCompletedMilestone = trac.Milestone{
		Name:        completedMilestoneName,
		Description: completedMilestoneDescription,
		Due:         completedMilestoneDueTime,
		Completed:   completedMilestoneCompletedTime}

	tracUncompletedMilestone = trac.Milestone{
		Name:        uncompletedMilestoneName,
		Description: uncompletedMilestoneDescription,
		Due:         uncompletedMilestoneDueTime,
		Completed:   uncompletedMilestoneCompletedTime}

	// expect trac accessor to return each of our trac milestones
	mockTracAccessor.
		EXPECT().
		GetMilestones(gomock.Any()).
		DoAndReturn(func(handlerFn func(milestone *trac.Milestone) error) error {
			handlerFn(&tracUnnamedMilestone)
			handlerFn(&tracCompletedMilestone)
			handlerFn(&tracUncompletedMilestone)
			return nil
		})
}

// gomock Matcher for milestone names
type milestoneNameMatcher struct{ name string }

func isMilestone(milestoneName string) gomock.Matcher {
	return milestoneNameMatcher{name: milestoneName}
}

func (matcher milestoneNameMatcher) Matches(arg interface{}) bool {
	giteaMilestone := arg.(*gitea.Milestone)
	return giteaMilestone.Name == matcher.name
}

func (matcher milestoneNameMatcher) String() string {
	return "is Gitea milestone " + matcher.name
}

func TestMilestonesWhereNoneExist(t *testing.T) {
	setUpMilestones(t)
	defer tearDown(t)

	// expect call to lookup ids of each of our (non-unnamed) milestones, return -1 as they don't exist
	mockGiteaAccessor.EXPECT().GetMilestoneID(completedMilestoneName).Return(int64(-1), nil)
	mockGiteaAccessor.EXPECT().GetMilestoneID(uncompletedMilestoneName).Return(int64(-1), nil)

	// expect to add new milestones on the basis of them not existing above
	mockGiteaAccessor.
		EXPECT().
		AddMilestone(isMilestone(completedMilestoneName)).
		DoAndReturn(func(giteaMilestone *gitea.Milestone) (int64, error) {
			assertEquals(t, giteaMilestone.Description, completedMilestoneDescription)
			assertEquals(t, giteaMilestone.Closed, true)
			assertEquals(t, giteaMilestone.DueTime, completedMilestoneDueTime)
			assertEquals(t, giteaMilestone.ClosedTime, completedMilestoneCompletedTime)
			return completedMilestoneID, nil
		})
	mockGiteaAccessor.
		EXPECT().
		AddMilestone(isMilestone(uncompletedMilestoneName)).
		DoAndReturn(func(giteaMilestone *gitea.Milestone) (int64, error) {
			assertEquals(t, giteaMilestone.Description, uncompletedMilestoneDescription)
			assertEquals(t, giteaMilestone.Closed, false)
			assertEquals(t, giteaMilestone.DueTime, uncompletedMilestoneDueTime)
			assertEquals(t, giteaMilestone.ClosedTime, uncompletedMilestoneCompletedTime)
			return uncompletedMilestoneID, nil
		})

	dataImporter.ImportMilestones()
}

func TestMilestonesWhereAlreadyExist(t *testing.T) {
	setUpMilestones(t)
	defer tearDown(t)

	// expect call to lookup ids of each of our (non-unnamed) milestones, return ids because they exist
	mockGiteaAccessor.EXPECT().GetMilestoneID(completedMilestoneName).Return(completedMilestoneID, nil)
	mockGiteaAccessor.EXPECT().GetMilestoneID(uncompletedMilestoneName).Return(uncompletedMilestoneID, nil)

	// do not expect to add new milestones...

	dataImporter.ImportMilestones()
}
