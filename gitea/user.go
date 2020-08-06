package gitea

import (
	"log"
	"strings"
)

func (accessor *Accessor) FindUserID(nameOrAddress string) int64 {
	if strings.Trim(nameOrAddress, " ") == "" {
		return -1
	}

	var id int64
	err := accessor.db.QueryRow(`SELECT id FROM user WHERE lower_name = $1 or email = $1`, nameOrAddress).Scan(&id)
	if err != nil {
		return -1
	}

	return id
}

func (accessor *Accessor) findAdminUserID() int64 {
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

func (accessor *Accessor) findAdminDefaultingUserID(userName string, adminUserID int64) int64 {
	userID := adminUserID
	if userName != "" {
		userID = accessor.FindUserID(userName)
		if userID == -1 {
			log.Fatal("Cannot find gitea user ", userName)
		}
	}

	return userID
}
