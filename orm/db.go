package orm

import (
	"database/sql"
	"sync"
)

type DB struct {
	db     *sql.DB
	lock   sync.RWMutex
	tables map[string]*Table
}

func NewDB(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{
		db:     db,
		tables: make(map[string]*Table),
	}, nil
}

func (db *DB) GetTable(tableName string, modelBean interface{}) *Table {
	db.lock.RLock()
	table, ok := db.tables[tableName]
	db.lock.RUnlock()
	if ok {
		return table
	}
	db.lock.Lock()
	defer db.lock.Unlock()
	table = newTable(db, tableName, modelBean)
	db.tables[tableName] = table
	return table
}
