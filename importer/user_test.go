// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer_test

import (
	"testing"

	"github.com/golang/mock/gomock"
)

const (
	noEmailUserName    = "user1"
	noEmailUser        = noEmailUserName
	matchedNoEmailUser = "matched-user1"

	emailUserName    = "user2"
	emailUserEmail   = "u2@abc.def"
	emailUser        = emailUserName + " <" + emailUserEmail + ">"
	matchedEmailUser = "matched-user2"

	noMatchUserName  = "user3"
	noMatchUserEmail = "u3@ghi.jkl"
	noMatchUser      = noMatchUserName + " <" + noMatchUserEmail + ">"
)

func expectToRetrieveTracUsers(t *testing.T, users ...string) {
	mockTracAccessor.
		EXPECT().
		GetUserNames(gomock.Any()).
		DoAndReturn(func(handlerFn func(u string) error) error {
			for _, user := range users {
				handlerFn(user)
			}
			return nil
		})
}

func expectMatchUser(t *testing.T, userName string, userEMail string, matchedUser string) {
	mockGiteaAccessor.
		EXPECT().
		MatchUser(gomock.Eq(userName), gomock.Eq(userEMail)).
		Return(matchedUser, nil)
}

func TestDefaultUserMapForUserWithNoEmail(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	expectToRetrieveTracUsers(t, noEmailUser)
	expectMatchUser(t, noEmailUserName, "", matchedNoEmailUser)
	userMap, _ := dataImporter.DefaultUserMap()
	assertEquals(t, userMap[noEmailUser], matchedNoEmailUser)
}

func TestDefaultUserMapForUserWithEmail(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	expectToRetrieveTracUsers(t, emailUser)
	expectMatchUser(t, emailUserName, emailUserEmail, matchedEmailUser)
	userMap, _ := dataImporter.DefaultUserMap()
	assertEquals(t, userMap[emailUser], matchedEmailUser)
}

func TestDefaultUserMapForUnmatchedUser(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	expectToRetrieveTracUsers(t, noMatchUser)
	expectMatchUser(t, noMatchUserName, noMatchUserEmail, "")

	userMap, _ := dataImporter.DefaultUserMap()
	assertEquals(t, userMap[noMatchUser], "")
}

func TestDefaultUserMapForMultipleUsers(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	expectToRetrieveTracUsers(t, noEmailUser, emailUser, noMatchUser)
	expectMatchUser(t, noEmailUserName, "", matchedNoEmailUser)
	expectMatchUser(t, emailUserName, emailUserEmail, matchedEmailUser)
	expectMatchUser(t, noMatchUserName, noMatchUserEmail, "")

	userMap, _ := dataImporter.DefaultUserMap()
	assertEquals(t, userMap[noEmailUser], matchedNoEmailUser)
	assertEquals(t, userMap[emailUser], matchedEmailUser)
	assertEquals(t, userMap[noMatchUser], "")
}
