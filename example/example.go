package main

import (
	"errors"
	"fmt"
	mg "github.com/eventials/gomigration"
	"github.com/jmoiron/sqlx"
)

func ExecuteAlterTableTestTable(tx *sqlx.Tx) error {
	tx.MustExec(`ALTER TABLE test ADD COLUMN test int;`)

	value := 1

	_, err := tx.NamedExec(`UPDATE test SET test = :value`, map[string]interface{}{
		"value": value})

	return err
}

func ExecuteInsertCourses(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		INSERT INTO courses (name) values ('assembly');
		INSERT INTO courses (name) values ('go-lang');
		INSERT INTO courses (name) values ('python');
		INSERT INTO courses (name) values ('java');
	`)

	return err
}

func ExecuteDeleteInvalidEntity(tx *sqlx.Tx) error {
	return nil
}

func ExecuteAnotherChangeInTestTable(tx *sqlx.Tx) error {
	return errors.New("Failed to execute ExecuteAnotherChangeInTestTable")
}

func ExecuteCreateDatabase(tx *sqlx.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE test (
			id BIGSERIAL PRIMARY KEY
		);
		CREATE TABLE courses (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);
	`)
	return err
}

func DropTables(tx *sqlx.Tx) error {
	tx.MustExec(`
		DROP TABLE test;
		DROP TABLE courses;
	`)

	return nil
}

func main() {
	storage := mg.NewPostgresStore("postgres://postgres:postgres@db/postgres?sslmode=disable")
	migrations := mg.NewMigrationsTree(storage, "app-name")

	db := migrations.Add("26d766a2-a593-48e4-b53a-4150d8113cc7", ExecuteCreateDatabase)

	testTable := db.Add("4666a9e5-aa42-4e0c-888f-60c4d4c0cb15", ExecuteAlterTableTestTable)
	testTableWithTestColumn := testTable.Add("691e6c54-00c2-457d-84f8-f15aeba05ded", ExecuteAnotherChangeInTestTable)
	testTableWithTestColumn.Add("random keyword to identify migration", ExecuteDeleteInvalidEntity)

	db.Add("d493d3d9-18c1-4519-8053-b9832b33b619", ExecuteInsertCourses)

	db.Add("f0653f39-15e7-4517-bd3d-2b180825fb29", ExecuteDeleteInvalidEntity)

	err := migrations.Execute()

	if err != nil {
		fmt.Println(err)
	}
}
