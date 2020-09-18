// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"regexp"
	"strings"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/log"
)

// regexp for matching a user: $1=username (may have space padding) $2=user email (optional)
var userRegexp = regexp.MustCompile(`([^<]*)(?:<([^>]+)>)?`)

// DefaultUserMap retrieves the default mapping between Trac users and Gitea users
func (importer *Importer) DefaultUserMap() (map[string]string, error) {
	userMap := make(map[string]string)

	err := importer.tracAccessor.GetUserNames(func(user string) error {
		userName := userRegexp.ReplaceAllString(user, `$1`)
		trimmedUserName := strings.Trim(userName, " ")
		userEmail := userRegexp.ReplaceAllString(user, `$2`)

		matchedUserName, err := importer.giteaAccessor.MatchUser(trimmedUserName, userEmail)
		if err != nil {
			return err
		}
		log.Debug("matched user \"%s\", email \"%s\" to \"%s\"", userName, userEmail, matchedUserName)

		userMap[user] = matchedUserName

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userMap, nil
}

// getUserID retrieves the Gitea user ID corresponding to a Trac user name
func (importer *Importer) getUserID(tracUser string, userMap map[string]string) (int64, error) {
	giteaUserName := userMap[tracUser]
	if giteaUserName == "" {
		return gitea.NullID, nil
	}

	userID, err := importer.giteaAccessor.GetUserID(giteaUserName)
	if err != nil {
		return gitea.NullID, err
	}

	log.Debug("mapped Trac user %s onto Gitea user %s", tracUser, giteaUserName)
	return userID, nil
}
