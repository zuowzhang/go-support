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

func parserRows(rows *sql.Rows, table *Table) []interface{} {
	columnNames, err := rows.Columns()
	if err != nil {
		return nil
	}
	var beans []interface{}
	for rows.Next() {
		var fieldValues []interface{}
		bean := reflect.New(table.structType).Elem()
		for _, columnName := range columnNames {
			fieldValues = append(fieldValues,
				bean.FieldByName(table.mapper.FieldName(columnName)).Addr().Interface())
		}
		err := rows.Scan(fieldValues...)
		if err == nil {
			beans = append(beans, bean)
		}
	}
	return beans
}

func (q *Query) GetOne() interface{} {
	beans := q.Get()
	if beans != nil && len(beans) > 0 {
		return beans[0]
	}
	return nil
}

func (q *Query) Get() []interface{} {
	sql, args := q.sql()
	rows, err := q.table.db.db.Query(sql, args...)
	if err != nil {
		//show error
		return nil
	}
	return parserRows(rows, q.table)
}
