package service

import (
	"database/sql"
	"fmt"
	"log"
)

// Database wraps our SQL database. Defining our own type allows us to define helper functions on the Database.
type Database struct {
	DB *sql.DB
}

func (db *Database) Close() {
	db.DB.Close()
}

// ===== TRANSACTIONS ==================================================================================================

// Transaction wraps a SQL transaction. Defining our own type allows functions to be defined on the Transaction.
type Transaction struct {
	*sql.Tx
	db *Database
}

type TransactionFunc func(*Transaction)

func (db *Database) begin() (*Transaction, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Transaction{tx, db}, nil
}

// Read begins a read-only transaction and passes it to the given function. The transaction will be rolled back after
// the function returns. Any panics will be handled, and returned as an error.
func (db *Database) Read(reader TransactionFunc) (err error) {
	tx, err := db.begin()
	if err != nil {
		return err
	}

	// A read should always rollback the transaction
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			// Ignore errors rolling back
			log.Println(rollbackErr.Error())
		}
	}()

	// recover any panics during the transaction, and return it as an error to the caller
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Database.Read: %v", r)
			}
		}
	}()

	// Mark the transaction as read only
	_, err = tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ, READ ONLY")
	if err != nil {
		panic(fmt.Errorf("Unable to mark transaction read-only"))
	}

	reader(tx) // Code in this function can panic
	return err
}

// Write begins a transaction and passes it to the given function. The transaction will be committed when the function
// returns. If the function panics, the transaction is rolled back, and the error provided to panic is returned.
func (db *Database) Write(writer TransactionFunc) (err error) {
	tx, err := db.begin()
	if err != nil {
		return err
	}

	didPanic := false

	// write operations commit or rollback the transaction
	defer func() {
		if didPanic {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Ignore errors rolling back
				log.Println(rollbackErr.Error())
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = commitErr
			}
		}
	}()

	// recover any panics during the transaction, and return it as an error to the caller
	defer func() {
		if r := recover(); r != nil {
			didPanic = true

			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("Database.Write: %v", r)
			}
		}
	}()

	writer(tx) // If the function panics, the transaction will be rolled back

	return err
}
