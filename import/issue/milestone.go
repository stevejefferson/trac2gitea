package issue

import "stevejefferson.co.uk/trac2gitea/log"

// ImportMilestones imports Trac milestones as Gitea milestones.
func (importer *Importer) ImportMilestones() {
	importer.tracAccessor.GetMilestones(func(name string, description string, due int64, completed int64) {
		if importer.giteaAccessor.GetMilestoneID(name) != -1 {
			log.Debugf("milestone %s already exists - skipping...\n", name)
			return
		}
		milestoneID := importer.giteaAccessor.AddMilestone(name, description, completed != 0, due, completed)
		log.Debugf("Added milestone (id %d) %s\n", milestoneID, name)
	})
}
