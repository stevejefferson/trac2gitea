package issue

import (
	"fmt"
	"log"
)

// ImportMilestones imports Trac milestones as Gitea milestones.
func (importer *Importer) ImportMilestones() {
	// NOTE: trac timestamps are to the microseconds, we just need seconds
	rows := importer.tracAccessor.Query(`
		SELECT COALESCE(name,''), CAST(due*1e-6 AS int8), CAST(completed*1e-6 AS int8), description
			FROM milestone UNION
			SELECT distinct(COALESCE(milestone,'')),0,0,''
				FROM ticket
				WHERE COALESCE(milestone,'') NOT IN ( select COALESCE(name,'') from milestone )`)

	for rows.Next() {
		var completed, due int64
		var nam, desc string
		if err := rows.Scan(&nam, &due, &completed, &desc); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Adding milestone", nam)
		importer.giteaAccessor.AddMilestone(nam, desc, completed != 0, due, completed)
	}
}
