package orm

import "reflect"

const TAG_ORM string = "orm"

type Table struct {
	name  string
	structType reflect.Type
	mapper     Mapper
}

func NewTable(name string, modelBean interface{}) *Table {
	t := reflect.TypeOf(modelBean)
	mapper := NewMemoryMapper()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get(TAG_ORM)
		if tag != "" {

		}
	}
	return &Table{
		name:name,
		structType:t,
		mapper:mapper,
	}
}
