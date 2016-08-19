package drops

import (
	"html/template"
	"net/http"

	"golang.org/x/net/context"
)

type MenuId uint8

type App struct {
	Menu          Menu
	Pages         map[string]Page
	Widgets       map[string]Widget
	TemplatePath  string
	Subdirectory  string
	Dev           bool // development mode, loads template from file
	TemplateFuncs template.FuncMap
}

func (t *App) Init() {
	loadIds(t)
	loadHandlers(t)
	loadTemplates(t, t.TemplatePath)
}

func (t *App) GetTemplate(id string) Template {
	tpl := t.Pages[id].Template
	if t.Dev {
		th := tpl.(*HtmlTemplate)
		tpl, s := loadTemplate(t.TemplatePath, th.File, id, t.TemplateFuncs)
		return NewHtmlTemplate(th.Name(), th.File, s, tpl)
	}
	return tpl
}

func (t *App) GetWidget(id string) Template {
	tpl := t.Widgets[id].Template
	if t.Dev {
		th := tpl.(*HtmlTemplate)
		tpl, s := loadTemplate(t.TemplatePath, th.File, id, t.TemplateFuncs)
		return NewHtmlTemplate(th.Name(), th.File, s, tpl)
	}
	return tpl
}

type Page struct {
	Id           string
	File         string
	Name         string
	Url          string
	Label        string
	Handler      http.Handler
	Menu         MenuId
	Data         func(context.Context) (interface{}, error)
	Template     Template
	Ordinal      int
	HtmlMenuItem string
	Parent       string
	Permission   int
	Submenu      MenuId
}
type PageOrdered []Page

func (a PageOrdered) Len() int           { return len(a) }
func (a PageOrdered) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PageOrdered) Less(i, j int) bool { return a[i].Ordinal < a[j].Ordinal }

type Menu struct {
	Items []MenuItem
}

type MenuItem struct {
	Label string
	Href  string
	Html  template.HTML
	Class string
	Data  map[string]string
}

type Footer struct {
	Template
}

type Widget struct {
	Id       string
	Name     string
	File     string
	Template Template
}

type SpfResponse struct {
	Title string                       `json:"title"`
	Url   string                       `json:"url"`
	Head  string                       `json:"head"`
	Body  map[string]string            `json:"body"`
	Attr  map[string]map[string]string `json:"attr"`
	Foot  string                       `json:"foot"`
}
