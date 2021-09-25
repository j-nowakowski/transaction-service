package transactions

import (
	"fmt"
	"sort"
)

// Aggregates the transaction records by month. Records are recorded sorted by year, month ascending.
func (e *EmailWriter) getMonthlyTxnMetadata(txnRecords []transactionRecord) []transactionMetadata {
	// Intermediate data type used for aggregating
	type interMetadata struct {
		NumOfCredits      int
		SumCredits        float64
		NumOfDebits       int
		SumDebits         float64
		NumOfTransactions int
	}

	periodToMetadata := make(map[string]*interMetadata)
	for _, txn := range txnRecords {
		var month string
		if txn.Date.Month() <= 9 {
			month = "0" + fmt.Sprint(int(txn.Date.Month()))
		} else {
			month = fmt.Sprint(txn.Date.Month())
		}
		txnPeriod := fmt.Sprintf("%v/%v", txn.Date.Year(), month)
		_, ok := periodToMetadata[txnPeriod]
		if !ok {
			periodToMetadata[txnPeriod] = &interMetadata{}
		}
		metadata := periodToMetadata[txnPeriod]
		metadata.NumOfTransactions++
		if txn.Amount > 0 {
			metadata.NumOfCredits++
			metadata.SumCredits += txn.Amount
		} else if txn.Amount < 0 {
			metadata.NumOfDebits++
			metadata.SumDebits += txn.Amount
		}
	}

	result := make([]transactionMetadata, 0, len(periodToMetadata))
	for txnPeriod, metadata := range periodToMetadata {
		var avgCredit float64
		if metadata.NumOfCredits == 0 {
			avgCredit = 0
		} else {
			avgCredit = metadata.SumCredits / float64(metadata.NumOfCredits)
		}
		var avgDebit float64
		if metadata.NumOfDebits == 0 {
			avgDebit = 0
		} else {
			avgDebit = metadata.SumDebits / float64(metadata.NumOfDebits)
		}
		result = append(result, transactionMetadata{
			Period:            txnPeriod,
			NumOfTransactions: metadata.NumOfTransactions,
			AvgCredit:         avgCredit,
			AvgDebit:          avgDebit,
			Sum:               metadata.SumCredits + metadata.SumDebits,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Period < result[j].Period
	})

	return result
}
