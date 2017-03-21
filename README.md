# migration

Database migration library for Golang projects.

Any added migration will be performed within an individual transaction, if it returns an error, they will be rolled back.
All "trunk" migrations will be perfomed within an go routine, and his method will be sequentially executed.

## Database Support

- PostgreSQL

## Example

Every node returned from Add() will be execute in a go routine.

```go
import (
	"database/sql"

	"github.com/eventials/gomigration"
	"github.com/eventials/gomigration/postgres"
)

func main() {
	storage := postgres.NewStorage("postgres://postgres:postgres@db/postgres?sslmode=disable")
	migrations := gomigration.NewMigrationsTree(storage, "example-app")

	db := migrations.Add("create languages table", func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			CREATE TABLE languages (
				id BIGSERIAL PRIMARY KEY,
				name TEXT NOT NULL
			);
		`)

		return err
	})

	db.Add("insert default languages", func(tx *sql.Tx) error {
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

More examples in `examples` directory.
