// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"fmt"
	"time"

	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

// Importer of issue data from Trac tickets.
type Importer struct {
	giteaAccessor gitea.Accessor
	tracAccessor  trac.Accessor
	userMap       map[string]string
}

// CreateImporter returns a new Trac ticket to Gitea issue importer.
func CreateImporter(
	tAccessor trac.Accessor,
	gAccessor gitea.Accessor,
	uMap map[string]string) (*Importer, error) {
	importer := Importer{tracAccessor: tAccessor, giteaAccessor: gAccessor, userMap: uMap}
	return &importer, nil
}

// addTracContext adds context information to the provided message giving the original Trac context of the data in the message.
func addTracContext(tracContext string, tracUpdateTime int64, message string) string {
	updateTimeStr := time.Unix(tracUpdateTime, 0)
	contextMessage := fmt.Sprintf("[Imported from Trac: %s at %s]\n\n%s", tracContext, updateTimeStr, message)
	return contextMessage
}
