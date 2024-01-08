package MySql

import (
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
)

type TestEntity struct {
    ID    int
    Name  string
    Value string
}

func TestSave(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    entity := TestEntity{
        ID:    1,
        Name:  "Test Name",
        Value: "Test Value",
    }

    mock.ExpectExec("^INSERT INTO test_table \\(ID, Name, Value\\) VALUES \\(\\?, \\?, \\?\\)$").
        WithArgs(entity.ID, entity.Name, entity.Value).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = Save(db, "test_table", &entity)
    if err != nil {
        t.Errorf("Error in Save: %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}
