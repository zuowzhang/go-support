package simplehttp

import (
	"net/http"
	"net/url"
	"errors"
	"bytes"
	"strings"
	"log"
)

type Context interface {
	Simple() *Simple
	Request() *http.Request
	Writer() http.ResponseWriter
	Get(name string) interface{}
	Set(name string, value interface{})
	Redirect(code int, url string) error
	Render(code int, name string, data interface{}) error
	String(code int, data string) error
	Blob(code int, contentType string, data []byte) error
	FormValue(name string) string
	FormValues() (url.Values, error)
	Reset(w http.ResponseWriter, r *http.Request)
	Handler() HandlerFunc
}

type context struct {
	simple  *Simple
	request *http.Request
	writer  http.ResponseWriter
}

func (c *context)Simple() *Simple {
	return c.simple
}

func (c *context)Request() *http.Request {
	return c.request
}

func (c *context)Writer() http.ResponseWriter {
	return c.writer
}

func (c *context)Get(name string) interface{} {
	if c.simple.sessionManager != nil {
		session := c.simple.sessionManager.SessionStart(c.writer, c.request)
		if session != nil {
			return session.Get(name)
		}
	}
	return nil
}

func (c *context)Set(name string, value interface{}) {
	if c.simple.sessionManager != nil {
		session := c.simple.sessionManager.SessionStart(c.writer, c.request)
		if session != nil {
			session.Set(name, value)
		}
	}
}

func (c *context)Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return errors.New("code is invalied for redirect")
	}
	c.writer.Header().Set("Location", url)
	c.writer.WriteHeader(code)
	return nil
}

func (c *context)Render(code int, name string, data interface{}) error {
	buffer := new(bytes.Buffer)
	if c.simple.render != nil {
		return c.simple.render.Render(buffer, name, data)
	}
	return c.Blob(code, "text/html; charset=UTF-8", buffer.Bytes())
}

func (c *context)String(code int, data string) error {
	return c.Blob(code, "text/plain; charset=UTF-8", []byte(data))
}

func (c *context)Blob(code int, contentType string, data []byte) error {
	c.writer.Header().Set("Content-Type", contentType)
	c.writer.WriteHeader(code)
	_, err := c.writer.Write(data)
	return err
}

func (c *context)FormValue(name string) string {
	return c.request.Form.Get(name)
}

func (c *context)FormValues() (url.Values, error) {
	if strings.HasSuffix(c.request.Header.Get("Content-Type"), "multipart/form-data") {
		if err := c.request.ParseMultipartForm(64 * 1024); err != nil {
			return nil, err
		}
	} else {
		if err := c.request.ParseForm(); err != nil {
			return nil, err
		}
	}
	return c.request.Form, nil
}

func (c *context)Reset(w http.ResponseWriter, r *http.Request) {
	c.request = r
	c.writer = w
}

func (c *context)Handler() HandlerFunc {
	log.Printf("Handler %s\n", c.request.URL.Path)
	methods, ok := c.simple.router.items[c.request.URL.Path]
	if ok {
		switch c.request.Method {
		case http.MethodGet:
			return methods.get
		case http.MethodPost:
			return methods.post
		}
	}
	return func(c Context) error {
		return c.String(404, "NOT FOUND")
	}
}
