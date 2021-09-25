package transactions

import "time"

type transactionRecord struct {
	ID     int
	Date   time.Time
	Amount float64
}
