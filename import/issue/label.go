package issue

// ImportComponents imports Trac components as Gitea labels.
func (importer *Importer) ImportComponents() {
	importer.tracAccessor.GetComponentNames(func(cmptName string) {
		label := "Component/" + cmptName
		importer.giteaAccessor.AddLabel(label, "#fbca04")
	})
}

// ImportPriorities imports Trac priorities as Gitea labels.
func (importer *Importer) ImportPriorities() {
	importer.tracAccessor.GetPriorityNames(func(priorName string) {
		label := "Priority/" + priorName
		importer.giteaAccessor.AddLabel(label, "#207de5")
	})
}

// ImportSeverities imports Trac severities as Gitea labels.
func (importer *Importer) ImportSeverities() {
	importer.tracAccessor.GetSeverityNames(func(sevName string) {
		label := "Severity/" + sevName
		importer.giteaAccessor.AddLabel(label, "#eb6420")
	})
}

// ImportVersions imports Trac versions as Gitea labels.
func (importer *Importer) ImportVersions() {
	importer.tracAccessor.GetVersionNames(func(verName string) {
		label := "Version/" + verName
		importer.giteaAccessor.AddLabel(label, "#009800")
	})
}

// ImportTypes imports Trac types as Gitea labels.
func (importer *Importer) ImportTypes() {
	importer.tracAccessor.GetTypeNames(func(typeName string) {
		label := "Type/" + typeName
		importer.giteaAccessor.AddLabel(label, "#e11d21")
	})
}

// ImportResolutions imports Trac resolutions as Gitea labels.
func (importer *Importer) ImportResolutions() {
	importer.tracAccessor.GetResolutionNames(func(typeName string) {
		label := "Resolution/" + typeName
		importer.giteaAccessor.AddLabel(label, "#9e9e9e")
	})
}
