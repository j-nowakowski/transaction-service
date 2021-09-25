package transactions

import "database/sql"

// EmailWriter is responsible for fetching transactions from the database, aggregating that
// data, and generating an output file.
type EmailWriter struct {
	DB interface {
		Query(query string, args ...interface{}) (*sql.Rows, error)
	}
}
