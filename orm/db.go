package orm

import "database/sql"

type DB struct {
	db     *sql.DB
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

func (db *DB) getTable(tableName string, modelBean interface{}) *Table {
	return newTable(db, tableName, modelBean)
}
