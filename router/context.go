package router

import (
	"net/http"
	"github.com/zuowzhang/go-support/orm"
	"github.com/zuowzhang/go-support/session"
)

type Context interface {
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Table(name string) *orm.Table
	SessionManager() *session.SessionManager
}

type context struct {

}

func (c *context)Request() *http.Request {
	return c.Request()
}

func (c *context)ResponseWriter() http.ResponseWriter {
	return c.ResponseWriter()
}


