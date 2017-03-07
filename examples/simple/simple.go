package main

import (
	"github.com/eventials/gomigration"
	"github.com/jmoiron/sqlx"
)

func main() {
	storage := gomigration.NewPostgresStore("postgres://postgres:postgres@db/postgres?sslmode=disable")
	migrations := gomigration.NewMigrationsTree(storage, "example-app")

	db := migrations.Add("create languages table", func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			CREATE TABLE languages (
				id BIGSERIAL PRIMARY KEY,
				name TEXT NOT NULL
			);
		`)

		return err
	})

	db.Add("insert default languages", func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO languages (name) values ('assembly');
			INSERT INTO languages (name) values ('go-lang');
			INSERT INTO languages (name) values ('python');
			INSERT INTO languages (name) values ('java');
		`)

		return err
	})

	err := migrations.Execute()

	if err != nil {
		panic(err)
	}
}
