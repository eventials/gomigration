package gomigration

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
)

type MigrationCallBack func(tx *sqlx.Tx) error

type Migration interface {
	Add(id string, callback MigrationCallBack) *Node
	Execute() error
	Clear()
}

type MigrationImpl struct {
	trunk     *Node
	storage   Storage
	wg        *sync.WaitGroup
	appName   string
	errors    map[string]error
	processed []string
}

func (m *MigrationImpl) Add(id string, callback MigrationCallBack) *Node {
	return m.trunk.Add(id, callback)
}

func (m *MigrationImpl) setup() {
	m.storage.CreateMigrationTable()
}

func (m *MigrationImpl) callMethods(n *Node) {
	if n == nil {
		return
	}
	for _, method := range n.Methods {
		transaction := m.storage.GetTransaction()

		if transaction.InsertId(m.appName, method.Id) {
			err := method.Callback(transaction.GetTx())

			if err == nil {
				transaction.Commit()

				m.processed = append(m.processed, method.Id)
			} else {
				m.errors[method.Id] = err

				transaction.Rollback()
			}
		} else {
			transaction.Rollback()
		}

		m.wg.Add(1)
		go func(method Method) {
			m.callMethods(method.Next)

			m.wg.Done()
		}(method)
	}
}

func getMapKeys(m map[string]bool) []string {
	r := []string{}

	for k, v := range m {
		if v {
			r = append(r, k)
		}
	}

	return r
}

func duplicatedIds(n *Node, duplicated map[string]bool) {
	if n == nil {
		return
	}

	for _, method := range n.Methods {
		_, exists := duplicated[method.Id]
		if !exists {
			duplicated[method.Id] = false
		} else {
			duplicated[method.Id] = true
		}

		duplicatedIds(method.Next, duplicated)
	}
}

func (m *MigrationImpl) validate() error {
	duplicated := make(map[string]bool)
	duplicatedIds(m.trunk, duplicated)
	keys := getMapKeys(duplicated)

	if len(keys) > 0 {
		return fmt.Errorf("Migration has duplicated ids: %s", keys)
	}

	return nil
}

func (m *MigrationImpl) Clear() {
	m.storage.DeleteMigrations(m.appName)
}

func (m *MigrationImpl) Execute() error {
	err := m.validate()

	if err != nil {
		return err
	}

	m.errors = make(map[string]error)

	m.callMethods(m.trunk)

	m.wg.Wait()

	if len(m.errors) > 0 {
		return fmt.Errorf("Some migrations failed: %s", m.errors)
	}

	return nil
}

func NewMigration(storage Storage, appName string) Migration {
	m := &MigrationImpl{
		storage: storage,
		trunk:   &Node{nil, []Method{}},
		wg:      &sync.WaitGroup{},
		appName: appName,
		errors:  nil,
	}
	m.setup()
	return m
}
