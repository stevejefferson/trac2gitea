package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetVersionNames retrieves all version names used in Trac, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetVersionNames(handlerFn func(version string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT COALESCE(version,'') FROM ticket UNION SELECT COALESCE(name,'') FROM version`)
	if err != nil {
		log.Error(err)
		return err
	}

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			log.Error(err)
			return err
		}

		err = handlerFn(version)
		if err != nil {
			return err
		}
	}

	return nil
}
