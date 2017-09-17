package orm

import "bytes"

type Where struct {
	condition string
	args      []interface{}
	next      *Where
}

func NewWhere(condition string, args ...interface{}) *Where {
	return &Where{
		condition:condition,
		args:args,
	}
}

func (w *Where)sql() (string, []interface{}) {
	var buffer bytes.Buffer
	args := make([]interface{}, len(w.args))
	args = append(args, w.args...)
	buffer.WriteString(" WHERE ")
	buffer.WriteString(w.condition)
	next := w.next
	for {
		if next == nil {
			break
		}
		buffer.WriteString(next.condition)
		args = append(args, next.args...)
		next = next.next
	}
	return buffer.String(), args
}

func (w *Where)appendKeyWord(key string, condition string, args ...interface{}) *Where {
	next := w
	for {
		if next.next == nil {
			break
		} else {
			next = next.next
		}
	}
	next.next = &Where{
		condition:key + condition,
		args:args,
	}
	return w
}

func (w *Where)And(condition string, args ...interface{}) *Where {
	return w.appendKeyWord(" AND ", condition, args...)
}

func (w *Where)Or(condition string, args ...interface{}) *Where {
	return w.appendKeyWord(" OR ", condition, args...)
}


