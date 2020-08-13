package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetSeverityNames retrieves all severity names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetSeverityNames(handlerFn func(severityName string)) {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(severity,'') FROM ticket`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var severityName string
		if err := rows.Scan(&severityName); err != nil {
			log.Fatal(err)
		}

		handlerFn(severityName)
	}
}
