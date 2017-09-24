package orm

import (
	"database/sql"
	"sync"
	"fmt"
)

type Log interface {
	D(string, ...interface{})
	I(string, ...interface{})
	W(string, ...interface{})
	E(string, ...interface{})
}

type DB struct {
	db      *sql.DB
	lock    sync.RWMutex
	tables  map[string]*Table
	showSql bool
	logger  Log
}

func NewDB(driverName, dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		db:     db,
		tables: make(map[string]*Table),
	}, nil
}

func (db *DB)ShowSql(show bool) *DB {
	db.showSql = show
	return db
}

func (db *DB)Logger(logger Log) *DB {
	db.logger = logger
	return db
}

func (db *DB)printSql(sql string, args ...interface{}) *DB {
	if db.showSql && db.logger != nil {
		db.logger.D(fmt.Sprintf("<%s> with args%v\n", sql, args))
	}
	return db
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
