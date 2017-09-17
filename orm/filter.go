package orm

type filter interface {
	sql() (string, []interface{})
}
