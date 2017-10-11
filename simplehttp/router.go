package simplehttp

import (
	"net/http"
)

type methodFunc struct {
	get  HandlerFunc
	post HandlerFunc
}

type Router struct {
	items map[string]*methodFunc
}

func NewRouter() *Router {
	return &Router{
		items:make(map[string]*methodFunc),
	}
}

func (r *Router)Add(method, path string, h HandlerFunc, filters... FilterFunc) {
	methods, ok := r.items[path]
	if !ok {
		methods = new(methodFunc)
		r.items[path] = methods
	}
	handler := func(c Context) error {
		for i := len(filters) - 1; i >= 0; i-- {
			h = filters[i](h)
		}
		return h(c)
	}
	switch method {
	case http.MethodGet:
		methods.get = handler
	case http.MethodPost:
		methods.post = handler
	}
}
