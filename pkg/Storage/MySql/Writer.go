package MySql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func Save(db *sql.DB, tableName string, structPtr interface{}) error {
	// Get the type of the struct
	val := reflect.ValueOf(structPtr).Elem()
	typ := val.Type()

	// Prepare a slice to hold the field names and values for the query
	var fieldNames []string
	var placeholders []string
	var values []interface{}

	// Iterate over the struct fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Add the field name and a placeholder to the slices
		fieldNames = append(fieldNames, field.Name)
		placeholders = append(placeholders, "?")
		values = append(values, value.Interface())
	}

	// Construct the query string
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(fieldNames, ", "),
		strings.Join(placeholders, ", "),
	)

	// Execute the query
	_, err := db.Exec(query, values...)
	return err
}