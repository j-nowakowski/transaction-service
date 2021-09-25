package transactions

import (
	"database/sql"
)

// Parser is responsible for parsing and inserting transaction records from an input.
type Parser struct {
	DB interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
	}
}
