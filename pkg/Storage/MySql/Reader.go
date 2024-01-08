package MySql

import (
	"database/sql"
	"fmt"
	"reflect"
)

func GetAll(db *sql.DB, tableName string, slicePtr interface{}) error {
	// Check that slicePtr is a pointer to a slice
	sliceVal := reflect.ValueOf(slicePtr)
	if sliceVal.Kind() != reflect.Ptr || sliceVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("slicePtr must be a pointer to a slice")
	}

	// Check that the slice element is a struct
	elementType := sliceVal.Elem().Type().Elem()
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		element := reflect.New(elementType).Elem()
		fieldValues := make([]interface{}, element.NumField())
		for i := 0; i < element.NumField(); i++ {
			fieldValues[i] = element.Field(i).Addr().Interface()
		}

		if err := rows.Scan(fieldValues...); err != nil {
			return err
		}

		sliceVal.Elem().Set(reflect.Append(sliceVal.Elem(), element))
	}

	return rows.Err()
}

func GetByWhere(db *sql.DB, tableName string, whereClause string, args []interface{}, objPtr interface{}) error {
    objVal := reflect.ValueOf(objPtr)
    if objVal.Kind() != reflect.Ptr || objVal.Elem().Kind() != reflect.Struct {
        return fmt.Errorf("objPtr must be a pointer to a struct")
    }

    query := fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, whereClause)
    row := db.QueryRow(query, args...) // Pass the arguments to QueryRow

    fieldValues := make([]interface{}, objVal.Elem().NumField())
    for i := 0; i < objVal.Elem().NumField(); i++ {
        fieldValues[i] = objVal.Elem().Field(i).Addr().Interface()
    }

    err := row.Scan(fieldValues...)
    if err != nil {
        return fmt.Errorf("error scanning row: %v", err)
    }

    return nil
}