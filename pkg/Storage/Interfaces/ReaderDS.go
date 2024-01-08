package StorageInterfaces

import (
	"database/sql"
)

type ReaderDS interface {
	GetAll(db *sql.DB, tableName string, slicePtr interface{}) error
	GetByWhere(db *sql.DB, tableName string, whereClause string, objPtr interface{}) error
}