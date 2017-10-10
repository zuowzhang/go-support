package simplehttp

import "net/http"

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

func (r *Router)Add(method, path string, h HandlerFunc) {
	methods, ok := r.items[path]
	if !ok {
		methods = new(methodFunc)
	}
	switch method {
	case http.MethodGet:
		methods.get = h
	case http.MethodPost:
		methods.post = h
	}
}
