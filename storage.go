package gomigration

import (
	"database/sql"
)

type Transaction interface {
	InsertId(appName string, id string) bool
	Commit()
	Rollback()
	GetTx() *sql.Tx
}

type Storage interface {
	CreateMigrationTable()
	DropMigrationTable()
	DeleteMigrations(appName string)
	GetTransaction() Transaction
}
