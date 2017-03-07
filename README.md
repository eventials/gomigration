# migration
Migration lib for Golang projects.

Any added migration will be performed within an individual transaction, if it returns an error, they will be rolled back.
All "trunk" migrations will be perfomed within an go routine, and his method will be sequentially executed.

## Preparing the environment

Install Docker:

* [Install steps](https://docs.docker.com/engine/installation/)

## Compiling and running the app

To compile and run the example, just run:

```
docker-compose run app go run ./examples/simple/simple.go
```

## Example

Every node returned from Add() will be execute in a go routine.

```go
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
```

If you need to clear all migrations form the current app, just do:

```go
import (
	"github.com/eventials/gomigration"
)
func main() {
	storage := gomigration.NewPostgresStore("postgres://postgres:postgres@db/postgres?sslmode=disable")
	migrations := gomigration.NewMigrationsTree(storage, "example-app")

	migrations.Clear()
}

```

Another clear example:

```go
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
```

If you want to drop the migration control table, just do:

```go
import (
	"github.com/eventials/gomigration"
)

func main() {
	storage := gomigration.NewPostgresStore("postgres://postgres:postgres@db/postgres?sslmode=disable")
	storage.DropMigrationTable()
}
```
