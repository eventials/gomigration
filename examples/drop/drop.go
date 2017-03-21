package main

import (
	"github.com/eventials/gomigration/postgres"
)

func main() {
	storage := postgres.NewStorage("postgres://postgres:postgres@db/postgres?sslmode=disable")
	storage.DropMigrationTable()
}
