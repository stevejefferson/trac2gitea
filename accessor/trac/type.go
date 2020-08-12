package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetTypeNames retrieves all type names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *Accessor) GetTypeNames(handlerFn func(string)) {
	rows, err := accessor.db.Query(`SELECT DISTINCT type FROM ticket`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var typ string
		if err := rows.Scan(&typ); err != nil {
			log.Fatal(err)
		}

		handlerFn(typ)
	}
}
