package orm

import (
	"bytes"
	"database/sql"
	"reflect"
)

type Query struct {
	table *Table
	where *Where
}

func (q *Query) Where(where *Where) *Query {
	q.where = where
	return q
}

func (q *Query) sql() (string, []interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString("SELECT * FROM ")
	buffer.WriteString(q.table.name)
	var args []interface{}
	if q.where != nil {
		where, tempArgs := q.where.sql()
		buffer.WriteString(where)
		args = tempArgs
	}
	return buffer.String(), args
}

func parserRows(rows *sql.Rows, table *Table) ([]interface{}, error) {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var beans []interface{}
	for rows.Next() {
		var fieldValues []interface{}
		bean := reflect.New(table.structType).Elem()
		for _, columnName := range columnNames {
			fieldValues = append(fieldValues,
				bean.FieldByName(table.column2Field[columnName]).Addr().Interface())
		}
		err := rows.Scan(fieldValues...)
		if err == nil {
			beans = append(beans, bean)
		}
	}
	return beans, err
}

func (q *Query) GetOne() (interface{}, error) {
	beans, err := q.Get()
	if beans != nil && len(beans) > 0 {
		return beans[0], err
	}
	return nil, err
}

func (q *Query) Get() ([]interface{}, error) {
	sql, args := q.sql()
	q.table.db.printSql(sql, args...)
	rows, err := q.table.db.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return parserRows(rows, q.table)
}
