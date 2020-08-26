// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"

	"github.com/stevejefferson/trac2gitea/log"
)

// getUserID retrieves the Gitea user ID and name corresponding to a Trac user name
func (importer *Importer) getUser(tracUser string) (int64, string, error) {
	// lookup Gitea user in map - the only reason for there not to be a mapping is with a faulty user-supplied map
	giteaUserName := importer.userMap[tracUser]
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
