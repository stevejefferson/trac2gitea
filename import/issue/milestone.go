package issue

// ImportMilestones imports Trac milestones as Gitea milestones.
func (importer *Importer) ImportMilestones() {
	importer.tracAccessor.GetMilestones(func(name string, description string, due int64, completed int64) {
		importer.giteaAccessor.AddMilestone(name, description, completed != 0, due, completed)
	})
}
