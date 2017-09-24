package orm

import "bytes"

type Delete struct {
	table *Table
	where *Where
}

func (d *Delete) Where(where *Where) *Delete {
	d.where = where
	return d
}

func (d *Delete) Exec() (int64, error) {
	var buffer bytes.Buffer
	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(d.table.name)
	var args []interface{}
	if d.where != nil {
		sql, wArgs := d.where.sql()
		buffer.WriteString(sql)
		args = wArgs
	}
	sql := buffer.String()
	d.table.db.printSql(sql, args...)
	result, err := d.table.db.db.Exec(sql, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
