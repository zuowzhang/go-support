package router

import (
	"net/http"
	"strings"
)

type methodHandler struct {
	get  http.HandlerFunc
	post http.HandlerFunc
}

type Location struct {
	app      *App
	handlers map[string]*methodHandler
}

func NewLocation(app *App) *Location {
	return &Location{
		app:      app,
		handlers: make(map[string]*methodHandler),
	}
}

func (l *Location) handlerByPath(path string) *methodHandler {
	handler, ok := l.handlers[path]
	if !ok {
		handler = new(methodHandler)
	}
	return handler
}

func (l *Location) Register(method, path string, handlerFunc http.HandlerFunc) *Location {
	handler := l.handlerByPath(path)
	switch strings.ToUpper(method) {
	case METHOD_GET:
		handler.get = handlerFunc
	case METHOD_POST:
		handler.post = handlerFunc
	}
	return l
}

func (l *Location) Get(path string, handlerFunc http.HandlerFunc) *Location {
	handler := l.handlerByPath(path)
	handler.get = handlerFunc
	return l
}

func (l *Location) Post(path string, handlerFunc http.HandlerFunc) *Location {
	handler := l.handlerByPath(path)
	handler.post = handlerFunc
	return l
}
