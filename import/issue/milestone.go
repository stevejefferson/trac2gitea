// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
)

// ImportMilestones imports Trac milestones as Gitea milestones.
func (importer *Importer) ImportMilestones() error {
	return importer.tracAccessor.GetMilestones(func(milestone *trac.Milestone) error {
		milestoneID, err := importer.giteaAccessor.GetMilestoneID(milestone.Name)
		if err != nil {
			return err
		}
		if milestoneID != -1 {
			log.Debug("milestone %s already exists - skipping...", milestone.Name)
			return nil
		}

		milestoneID, err = importer.giteaAccessor.AddMilestone(
			milestone.Name, milestone.Description, milestone.Completed != 0, milestone.Due, milestone.Completed)
		if err != nil {
			return err
		}

		log.Debug("added milestone (id %d) %s", milestoneID, milestone.Name)
		return nil
	})
}
