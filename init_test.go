package migration

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

var (
	storage Storage
	db      *sqlx.DB
)

type MockedStorage struct {
	mock.Mock
}

func (m *MockedStorage) CreateMigrationTable() {
	m.Called()
}

func (m *MockedStorage) DropMigrationTable() {
	m.Called()
}

func (m *MockedStorage) DeleteMigrations(appName string) {
	m.Called(appName)
}

func (m *MockedStorage) GetTransaction() Transaction {
	args := m.Called()
	return args.Get(0).(Transaction)
}

type MockedTransaction struct {
	mock.Mock
}

func (m *MockedTransaction) InsertId(appName string, id string) bool {
	args := m.Called(appName, id)
	return args.Bool(0)
}

func (m *MockedTransaction) Commit() {
	m.Called()
}

func (m *MockedTransaction) Rollback() {
	m.Called()
}

func (m *MockedTransaction) GetTx() *sqlx.Tx {
	args := m.Called()
	tx := args.Get(0)

	if tx == nil {
		return nil
	}

	return tx.(*sqlx.Tx)
}

func setupTestStorage() {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://postgres:postgres@db/postgres?sslmode=disable")

	if err != nil {
		panic(err)
	}

	storage = &StoragePostgres{db: db}
}

func setupTest() {

}

func TestMain(m *testing.M) {
	setupTest()
	setupTestStorage()
	os.Exit(m.Run())
}
