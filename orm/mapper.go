package orm

type Mapper interface {
	ColumnName(fieldName string) string
	FieldName(columnName string) string
}

type MemoryMapper struct {

}
