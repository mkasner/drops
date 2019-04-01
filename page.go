package drops

import (
	"html/template"
	"net/http"

	"context"
)

type MenuId uint8

// Menu configuration
type MenuCfg struct {
	MenuId  MenuId
	Ordinal int
}

type App struct {
	Menu          Menu
	Pages         map[string]Page
	Widgets       map[string]Widget
	TemplatePath  string
	Subdirectory  string
	Dev           bool             // development mode, loads template from file
	TemplateFuncs template.FuncMap `json:"-"`
	WildcardMux   *WildcardMux
	NotFound      http.Handler
}

func (t *App) Init() {
	loadIds(t)
	loadHandlers(t)
	loadTemplates(t, t.TemplatePath)
	loadSubcontents(t)
}

// Plugin only instantiates handlers, not templates
// Useful when application is behaving like a plugin, and not rendering templates
func (t *App) Plugin() {
	loadIds(t)
	loadHandlers(t)
}

func (t *App) GetTemplate(id string) Template {
	tpl := t.Pages[id].Template
	if t.Dev {
		th := tpl.(*HtmlTemplate)
		tpl, s := loadTemplate(t.TemplatePath, th.File, id, t.TemplateFuncs, t.Subdirectory)
		return NewHtmlTemplate(th.Name(), th.File, s, tpl)
	}
	return tpl.Clone()
}

func (t *App) GetWidget(id string) Template {
	tpl := t.Widgets[id].Template
	if t.Dev {
		th := tpl.(*HtmlTemplate)
		tpl, s := loadTemplate(t.TemplatePath, th.File, id, t.TemplateFuncs, t.Subdirectory)
		return NewHtmlTemplate(th.Name(), th.File, s, tpl)
	}
	return tpl.Clone()
}

type Page struct {
	Id           string
	File         string
	Name         string
	Url          string
	Label        string
	Handler      http.Handler `json:"-"`
	Menu         []MenuCfg
	Data         func(context.Context) (interface{}, error) `json:"-"`
	Template     Template
	HtmlMenuItem string
	Parent       string
	Permission   Permission
	Submenu      MenuId
	Description  string
	Subcontent   []string // subcontent widgets which can be included on page
	Inactive     bool
	Attrs        map[string]string // custom configuration for page
}

func (t *Page) HasMenu(mid MenuId) bool {
	for _, m := range t.Menu {
		if mid == m.MenuId {
			return true
		}
	}
	return false
}

func (t *Page) Ordinal(mid MenuId) int {
	for _, m := range t.Menu {
		if mid == m.MenuId {
			return m.Ordinal
		}
	}
	return 0
}

type Menu struct {
	Items []MenuItem
}

type MenuItem struct {
	Label   string
	Href    string
	Html    template.HTML
	Class   string
	Data    map[string]string
	Ordinal int
	Inner   template.HTML
}

type MenuItemOrdered []MenuItem

func (a MenuItemOrdered) Len() int           { return len(a) }
func (a MenuItemOrdered) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MenuItemOrdered) Less(i, j int) bool { return a[i].Ordinal < a[j].Ordinal }

type Footer struct {
	Template
}

type Widget struct {
	Id         string
	Name       string
	File       string
	Template   Template
	Subcontent []string // subcontent widgets which can be included with this widget
}

type SpfResponse struct {
	Title string                       `json:"title"`
	Url   string                       `json:"url"`
	Head  string                       `json:"head"`
	Body  map[string]string            `json:"body"`
	Attr  map[string]map[string]string `json:"attr"`
	Foot  string                       `json:"foot"`
}

func MergeApps(app *App, merges ...App) *App {
	if app.Pages == nil {
		app.Pages = make(map[string]Page)
	}
	if app.Widgets == nil {
		app.Widgets = make(map[string]Widget)
	}
	if app.TemplateFuncs == nil {
		app.TemplateFuncs = make(template.FuncMap)
	}
	for _, merge := range merges {
		for id, page := range merge.Pages {
			app.Pages[id] = page
		}
		for id, widget := range merge.Widgets {
			app.Widgets[id] = widget
		}
		for id, fn := range merge.TemplateFuncs {
			app.TemplateFuncs[id] = fn
		}
	}
	return app
}

type Permission struct {
	Page       uint     // used to group permissions of certain part of application
	Permission []uint64 // concrete permission, checks if one of permissions is active
}

func Perm(page uint, permission ...uint64) Permission {
	return Permission{Page: page, Permission: permission}
}

func (p *Permission) IsZero() bool {
	if p.Page > 0 && len(p.Permission) > 0 {
		return false
	}
	return true
}
