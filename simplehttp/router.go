package simplehttp

import (
	"net/http"
	"log"
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

func (r *Router)Add(method, path string, h HandlerFunc) {
	methods, ok := r.items[path]
	log.Printf("Add %s %s %t\n", method, path, ok)
	if !ok {
		methods = new(methodFunc)
	}
	switch method {
	case http.MethodGet:
		methods.get = h
	case http.MethodPost:
		methods.post = h
	}
	log.Printf("Add end")
}
