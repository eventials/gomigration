package main

import (
	"github.com/eventials/gomigration"
)

func main() {
	storage := gomigration.NewPostgresStore("postgres://postgres:postgres@db/postgres?sslmode=disable")
	storage.DropMigrationTable()
}
