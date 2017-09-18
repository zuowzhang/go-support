package orm

import (
	"bytes"
	"reflect"
)

type Insert struct {
	table *Table
	bean  interface{}
}

func (i *Insert) Exec() (int64, error) {
	var buffer, valueBuffer bytes.Buffer
	buffer.WriteString("INSERT INTO ")
	buffer.WriteString(i.table.name)
	buffer.WriteString(" (")
	valueBuffer.WriteString(" VALUES(")
	idx := 0
	var args []interface{}
	v := reflect.Indirect(reflect.ValueOf(i.bean))
	for cName, fName := range i.table.column2Field {
		if idx != 0 {
			buffer.WriteByte(',')
			valueBuffer.WriteByte(',')
		}
		buffer.WriteString(cName)
		valueBuffer.WriteByte('?')
		args = append(args, v.FieldByName(fName).Interface())
		idx++
	}
	buffer.WriteByte(')')
	valueBuffer.WriteByte(')')
	buffer.Write(valueBuffer.Bytes())
	sql := buffer.String()
	i.table.db.printSql(sql, args...)
	result, err := i.table.db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
