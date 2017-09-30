package simplehttp

import (
	"net/http"
	"net/url"
)

type Context interface {
	Simple() *Simple
	Request() *http.Request
	Writer() http.ResponseWriter
	Get(name string) interface{}
	Set(name string, value interface{})
	Redirect(code int, url string)
	Render(code int, name string, data interface{})
	FormValue(name string) string
	FormValues() (url.Values, error)
	Reset(w http.ResponseWriter, r *http.Request)
}

type context struct {
	simple  *Simple
	request *http.Request
	writer  http.ResponseWriter
}
