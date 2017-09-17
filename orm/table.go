package orm

import (
	"reflect"
	"strings"
	"bytes"
)

const (
	TAG_COLUMN_NAME = "name"
	TAG_IGNORE      = "ignore"
	TAG_TYPE        = "type"
)

type Table struct {
	db           *DB
	name         string
	structType   reflect.Type
	column2Field map[string]string
	field2Column map[string]string
	charset      string
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
	column2Field := make(map[string]string)
	field2Column := make(map[string]string)
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
		column2Field[name] = t.Field(i).Name
		field2Column[t.Field(i).Name] = name
	}
	return &Table{
		db:           db,
		name:         name,
		structType:   t,
		column2Field: column2Field,
		field2Column: field2Column,
	}
}

func (t *Table) Charset(charset string) *Table {
	t.charset = charset
	return t
}

func (t *Table) sqlType(fieldName string) string {
	if field, ok := t.structType.FieldByName(fieldName); ok {
		tag := field.Tag
		sqlType := tag.Get(TAG_TYPE)
		if sqlType != "" {
			return sqlType
		}
		switch field.Type.Kind() {
		case reflect.String:
			return "TEXT"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			return "INT"
		case reflect.Int64, reflect.Uint64:
			return "BIGINT"
		}
	}
	return ""
}

func (t *Table) CreateIfNotExists() error {
	var buffer bytes.Buffer
	buffer.WriteString("CREATE TABLE IF NOT EXISTS ")
	buffer.WriteString(t.name)
	buffer.WriteByte('(')
	idx := 0
	for cName, fName := range t.column2Field {
		if idx != 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(cName)
		buffer.WriteByte(' ')
		buffer.WriteString(t.sqlType(fName))
	}
	buffer.WriteByte(')')
	if t.charset != "" {
		buffer.WriteString(" CHARSET SET ")
		buffer.WriteString(t.charset)
	}
	sql := buffer.String()
	_, err := t.db.db.Exec(sql)
	return err
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

func (t *Table) NewInsert(modelBean interface{}) *Insert {
	return &Insert{
		table: t,
		bean:  modelBean,
	}
}

func (t *Table) Drop() error {
	var buffer bytes.Buffer
	buffer.WriteString("DROP TABLE ")
	buffer.WriteString(t.name)
	sql := buffer.String()
	_, err := t.db.db.Exec(sql)
	return err
}
