package main

import (
	"app/transactions"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	transactionsFileLoc string = "/tmp/transactions.csv"
	outputDirLoc        string = "/tmp/output"
	dbUsernameEnvVar    string = "MYSQL_USER"
	dbPasswordEnvVar    string = "MYSQL_PASSWORD"
	dbHostEnvVar        string = "MYSQL_HOST"
)

func main() {
	// Bootstrap the DB.
	db, err := initializeDbConn()
	if err != nil {
		panic(fmt.Errorf("in initializeDbConn: %w", err))
	}
	defer db.Close()
	err = initializeDbSchema(db)
	if err != nil {
		panic(fmt.Errorf("in initializeDbSchema: %w", err))
	}

	// Parse the transactions file.
	txnsFile, err := os.Open(transactionsFileLoc)
	if err != nil {
		panic(fmt.Errorf("while opening %v: %w", transactionsFileLoc, err))
	}
	defer txnsFile.Close()
	txnsFileReader := bufio.NewReader(txnsFile)
	parseLogLoc := fmt.Sprintf("%v/errors.txt", outputDirLoc)
	parseLog, err := os.OpenFile(parseLogLoc, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(fmt.Errorf("while opening %v: %w", parseLogLoc, err))
	}
	defer parseLog.Close()
	parser := &transactions.Parser{
		DB: db,
	}
	err = parser.ParseFile(txnsFileReader, parseLog)
	if err != nil {
		panic(fmt.Errorf("in parser.Parse: %w", err))
	}

	// Create the email file.
	emailFileLoc := fmt.Sprintf("%v/email.csv", outputDirLoc)
	emailFile, err := os.OpenFile(emailFileLoc, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		panic(fmt.Errorf("while opening %v: %w", emailFileLoc, err))
	}
	defer emailFile.Close()
	emailWriter := &transactions.EmailWriter{
		DB: db,
	}
	err = emailWriter.Write(emailFile)
	if err != nil {
		panic(fmt.Errorf("in emailWriter.Create: %w", err))
	}
}

// Establishes a connection to the database.
func initializeDbConn() (*sql.DB, error) {
	username := os.Getenv(dbUsernameEnvVar)
	password := os.Getenv(dbPasswordEnvVar)
	host := os.Getenv(dbHostEnvVar)
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/?parseTime=true", username, password, host))
	time.Sleep(15 * time.Second) // This is to wait for MySQL to initialize.
	return db, err
}

// Creates the schema for database.
//
// For the sake of this example project, no volume is being used to store the DB outside of the container, so the DB
// will be destroyed every time the container resets, hence the need for the schema to be initialized here.
// If this function existed in a real application, it would need to be written idempotently.
func initializeDbSchema(db *sql.DB) error {
	_, err := db.Exec("CREATE DATABASE myDB;")
	if err != nil {
		return fmt.Errorf("while creating database: %v", err)
	}

	_, err = db.Exec("USE myDB;")
	if err != nil {
		return fmt.Errorf("while using database: %v", err)
	}

	_, err = db.Exec("CREATE TABLE transactions(" +
		"`id` INT NOT NULL, " +
		"`date` DATE NOT NULL, " +
		"`amount` DECIMAL(11, 2) NOT NULL, " +
		"PRIMARY KEY (`id`)" +
		") ENGINE=InnoDB")
	if err != nil {
		return fmt.Errorf("while creating transactions table: %v", err)
	}

	return nil
}
