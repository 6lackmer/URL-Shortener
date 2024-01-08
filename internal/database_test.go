package internal

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInitMySqlDB(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Set expectations
	mock.ExpectExec("CREATE DATABASE IF NOT EXISTS final_project").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("USE final_project").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS url_shortener").WillReturnResult(sqlmock.NewResult(0, 0))

	InitMySqlDB(db)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestConnectToMySqlDB(t *testing.T) {
	// Mock successful database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a mock database connection", err)
	}
	defer mockDB.Close()

	// Replace sqlOpen with a mock
	originalSqlOpen := sqlOpen
	sqlOpen = func(driverName string, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}
	defer func() { sqlOpen = originalSqlOpen }()

	// Call the function - expecting no error as the connection is successful
	if _, err := ConnectToMySqlDB(); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Now, simulate a connection error
	sqlOpen = func(driverName string, dataSourceName string) (*sql.DB, error) {
		return nil, sql.ErrConnDone
	}

	// Call the function again - expecting an error due to connection failure
	if _, err := ConnectToMySqlDB(); err == nil {
		t.Errorf("Expected an error, but got none")
	}
}
