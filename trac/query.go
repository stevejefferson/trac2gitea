package trac

import (
	"database/sql"
	"log"
)

func (accessor *Accessor) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := accessor.db.Query(query, args)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}
