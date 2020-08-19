package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetPriorityNames retrieves all priority names used in Trac tickets, passing each one to the provided "handler" function.
func (accessor *DefaultAccessor) GetPriorityNames(handlerFn func(priorityName string) error) error {
	rows, err := accessor.db.Query(`SELECT DISTINCT priority FROM ticket`)
	if err != nil {
		log.Error(err)
		return err
	}

	for rows.Next() {
		var priorityName string
		if err := rows.Scan(&priorityName); err != nil {
			log.Error(err)
			return err
		}

		err = handlerFn(priorityName)
		if err != nil {
			return err
		}
	}

	return nil
}
