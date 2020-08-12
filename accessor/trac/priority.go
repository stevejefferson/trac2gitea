package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetPriorityNames retrieves all priority names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *Accessor) GetPriorityNames(handlerFn func(string)) {
	rows, err := accessor.db.Query(`SELECT DISTINCT priority FROM ticket`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var priorityName string
		if err := rows.Scan(&priorityName); err != nil {
			log.Fatal(err)
		}

		handlerFn(priorityName)
	}
}
