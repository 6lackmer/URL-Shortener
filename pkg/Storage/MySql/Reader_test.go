package MySql
import (
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
)

type TestStruct struct {
    ID    int
    Name  string
    Value string
}

func TestGetAllWithSqlmock(t *testing.T) {
    db, mock, err := sqlmock.New() 
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    columns := []string{"id", "name", "value"}
    mock.ExpectQuery("^SELECT \\* FROM test_table$").WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "testName", "testValue"))

    var results []TestStruct
    err = GetAll(db, "test_table", &results)
    if err != nil {
        t.Errorf("Error in GetAll: %v", err)
    }

    if len(results) == 0 {
        t.Error("Expected non-empty result set")
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}

func TestGetByWhereWithSqlmock(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    columns := []string{"id", "name", "value"}
    mock.ExpectQuery("^SELECT \\* FROM test_table WHERE id = \\?$").
        WithArgs(1).
        WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "testName", "testValue"))

    var result TestStruct
    err = GetByWhere(db, "test_table", "id = ?", []interface{}{1}, &result)
    if err != nil {
        t.Errorf("Error in GetByWhere: %v", err)
    }

    if result.ID != 1 {
        t.Errorf("Expected ID to be 1, got %d", result.ID)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}