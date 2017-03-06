package migration

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateMigrationTable(t *testing.T) {
	storage := new(MockedStorage)

	storage.On("CreateMigrationTable")

	migration := NewMigration(storage, "test-app")

	assert.NotNil(t, migration)
	storage.AssertExpectations(t)
}

func TestClearMigrationTable(t *testing.T) {
	storage := new(MockedStorage)

	storage.On("CreateMigrationTable")
	storage.On("DeleteMigrations", "clear-app")

	migration := NewMigration(storage, "clear-app")
	migration.Clear()

	storage.AssertExpectations(t)
}

func TestDuplicatedIds(t *testing.T) {
	storage := new(MockedStorage)

	storage.On("CreateMigrationTable")

	migration := NewMigration(storage, "test-app")

	migration.Add("1", nil)
	migration.Add("1", nil)

	assert.NotNil(t, migration.Execute())

	storage.AssertExpectations(t)
}

func TestInsertId(t *testing.T) {
	storage := new(MockedStorage)
	transaction := new(MockedTransaction)

	transaction.On("InsertId", "test-app", "1").Return(true)
	transaction.On("Commit")
	transaction.On("GetTx").Return(nil)

	storage.On("CreateMigrationTable")
	storage.On("GetTransaction").Return(transaction)

	migration := NewMigration(storage, "test-app")

	assert.NotNil(t, migration)

	called := false

	migration.Add("1", func(tx *sqlx.Tx) error {
		called = true
		return nil
	})

	migration.Execute()

	assert.True(t, called)
	transaction.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestInsertIdFail(t *testing.T) {
	storage := new(MockedStorage)
	transaction := new(MockedTransaction)

	transaction.On("InsertId", "test-app", "1").Return(false)
	transaction.On("Rollback")

	storage.On("CreateMigrationTable")
	storage.On("GetTransaction").Return(transaction)

	migration := NewMigration(storage, "test-app")

	assert.NotNil(t, migration)

	called := false

	migration.Add("1", func(tx *sqlx.Tx) error {
		called = true
		return nil
	})

	migration.Execute()

	assert.False(t, called)
	transaction.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestCallbackReturningError(t *testing.T) {
	storage := new(MockedStorage)
	transaction := new(MockedTransaction)

	transaction.On("InsertId", "test-app", "1").Return(true)
	transaction.On("Rollback")
	transaction.On("GetTx").Return(nil)

	storage.On("CreateMigrationTable")
	storage.On("GetTransaction").Return(transaction)

	migration := NewMigration(storage, "test-app")

	assert.NotNil(t, migration)

	called := false

	migration.Add("1", func(tx *sqlx.Tx) error {
		called = true
		return errors.New("Callback failed")
	})

	migration.Execute()

	assert.True(t, called)
	transaction.AssertExpectations(t)
	storage.AssertExpectations(t)
}
