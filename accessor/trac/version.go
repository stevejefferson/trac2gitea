package trac

import "log"

// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
func (accessor *Accessor) GetVersionNames(handlerFn func(string)) {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(version,'') FROM ticket UNION SELECT COALESCE(name,'') FROM version`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			log.Fatal(err)
		}

		handlerFn(version)
	}
}
