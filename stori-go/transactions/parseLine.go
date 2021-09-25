package transactions

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	delimiter           string = ","
	expectedNumOfTokens int    = 3
	dateFormat          string = "2006/01/02"
)

// Parses a string as a transaction record. Input string should be in the form Id,Date,Amount where
// Id is an int, Date is a date in the form YYYY/MM/DD, and Amount is numeric. Id should be unique.
// An error will be returned if any of the above rules are broken.
func (p *Parser) parseLine(line string) (transactionRecord, error) {
	tokens := strings.Split(line, delimiter)
	if len(tokens) != expectedNumOfTokens {
		return transactionRecord{}, fmt.Errorf("while parsing `%v`, expected %v tokens, instead counted %v",
			line, expectedNumOfTokens, len(tokens))
	}

	// Parse ID (token #0)
	id, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return transactionRecord{}, fmt.Errorf("while parsing id out of `%v`: %w", line, err)
	}

	// Parse Date (token #1)
	date, err := time.Parse(dateFormat, tokens[1])
	if err != nil {
		return transactionRecord{}, fmt.Errorf("while parsing date out of `%v`: %w", line, err)
	}

	// Parse Amount (token #2)
	amount, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return transactionRecord{}, fmt.Errorf("while parsing amount out of `%v`: %w", line, err)
	}

	return transactionRecord{
		ID:     int(id),
		Date:   date,
		Amount: amount,
	}, nil
}
