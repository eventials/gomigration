package gomigration

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
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

func (m *MockedTransaction) GetTx() *sql.Tx {
	args := m.Called()
	tx := args.Get(0)

	if tx == nil {
		return nil
	}

	return tx.(*sql.Tx)
}
