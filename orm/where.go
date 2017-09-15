package orm

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


