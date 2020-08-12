package trac

import "stevejefferson.co.uk/trac2gitea/log"

// GetComments retrieves all comments on a given Trac ticket, passing data from each one to the provided "handler" function.
func (accessor *Accessor) GetComments(ticketID int64, handlerFn func(ticketID int64, time int64, author string, comment string)) {
	rows, err := accessor.db.Query(`
		SELECT CAST(time*1e-6 AS int8) tim, COALESCE(author, '') author, COALESCE(newvalue, '') newval
			FROM ticket_change where ticket = $1 AND field = 'comment' AND trim(COALESCE(newvalue, ''), ' ') != ''
			ORDER BY time asc`, ticketID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var time int64
		var author, comment string
		if err := rows.Scan(&time, &author, &comment); err != nil {
			log.Fatal(err)
		}

		handlerFn(ticketID, time, author, comment)
	}
}
