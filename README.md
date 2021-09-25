# Summary
This application demonstrates a transaction file parsing and aggregation service written in Go. 
Transactions are parsed from an input file (`stori-go/transactions.csv`), loaded into a MySQL database, and then 
fetched and aggregated. The aggregation results are stored in an output file (`stori-go/output/email.csv`).
If errors occur while parsing the input file, they are written to a log (`stori-go/out/output/errors.txt`).
Both the Go application and MySQL instance are ran in Docker containers. Note that both output files are
truncated when the application boots.

# Usage

Run the following command to build and run this application:
    `docker-compose down -v && docker-compose up --build`

The Go application will wait 15 seconds before executing as it waits for the MySQL instance to initialize.

# Parsing
The app reads from a file of delimited transaction data (`stori-go/transactions.csv`) and parses
transactions out of that file. It is expected that the first line in this file is a header, so 
this line is skipped. It is expected that each transaction has its own line in this file,
and that a transaction is made up of three parts, each delimited by a comma:
    `Id,Date,Amount`
    Example: `1,2021/03/14,-1.5`
Note that `Id` should be an integer, `Date` should be a date in the form of `YYYY/MM/DD`, and `Amount`
should be numeric. `Id` should be unique. If a single record fails to parse, other records will not be impacted.
Duplicate transactions are not inserted. Records which could not be parsed will have a descriptive error written to
the error log (`stori-go/out/output/errors.txt`).

# Storage
The transaction records which were successfully parsed from the input file are loaded into a MySQL database.
The database's name is `stori`, and it contains a table `transactions` with columns `id`, `date`, and `amount`.
Note that since the database itself is stored within a Docker container, it will be destroyed when the volume
it's in is destroyed.

# Aggregation:
After loading transactions, all transactions are fetched from the database and aggregated into an output file 
(`stori-go/output/email.csv`). The output file always starts with a header:
    `Period,NumberOfTransactions,AverageCredit,AverageDebit,Sum`
The first record after the header is always the `All` period, meaning all transactions are considered.
All records thereafter are aggregated by year and month. Those records are sorted by date ascending, and are
written in the form `YYYY/MM`. `AverageCredit`, `AverageDebit`, and `Sum` are rounded to 2 decimals.
    Example: `2021/03,3,2.25,-9.00,-4.50`

# Areas for Improvement
- Transaction amounts are expressed as float64s, which may be insufficiently precise. In a real
    application, a data type which guarantees the precision of all digits should be used instead.
- Unit tests need to be written.
- Some code duplication exists in `transactions.getMonthlyTxMetadata.go` and `transactions.getAllTxnMetadata.go` which
    could be dried up.
- Duplicate transaction Ids are currently not handled if they occur within the same batch.
