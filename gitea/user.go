package gitea

import (
	"database/sql"
	"log"
	"strings"
)

// GetUserID retireves the id of a named Gitea user - returns -1 if no such user.
func (accessor *Accessor) GetUserID(name string) int64 {
	if strings.Trim(name, " ") == "" {
		return -1
	}

	var id int64 = -1
	err := accessor.db.QueryRow(`SELECT id FROM user WHERE lower_name = $1 or email = $1`, name).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return id
}

func (accessor *Accessor) getAdminUserID() int64 {
	row := accessor.db.QueryRow(`
		SELECT id FROM user WHERE is_admin ORDER BY id LIMIT 1;
		`)

	var adminID int64
	err := row.Scan(&adminID)
	if err != nil {
		log.Fatal("No admin user found in Gitea")
	}

	return adminID
}

func (accessor *Accessor) getAdminDefaultingUserID(userName string, adminUserID int64) int64 {
	userID := adminUserID
	if userName != "" {
		userID = accessor.GetUserID(userName)
		if userID == -1 {
			log.Fatal("Cannot find gitea user ", userName)
		}
	}

	return userID
}
