// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.
package issue

import "stevejefferson.co.uk/trac2gitea/log"

// ImportMilestones imports Trac milestones as Gitea milestones.
func (importer *Importer) ImportMilestones() error {
	return importer.tracAccessor.GetMilestones(func(name string, description string, due int64, completed int64) error {
		milestoneID, err := importer.giteaAccessor.GetMilestoneID(name)
		if err != nil {
			return err
		}
		if milestoneID != -1 {
			log.Debugf("milestone %s already exists - skipping...\n", name)
			return nil
		}

		milestoneID, err = importer.giteaAccessor.AddMilestone(name, description, completed != 0, due, completed)
		if err != nil {
			return err
		}

		log.Debugf("Added milestone (id %d) %s\n", milestoneID, name)
		return nil
	})
}
