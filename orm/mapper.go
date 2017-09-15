package orm

type Mapper interface {
	ColumnName(fieldName string) string
	FieldName(columnName string) string
}

type MemoryMapper struct {
	column2Field map[string]string
	field2Column map[string]string
}

func NewMemoryMapper() Mapper {
	return &MemoryMapper{
		column2Field:make(map[string]string),
		field2Column:make(map[string]string),
	}
}

func (m *MemoryMapper)ColumnName(fieldName string) string {
	return m.field2Column[fieldName]
}

func (m *MemoryMapper)FieldName(columnName string) string {
	return m.column2Field[columnName]
}
