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

func (g *Group)Get(path string, h HandlerFunc) *Group {
	g.simple.Get(g.prefix + os.PathSeparator + path, h)
	return g
}

func (g *Group)Post(path string, h HandlerFunc) *Group {
	g.simple.Post(g.prefix + os.PathSeparator + path, h)
	return g
}
