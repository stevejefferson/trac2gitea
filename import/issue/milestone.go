// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package issue

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
	"github.com/stevejefferson/trac2gitea/log"
)

// ImportMilestones imports Trac milestones as Gitea milestones.
func (importer *Importer) ImportMilestones() error {
	return importer.tracAccessor.GetMilestones(func(tracMilestone *trac.Milestone) error {
		milestoneID, err := importer.giteaAccessor.GetMilestoneID(tracMilestone.Name)
		if err != nil {
			return err
		}
		if milestoneID != -1 {
			log.Debug("milestone %s already exists - skipping...", tracMilestone.Name)
			return nil
		}

		giteaMilestone := gitea.Milestone{
			Name:        tracMilestone.Name,
			Description: tracMilestone.Description,
			Closed:      tracMilestone.Completed != 0,
			DueTime:     tracMilestone.Due,
			ClosedTime:  tracMilestone.Completed}

		milestoneID, err = importer.giteaAccessor.AddMilestone(&giteaMilestone)
		if err != nil {
			return err
		}

		log.Debug("added milestone (id %d) %s", milestoneID, tracMilestone.Name)
		return nil
	})
}
