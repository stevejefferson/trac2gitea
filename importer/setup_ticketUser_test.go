// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
)

/*
 * Set up for ticket/issue user parts of ticket tests.
 * Contains:
 * - ticket user types
 * - expectations for use with ticket users.
 */

// TicketUserImport holds the data on a user referenced by an imported ticket
type TicketUserImport struct {
	tracUser     string
	giteaUser    string
	giteaUserID  int64
	origTracUser string
}

func createTicketUserImport(tracUser string, giteaUser string) *TicketUserImport {
	// if we have a gitea user mapping use a unique user id but do not record the original trac user on issues
	// if not, expect default user id to be used and record original trac user
	var giteaUserID int64
	var origTracUser string
	if giteaUser != "" {
		giteaUserID = allocateID()
		origTracUser = ""
	} else {
		giteaUserID = defaultUserID
		origTracUser = tracUser
	}

	user := TicketUserImport{
		tracUser:     tracUser,
		giteaUser:    giteaUser,
		giteaUserID:  giteaUserID,
		origTracUser: origTracUser,
	}

	if tracUser != "" {
		userMap[user.tracUser] = user.giteaUser
	}

	return &user
}

func expectUserLookup(t *testing.T, user *TicketUserImport) {
	// only expect user lookup if we have a trac -> gitea user mapping
	if user.tracUser == "" || user.giteaUser == "" {
		return
	}

	mockGiteaAccessor.
		EXPECT().
		GetUserID(gomock.Eq(user.giteaUser)).
		Return(user.giteaUserID, nil)
}

func expectIssueUserToBeAdded(t *testing.T, ticket *TicketImport, user *TicketUserImport) {
	mockGiteaAccessor.
		EXPECT().
		AddIssueUser(gomock.Eq(ticket.issueID), gomock.Eq(user.giteaUserID)).
		Return(nil)
}

func expectIssueAssigneeToBeAdded(t *testing.T, ticket *TicketImport, user *TicketUserImport) {
	if user.giteaUser != "" {
		mockGiteaAccessor.
			EXPECT().
			AddIssueAssignee(gomock.Eq(ticket.issueID), gomock.Eq(user.giteaUserID)).
			Return(nil)
	}
}
