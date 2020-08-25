// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/stevejefferson/trac2gitea/log"
)

// GetUserID retrieves the id of a named Gitea user - returns -1 if no such user.
func (accessor *DefaultAccessor) GetUserID(userName string) (int64, error) {
	if strings.Trim(userName, " ") == "" {
		return -1, nil
	}

	var id int64 = -1
	err := accessor.db.QueryRow(`SELECT id FROM user WHERE lower_name = $1 or email = $1`, userName).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Error("Cannot look up user %s: %v\n", userName, err)
		return -1, err
	}

	return id, nil
}

// GetUserEMailAddress retrieves the email address of a given user
func (accessor *DefaultAccessor) GetUserEMailAddress(userName string) (string, error) {
	var emailAddress string = ""
	err := accessor.db.QueryRow(`SELECT email FROM user WHERE lower_name = $1`, userName).Scan(&emailAddress)
	if err != nil && err != sql.ErrNoRows {
		log.Error("Problem finding email address for user %s: %v\n", userName, err)
		return "", err
	}

	return emailAddress, nil
}

// getUserRepoURL retrieves the URL of the current repository for the current user
func (accessor *DefaultAccessor) getUserRepoURL() string {
	rootURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", rootURL, accessor.userName, accessor.repoName)
}

// matchUser retrieves the name of the user best matching a user name or email address
func (accessor *DefaultAccessor) matchUser(userName string, userEmail string) (string, error) {
	var matchedUserName = ""
	lcUserName := strings.ToLower(userName)
	err := accessor.db.QueryRow(`
		SELECT lower_name FROM user 
		WHERE lower_name = $1 
		OR full_name = $2 
		OR email = $3`, lcUserName, userName, userEmail).Scan(&matchedUserName)
	if err != nil && err != sql.ErrNoRows {
		log.Error("Problem matching user name %s, email %s: %v\n", userName, userEmail, err)
		return "", err
	}

	return matchedUserName, nil
}

// regexp for matching a user: $1=username (may have space padding) $2=user email (optional)
var userRegexp = regexp.MustCompile(`([^<]*)(?:<([^>]+)>)?`)

// GenerateDefaultUserMappings populates the provided user map with a default mapping for each user in the map.
func (accessor *DefaultAccessor) GenerateDefaultUserMappings(userMap map[string]string, defaultUserName string) error {
	for user := range userMap {
		userName := userRegexp.ReplaceAllString(user, `$1`)
		trimmedUserName := strings.Trim(userName, " ")
		userEmail := userRegexp.ReplaceAllString(user, `$2`)

		matchedUserName, err := accessor.matchUser(trimmedUserName, userEmail)
		log.Debug("Matched user \"%s\", email \"%s\" to \"%s\"\n", userName, userEmail, matchedUserName)
		if err != nil {
			return err
		}

		if matchedUserName == "" {
			matchedUserName = defaultUserName
			log.Debug("Mapping unmatched user \"%s\" to default user \"%s\"\n", userName, defaultUserName)
		}

		userMap[user] = matchedUserName
	}

	return nil
}
