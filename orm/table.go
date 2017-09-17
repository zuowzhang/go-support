package orm

import (
	"reflect"
	"strings"
	"bytes"
)

const (
	TAG_COLUMN_NAME = "name"
	TAG_IGNORE      = "ignore"
)

type Table struct {
	DB         *DB
	name       string
	structType reflect.Type
	mapper     Mapper
}

func getDefaultColumnName(fieldName string) string {
	buffer := bytes.NewBuffer([]byte{})
	for idx, r := range fieldName {
		if r >= 'A' && r <= 'Z' {
			if idx != 0 {
				buffer.WriteByte('_')
			}
			r += 'a' - 'A'
		}
		buffer.WriteRune(r)
	}
	return string(buffer.Bytes())
}

func newTable(db *DB, name string, modelBean interface{}) *Table {
	t := reflect.TypeOf(modelBean)
	mapper := NewMemoryMapper()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		ignore := tag.Get(TAG_IGNORE)
		if strings.ToLower(ignore) == "true" {
			continue
		}
		name := tag.Get(TAG_COLUMN_NAME)
		if name == "" {
			name = getDefaultColumnName(t.Field(i).Name)
		}
		mapper.(*MemoryMapper).column2Field[name] = t.Field(i).Name
		mapper.(*MemoryMapper).field2Column[t.Field(i).Name] = name
	}
	return &Table{
		name:       name,
		structType: t,
		mapper:     mapper,
	}
}

func (t *Table) Sync() {

}

func (t *Table) NewQuery() *Query {
	return &Query{
		table: t,
	}
}

func (t *Table) NewUpdate() *Update {
	return &Update{
		table: t,
	}
}

func (t *Table) NewDelete() *Delete {
	return &Delete{
		table: t,
	}
}

func (t *Table) NewInsert() *Insert {
	return &Insert{
		table: t,
	}
}

func (t *Table) Drop() error {

	return nil
}
