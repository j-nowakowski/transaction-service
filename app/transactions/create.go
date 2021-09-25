package transactions

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

const (
	emailHeader    string = "Period,NumberOfTransactions,AverageCredit,AverageDebit,Sum\n"
	emailDelimiter string = ","
)

// Fetches and aggregates the transaction data and writes it to the input writer.
// dest cannot be nil. emailWriter's DB cannot be nil.
func (e *EmailWriter) Write(dest io.Writer) error {
	// Check dependencies
	if e.DB == nil || reflect.ValueOf(e.DB).IsNil() {
		return errors.New("emailWriter DB cannot be nil")
	}

	// Check input
	if dest == nil || reflect.ValueOf(dest).IsNil() {
		return errors.New("input dest cannot be nil")
	}

	// Fetch transactions from the DB
	txnRecords, err := e.fetchAllTransactions()
	if err != nil {
		return fmt.Errorf("in fetchAllTransactions: %w", err)
	}

	// Aggregate the transaction data
	allTxnMetadata := e.getAllTxnMetadata(txnRecords)
	monthlyTxnMetadata := e.getMonthlyTxnMetadata(txnRecords)

	// Write to the destination file
	dest.Write([]byte(emailHeader))
	dest.Write([]byte(fmt.Sprintf("%v%v%v%v%.2f%v%.2f%v%.2f\n",
		allTxnMetadata.Period, emailDelimiter,
		allTxnMetadata.NumOfTransactions, emailDelimiter,
		allTxnMetadata.AvgCredit, emailDelimiter,
		allTxnMetadata.AvgDebit, emailDelimiter,
		allTxnMetadata.Sum)))
	for _, monthData := range monthlyTxnMetadata {
		dest.Write([]byte(fmt.Sprintf("%v%v%v%v%.2f%v%.2f%v%.2f\n",
			monthData.Period, emailDelimiter,
			monthData.NumOfTransactions, emailDelimiter,
			monthData.AvgCredit, emailDelimiter,
			monthData.AvgDebit, emailDelimiter,
			monthData.Sum)))
	}

	return nil
}
