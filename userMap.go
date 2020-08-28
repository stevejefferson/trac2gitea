// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/stevejefferson/trac2gitea/import/data"
)

// readUserMap reads the user map from the provided file, if no file provided, import a default map using the provided importer
func readUserMap(mapFile string, importer *data.Importer) (map[string]string, error) {
	if mapFile == "" {
		return importer.DefaultUserMap()
	}

	fd, err := os.Open(mapFile)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	userMap := make(map[string]string)
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		userMapLine := scanner.Text()
		equalsPos := strings.LastIndex(userMapLine, "=")
		if equalsPos == -1 {
			return nil, fmt.Errorf("badly formatted user map file %s: found line %s", mapFile, userMapLine)
		}

		tracUserName := strings.Trim(userMapLine[0:equalsPos], " ")
		giteaUserName := strings.Trim(userMapLine[equalsPos+1:], " ")
		userMap[tracUserName] = giteaUserName
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return userMap, nil
}

func writeUserMapToFile(mapFile string, userMap map[string]string) error {
	fd, err := os.Create(mapFile)
	if err != nil {
		return err
	}
	defer fd.Close()

	for tracUserName, giteaUserName := range userMap {
		if _, err := fd.WriteString(tracUserName + " = " + giteaUserName + "\n"); err != nil {
			return err
		}
	}

	return nil
}
