package transactions

import "strings"

const baseSQL string = "INSERT INTO `transactions`(`id`, `date`, `amount`) " +
	"SELECT * FROM (:transactions) `tx` " +
	"WHERE NOT EXISTS (" +
	"	SELECT 1 FROM `transactions` `tx2` " +
	"	WHERE `tx2`.`id` = `tx`.`id` " +
	")"

const recordSQL string = "SELECT ? AS `id`, ? AS `date`, ? AS `amount`"

// Inserts the input transaction records to the parser's DB.
// Does nothing if input slice is empty.
func (p *Parser) insertTransactions(txnRecords []transactionRecord) error {
	// Exit early?
	if len(txnRecords) < 1 {
		return nil
	}

	// Construct the SQL
	innerSQL := strings.Repeat(recordSQL+" UNION ALL ", len(txnRecords)-1) + recordSQL
	fullSQL := strings.Replace(baseSQL, ":transactions", innerSQL, 1)

	// Flatten the values
	values := make([]interface{}, 0, len(txnRecords)*3)
	for _, txnRec := range txnRecords {
		values = append(values, txnRec.ID, txnRec.Date, txnRec.Amount)
	}

	// Execute
	_, err := p.DB.Exec(fullSQL, values...)
	return err
}
