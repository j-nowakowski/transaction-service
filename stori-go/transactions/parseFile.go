package transactions

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

const batchSize int = 10

type fileReadLiner interface {
	ReadLine() (line []byte, isPrefix bool, err error)
}

// Parses and inserts the transactions from the input file. If there is an error with a transaction
// record in the input file, the error is logged to log and that record is skipped. Parsed transactions
// are inserted into the database in batches. Both file and log each cannot be nil. The parser's DB cannot be nil.
func (p *Parser) ParseFile(file fileReadLiner, log io.Writer) error {
	// Check dependencies
	if p.DB == nil || reflect.ValueOf(p.DB).IsNil() {
		return errors.New("parser DB cannot be nil")
	}

	// Check input
	if file == nil || reflect.ValueOf(file).IsNil() {
		return errors.New("input file cannot be nil")
	}
	if log == nil || reflect.ValueOf(log).IsNil() {
		return errors.New("input log cannot be nil")
	}

	// Process each line
	transactionBuffer := make([]transactionRecord, 0, batchSize)
	i := 0
	for {
		i++
		lineBytes, _, err := file.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Write([]byte(fmt.Sprintf("error reading line #%v: %v\n", i, err.Error())))
			continue
		}
		if i == 1 {
			continue // Skip first line
		}
		line := string(lineBytes)
		if line == "" {
			continue // Skip empty lines
		}
		txnRecord, err := p.parseLine(line)
		if err != nil {
			log.Write([]byte(fmt.Sprintf("error parsing line #%v, %v\n", i, err.Error())))
			continue
		}
		transactionBuffer = append(transactionBuffer, txnRecord)
		if len(transactionBuffer) == cap(transactionBuffer) {
			err = p.insertTransactions(transactionBuffer)
			if err != nil {
				log.Write([]byte(fmt.Sprintf("error inserting transactions near line #%v, %v\n", i, err)))
				continue
			}
		}
	}
	// Insert leftover records
	err := p.insertTransactions(transactionBuffer)
	if err != nil {
		log.Write([]byte(fmt.Sprintf("error inserting transactions near line #%v, %v\n", i, err)))
	}

	return nil
}
