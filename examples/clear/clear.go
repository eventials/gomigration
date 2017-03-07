package main

import (
	"github.com/eventials/gomigration"
	"github.com/jmoiron/sqlx"
)

func main() {
	storage := gomigration.NewPostgresStore("postgres://postgres:postgres@db/postgres?sslmode=disable")
	migrations := gomigration.NewMigrationsTree(storage, "example-app")

	trunk := migrations.Add("drop tables", func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`DROP TABLE IF EXISTS languages;`)

		return err
	})

	languages := trunk.Add("create tables", func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			CREATE TABLE languages (
				id BIGSERIAL PRIMARY KEY,
				name TEXT NOT NULL
			);
		`)

		return err
	})

	languages.Add("insert default languages", func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO languages (name) values ('assembly');
			INSERT INTO languages (name) values ('go-lang');
			INSERT INTO languages (name) values ('python');
			INSERT INTO languages (name) values ('java');
		`)

		return err
	})

	languages.Add("added type column to langue table", func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`ALTER TABLE languages ADD COLUMN type TEXT;`)

		return err
	})

	// all migrations from "example-app" will be deleted
	migrations.Clear()
	err := migrations.Execute()

	if err != nil {
		panic(err)
	}
}
