package issue

import (
	"log"
)

func (importer *Importer) importLabels(tracQuery string, labelPrefix string, labelColor string) {
	rows := importer.tracAccessor.Query(tracQuery)
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			log.Fatal(err)
		}
		lbl := labelPrefix + " / " + val
		importer.giteaAccessor.AddLabel(lbl, labelColor)
	}
}

// ImportComponents imports Trac components as Gitea labels.
func (importer *Importer) ImportComponents() {
	importer.importLabels(`SELECT name FROM component`, "Component", "#fbca04")
}

// ImportPriorities imports Trac priorities as Gitea labels.
func (importer *Importer) ImportPriorities() {
	importer.importLabels(`SELECT DISTINCT priority FROM ticket`, "Priority", "#207de5")
}

// ImportSeverities imports Trac severities as Gitea labels.
func (importer *Importer) ImportSeverities() {
	importer.importLabels(`SELECT DISTINCT COALESCE(severity,'') FROM ticket`, "Severity", "#eb6420")
}

// ImportVersions imports Trac versions as Gitea labels.
func (importer *Importer) ImportVersions() {
	importer.importLabels(`SELECT DISTINCT COALESCE(version,'') FROM ticket UNION
                        SELECT COALESCE(name,'') FROM version`, "Version", "#009800")
}

// ImportTypes imports Trac types as Gitea labels.
func (importer *Importer) ImportTypes() {
	importer.importLabels(`SELECT DISTINCT type FROM ticket`, "Type", "#e11d21")
}

// ImportResolutions imports Trac resolutions as Gitea labels.
func (importer *Importer) ImportResolutions() {
	importer.importLabels(`SELECT DISTINCT resolution FROM ticket WHERE trim(resolution) != ''`, "Resolution", "#9e9e9e")
}
