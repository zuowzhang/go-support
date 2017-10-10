package simplehttp

import (
	"io"
	"html/template"
	"errors"
)

type Render interface {
	Render(io.Writer, string, interface{}) error
}

type render struct {
	template *template.Template
}

func (r *render)ParseGlob(pattern string) error {
	if template, err := template.ParseGlob(pattern); err != nil {
		return err
	} else {
		r.template = template
	}
}

func (r *render)Render(writer io.Writer, name string, data interface{}) error {
	if r.template != nil {
		return r.template.ExecuteTemplate(writer, name, data)
	}
	return errors.New("Template not initialized...")
}
