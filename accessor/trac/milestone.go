package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetMilestones retrieves all Trac milestones, passing data from each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetMilestones(handlerFn func(name string, description string, due int64, completed int64)) {
	// NOTE: trac timestamps are to the microseconds, we just need seconds
	rows, err := accessor.db.Query(`
		SELECT COALESCE(name,''), description, CAST(due*1e-6 AS int8), CAST(completed*1e-6 AS int8)
			FROM milestone UNION
			SELECT distinct(COALESCE(milestone,'')),'',0,0
				FROM ticket
				WHERE COALESCE(milestone,'') NOT IN ( select COALESCE(name,'') from milestone )`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var completed, due int64
		var name, description string
		if err := rows.Scan(&name, &description, &due, &completed); err != nil {
			log.Fatal(err)
		}

		handlerFn(name, description, due, completed)
	}
}
