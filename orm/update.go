package orm

import (
	"errors"
	"bytes"
	"reflect"
)

type Update struct {
	table   *Table
	where   *Where
	kvValue map[string]interface{}
	bean    interface{}
}

func (u *Update) SetBeanValue(bean interface{}) *Update {
	u.bean = bean
	return u
}

func (u *Update) AddKV(columnName string, value interface{}) *Update {
	if u.kvValue == nil {
		u.kvValue = make(map[string]interface{})
	}
	u.kvValue[columnName] = value
	return u
}

func (u *Update) Where(where *Where) *Update {
	u.where = where
	return u
}

func (u *Update) execBean() (int64, error) {
	var args []interface{}
	var buffer, valueBuffer bytes.Buffer
	buffer.WriteString("UPDATE ")
	buffer.WriteString(u.table.name)
	buffer.WriteString(" SET(")
	valueBuffer.WriteString(" VALUES(")
	idx := 0
	v := reflect.ValueOf(u.bean)
	for cName, fName := range u.table.column2Field {
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
	if u.where != nil {
		w, wArgs := u.where.sql()
		buffer.WriteString(w)
		args = append(args, wArgs...)
	}
	sql := buffer.String()
	u.table.db.printSql(sql, args...)
	result, err := u.table.db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (u *Update) execKV() (int64, error) {
	var buffer, valueBuffer bytes.Buffer
	buffer.WriteString("UPDATE ")
	buffer.WriteString(u.table.name)
	buffer.WriteString(" SET(")
	valueBuffer.WriteString(" VALUES(")
	idx := 0
	var args []interface{}
	for k, v := range u.kvValue {
		if idx != 0 {
			buffer.WriteByte(',')
			valueBuffer.WriteByte(',')
		}
		buffer.WriteString(k)
		valueBuffer.WriteByte('?')
		args = append(args, v)
		idx++
	}
	buffer.WriteByte(')')
	valueBuffer.WriteByte(',')
	buffer.Write(valueBuffer.Bytes())
	if u.where != nil {
		w, wArgs := u.where.sql()
		buffer.WriteString(w)
		args = append(args, wArgs...)
	}
	sql := buffer.String()
	u.table.db.printSql(sql, args...)
	result, err := u.table.db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (u *Update) Exec() (int64, error) {
	if u.bean == nil && u.kvValue == nil {
		return -1, errors.New("no value for update operation")
	}
	if u.bean != nil {
		return u.execBean()
	} else {
		return u.execKV()
	}
}
