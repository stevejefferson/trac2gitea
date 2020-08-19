package gitea

import (
	"database/sql"
	"errors"
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
		log.Error(err)
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
func (accessor *DefaultAccessor) getAdminUserID() int64 {
	row := accessor.db.QueryRow(`
		SELECT id FROM user WHERE is_admin ORDER BY id LIMIT 1;
		`)

	var adminID int64
	err := row.Scan(&adminID)
	if err != nil {
		log.Fatal("No admin user found in Gitea\n")
	}

	return adminID
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
			err := errors.New("Cannot find gitea user  " + userName)
			log.Error(err)
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
		log.Error(err)
		return "", err
	}

	return emailAddress, nil
}
