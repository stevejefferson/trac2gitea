// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"fmt"
	"regexp"
	"strings"

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

		if matchedUserName == "" {
			matchedUserName = importer.giteaAccessor.GetCurrentUser()
			log.Debug("unmatched user \"%s\" mapped to default user \"%s\"", userName, matchedUserName)
		}

		userMap[user] = matchedUserName

		return nil
	})
	if err != nil {
		return nil, err
	}

	return userMap, nil
}

// getUserID retrieves the Gitea user ID and name corresponding to a Trac user name
func (importer *Importer) getUser(tracUser string, userMap map[string]string) (int64, string, error) {
	// lookup Gitea user in map - the only reason for there not to be a mapping is with a faulty user-supplied map
	giteaUserName := userMap[tracUser]
	if giteaUserName == "" {
		return -1, "", fmt.Errorf("cannot find mapping from Trac user %s to Gitea", tracUser)
	}

	userID, err := importer.giteaAccessor.GetUserID(giteaUserName)
	if err != nil {
		return -1, "", err
	}

	log.Debug("mapped Trac user %s onto Gitea user %s", tracUser, giteaUserName)
	return userID, giteaUserName, nil
}
