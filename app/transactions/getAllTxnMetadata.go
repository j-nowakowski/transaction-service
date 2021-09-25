package transactions

// Aggregates all transaction records into one period "All".
func (e *EmailWriter) getAllTxnMetadata(txnRecords []transactionRecord) transactionMetadata {
	numOfCredits := 0
	sumCredits := 0.0
	numOfDebits := 0
	sumDebits := 0.0
	for _, txn := range txnRecords {
		if txn.Amount > 0 {
			numOfCredits++
			sumCredits += txn.Amount
		} else if txn.Amount < 0 {
			numOfDebits++
			sumDebits += txn.Amount
		}
	}
	var avgCredit float64
	if numOfCredits == 0 {
		avgCredit = 0
	} else {
		avgCredit = sumCredits / float64(numOfCredits)
	}
	var avgDebit float64
	if numOfDebits == 0 {
		avgDebit = 0
	} else {
		avgDebit = sumDebits / float64(numOfDebits)
	}

	return transactionMetadata{
		Period:            "All",
		NumOfTransactions: len(txnRecords),
		AvgCredit:         avgCredit,
		AvgDebit:          avgDebit,
		Sum:               sumCredits + sumDebits,
	}
}
