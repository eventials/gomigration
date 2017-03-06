package migration

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Transaction interface {
	InsertId(appName string, id string) bool
	Commit()
	Rollback()

	GetTx() *sqlx.Tx
}

type Storage interface {
	CreateMigrationTable()
	DropMigrationTable()
	DeleteMigrations(appName string)

	GetTransaction() Transaction
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type TransactionImpl struct {
	tx *sqlx.Tx
}

type StoragePostgres struct {
	db *sqlx.DB
}

func (s *StoragePostgres) GetTransaction() Transaction {
	return &TransactionImpl{tx: s.db.MustBegin()}
}

func (s *StoragePostgres) CreateMigrationTable() {
	s.db.MustExec(`
	CREATE TABLE IF NOT EXISTS migration (
		id BIGSERIAL PRIMARY KEY NOT NULL,
		app_name TEXT NOT NULL,
		migration_id TEXT NOT NULL,
		date timestamp WITH TIME ZONE default (now() at TIME ZONE 'UTC')
	);
	CREATE UNIQUE INDEX IF NOT EXISTS migration_app_id ON migration (app_name, migration_id);
	`)
}

func (s *StoragePostgres) DropMigrationTable() {
	s.db.MustExec(`DROP TABLE IF EXISTS migration;`)
}

func (s *StoragePostgres) DeleteMigrations(appName string) {
	s.db.MustExec(`DELETE FROM migration WHERE app_name = $1`, appName)
}

func (t *TransactionImpl) InsertId(appName string, id string) bool {
	_, err := t.tx.Exec("INSERT INTO migration (app_name, migration_id) VALUES ($1, $2)", appName, id)
	return err == nil
}

func (t *TransactionImpl) Commit() {
	t.tx.Commit()
}

func (t *TransactionImpl) Rollback() {
	t.tx.Rollback()
}

func (t *TransactionImpl) GetTx() *sqlx.Tx {
	return t.tx
}

func NewPostgresStore(url string) Storage {
	db, err := sqlx.Connect("postgres", url)

	checkError(err)

	return &StoragePostgres{db: db}
}
