package sqlite3

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	storage *Storage
)

func TestMain(m *testing.M) {
	storage = NewStorage(".test.db")
	os.Exit(m.Run())
}

func ExistTable(table string) bool {
	sql := fmt.Sprintf(`SELECT COUNT(tbl_name) > 0 FROM sqlite_master WHERE tbl_name = '%s'`, table)
	exists := false
	err := storage.db.QueryRow(sql).Scan(&exists)
	if err != nil {
		panic(err)
	}

	return exists
}

func ExistsId(appName string, id string) bool {
	exists := false
	err := storage.db.QueryRow(
		`SELECT EXISTS (
                         SELECT 1
                           FROM migration
                          WHERE app_name = $1 AND migration_id = $2)`, appName, id).Scan(&exists)

	if err != nil {
		panic(err)
	}

	return exists
}

func TestStorageCreateMigrationTable(t *testing.T) {
	storage.CreateMigrationTable()
	assert.True(t, ExistTable("migration"))
}

func TestStorageDropMigrationTable(t *testing.T) {
	storage.CreateMigrationTable()
	storage.DropMigrationTable()
	assert.False(t, ExistTable("migration"))
}

func TestStorageInsertMigrationTable(t *testing.T) {
	storage.DropMigrationTable()
	storage.CreateMigrationTable()

	transaction := storage.GetTransaction()
	defer transaction.Rollback()

	transaction.InsertId("teste", "teste")

	exists := false
	err := transaction.GetTx().QueryRow(
		`SELECT EXISTS (SELECT 1 FROM migration WHERE app_name = $1 AND migration_id = $2)`, "teste", "teste").Scan(&exists)

	if err != nil {
		panic(err)
	}

	assert.True(t, exists)
}

func TestStorageClearMigrationTable(t *testing.T) {
	storage.DropMigrationTable()
	storage.CreateMigrationTable()

	tx := storage.GetTransaction()
	tx.InsertId("teste", "teste")
	tx.InsertId("teste2", "teste")
	tx.InsertId("teste3", "teste")

	tx.Commit()

	storage.DeleteMigrations("teste")

	assert.False(t, ExistsId("teste", "teste"))
	assert.True(t, ExistsId("teste2", "teste"))
	assert.True(t, ExistsId("teste3", "teste"))
}

func TestStorageMultiInstanceExistsId(t *testing.T) {
	wg := sync.WaitGroup{}

	storage.DropMigrationTable()
	storage.CreateMigrationTable()

	tx := storage.GetTransaction()
	tx2 := storage.GetTransaction()

	tx.InsertId("teste", "teste")

	inserted := false
	wg.Add(1)
	go func() {
		inserted = tx2.InsertId("teste", "teste")
		wg.Done()
	}()

	tx.Commit()
	tx2.Rollback()
	wg.Wait()

	assert.False(t, inserted)
}
