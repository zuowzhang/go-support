package simplehttp

import "os"

type Group struct {
	prefix  string
	filters []FilterFunc
	simple  *Simple
}

func (g *Group)Use(filter ...FilterFunc) *Group {
	g.filters = append(g.filters, filter...)
	return g
}

func (g *Group)Get(path string, h HandlerFunc, filters... FilterFunc) *Group {
	g.simple.Get(g.prefix + string(os.PathSeparator) + path, h, append(g.filters, filters...))
	return g
}

func (g *Group)Post(path string, h HandlerFunc, filters... FilterFunc) *Group {
	g.simple.Post(g.prefix + string(os.PathSeparator) + path, h, append(g.filters, filters...))
	return g
}
