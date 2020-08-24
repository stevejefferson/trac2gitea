// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package gitea

import (
	"database/sql"
	"fmt"
	"strings"

	"stevejefferson.co.uk/trac2gitea/log"
)

// GetUserID retrieves the id of a named Gitea user - returns -1 if no such user.
func (accessor *DefaultAccessor) GetUserID(name string) (int64, error) {
	if strings.Trim(name, " ") == "" {
		return -1, nil
	}

	var id int64 = -1
	err := accessor.db.QueryRow(`SELECT id FROM user WHERE lower_name = $1 or email = $1`, name).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Error("Cannot look up user %s: %v\n", name, err)
		return -1, err
	}

	return id, nil
}

// GetDefaultAssigneeID retrieves the id of the user to which to assign tickets/comments in the case where the Trac-supplied user id does not exist in Gitea.
func (accessor *DefaultAccessor) GetDefaultAssigneeID() int64 {
	return accessor.defaultAssigneeID
}

// GetDefaultAuthorID retrieves the id of the user to set as the author of tickets/comments in the case where the Trac-supplied user id does not exist in Gitea.
func (accessor *DefaultAccessor) GetDefaultAuthorID() int64 {
	return accessor.defaultAuthorID
}

// getAdminUserID retrieves the id of the project admin user.
func (accessor *DefaultAccessor) getAdminUserID() (int64, error) {
	row := accessor.db.QueryRow(`
		SELECT id FROM user WHERE is_admin ORDER BY id LIMIT 1;
		`)

	var adminID int64
	err := row.Scan(&adminID)
	if err != nil {
		err = fmt.Errorf("No admin user found in Gitea")
		log.Error("%v\n", err)
		return -1, err
	}

	return adminID, nil
}

// getAdminDefaultingUserID retrieves the id of a named user, defaulting to the admin user if that user does not exist.
func (accessor *DefaultAccessor) getAdminDefaultingUserID(userName string, adminUserID int64) (int64, error) {
	userID := adminUserID
	if userName != "" {
		userID, err := accessor.GetUserID(userName)
		if err != nil {
			return -1, err
		}
		if userID == -1 {
			err := fmt.Errorf("Cannot find gitea user %s", userName)
			log.Error("%v\n", err)
			return -1, err
		}
	}

	return userID, nil
}

// GetUserEMailAddress retrieves the email address of a given user
func (accessor *DefaultAccessor) GetUserEMailAddress(userID int64) (string, error) {
	var emailAddress string = ""
	err := accessor.db.QueryRow(`SELECT email FROM user WHERE id = $1`, userID).Scan(&emailAddress)
	if err != nil && err != sql.ErrNoRows {
		log.Error("Problem finding email address for user %d: %v\n", userID, err)
		return "", err
	}

	return emailAddress, nil
}

// getUserRepoURL retrieves the URL of the current repository for the current user
func (accessor *DefaultAccessor) getUserRepoURL() string {
	rootURL := accessor.GetStringConfig("server", "ROOT_URL")
	return fmt.Sprintf("%s/%s/%s", rootURL, accessor.userName, accessor.repoName)
}
