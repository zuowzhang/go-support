package orm

import "testing"

func TestNewWhere(t *testing.T) {
	where := NewWhere("id = ?", 1).
		And("name = ?", "zhangsan").
		Or("age > ?", 18)
	t.Log(where.sql())
}
