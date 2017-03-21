package postgres

import (
	"database/sql"

	"github.com/eventials/gomigration"
	_ "github.com/lib/pq"
)

type Transaction struct {
	tx *sql.Tx
}

type Storage struct {
	db *sql.DB
}

func (s *Storage) GetTransaction() gomigration.Transaction {
	tx, err := s.db.Begin()

	if err != nil {
		panic(err)
	}

	return &Transaction{tx: tx}
}

func (s *Storage) CreateMigrationTable() {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS migration
         (
             id BIGSERIAL PRIMARY KEY NOT NULL,
             app_name TEXT NOT NULL,
             migration_id TEXT NOT NULL,
             date timestamp WITH TIME ZONE default (now() at TIME ZONE 'UTC')
         );

         CREATE UNIQUE INDEX IF NOT EXISTS migration_app_id ON migration (app_name, migration_id);`)

	if err != nil {
		panic(err)
	}
}

func (s *Storage) DropMigrationTable() {
	_, err := s.db.Exec(`DROP TABLE IF EXISTS migration;`)

	if err != nil {
		panic(err)
	}
}

func (s *Storage) DeleteMigrations(appName string) {
	_, err := s.db.Exec(`DELETE FROM migration WHERE app_name = $1`, appName)

	if err != nil {
		panic(err)
	}
}

func (t *Transaction) InsertId(appName string, id string) bool {
	_, err := t.tx.Exec("INSERT INTO migration (app_name, migration_id) VALUES ($1, $2)", appName, id)
	return err == nil
}

func (t *Transaction) Commit() {
	t.tx.Commit()
}

func (t *Transaction) Rollback() {
	t.tx.Rollback()
}

func (t *Transaction) GetTx() *sql.Tx {
	return t.tx
}

func NewStorage(url string) *Storage {
	db, err := sql.Open("postgres", url)

	if err != nil {
		panic(err)
	}

	return &Storage{db: db}
}
