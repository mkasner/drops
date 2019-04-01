package drops

import (
	"bytes"
	"html/template"
)

type TemplateRepository struct {
	Templates map[string]Template
}

type Template interface {
	Render(*bytes.Buffer, interface{}) (*bytes.Buffer, error)
	Clone() Template
}

type HtmlTemplate struct {
	name    string
	File    string
	Content string
	Tpl     *template.Template
}

func (t *HtmlTemplate) Render(rcvr *bytes.Buffer, data interface{}) (*bytes.Buffer, error) {
	if len(t.name) == 0 {
		err := t.Tpl.Execute(rcvr, data)
		return rcvr, err
	}
	tpl, _ := t.Tpl.Clone()
	err := tpl.ExecuteTemplate(rcvr, t.name, data)
	return rcvr, err
}
func (t *HtmlTemplate) Name() string {
	return t.name
}

func (t *HtmlTemplate) Clone() Template {
	cln, _ := t.Tpl.Clone()
	return &HtmlTemplate{name: t.name, File: t.File, Content: t.Content, Tpl: cln}
}

func NewHtmlTemplate(name, file string, content string, t *template.Template) *HtmlTemplate {
	return &HtmlTemplate{
		name:    name,
		File:    file,
		Content: content,
		Tpl:     t,
	}
}

type SimpleTemplate struct {
	name       string
	RenderFunc func(*bytes.Buffer, interface{}) (*bytes.Buffer, error)
}

func (t *SimpleTemplate) Render(rcvr *bytes.Buffer, data interface{}) (*bytes.Buffer, error) {
	return t.RenderFunc(rcvr, data)
}

func (t *SimpleTemplate) Clone() Template {
	return t
}

func NewSimpleTemplate(name string, rf func(*bytes.Buffer, interface{}) (*bytes.Buffer, error)) *SimpleTemplate {
	return &SimpleTemplate{
		name:       name,
		RenderFunc: rf,
	}
}
