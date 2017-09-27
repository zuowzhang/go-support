package router

import (
	"sync"
	"net/http"
	"net"
)

type Router struct {
	pool       sync.Pool
	server     *http.Server
	listener   net.Listener
	preFilters []Filter
	filters    []Filter
}

func New() *Router {
	r := &Router{
		server:new(http.Server),
	}
	r.pool.New = r.newContext()
	return r
}

func (r *Router)Pre(filters ...Filter) *Router {
	r.preFilters = append(r.filters, filters...)
	return r
}

func (r *Router)User(filters ...Filter) *Router {
	r.filters = append(r.filters, filters...)
}

func (r *Router)run(address string) error {
	r.server.Addr = address
	if r.listener == nil {

	}
	return nil
}

func (r *Router)newContext() Context {
	return &context{

	}
}
