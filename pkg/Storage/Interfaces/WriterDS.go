package StorageInterfaces

import (
	"database/sql"
)

type WriterDS interface {
	Save(db *sql.DB, tableName string, structPtr interface{}) error 
}