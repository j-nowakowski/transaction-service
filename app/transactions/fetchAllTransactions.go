package transactions

import "fmt"

const fetchAllTransactionsSQL string = "SELECT `id`, `date`, `amount` FROM `transactions`"

// Fetches all transaction records from emailWriter's DB.
func (e *EmailWriter) fetchAllTransactions() ([]transactionRecord, error) {
	// Prepare query
	rows, err := e.DB.Query(fetchAllTransactionsSQL)
	if err != nil {
		return nil, fmt.Errorf("in DB.Query: %w", err)
	}
	defer rows.Close()

	// Read in records.
	txnRecords := []transactionRecord{}
	for rows.Next() {
		rec := transactionRecord{}
		err := rows.Scan(
			&rec.ID,
			&rec.Date,
			&rec.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("in rows.Next: %w", err)
		}
		txnRecords = append(txnRecords, rec)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("in rows.Err: %w", err)
	}
	return txnRecords, nil
}
