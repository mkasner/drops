package drops

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"
)

func GetData(ctx context.Context, app *App, contentKey string) (interface{}, error) {
	ctx = AppContext(ctx, app)
	if page, ok := app.Pages[contentKey]; ok {
		return page.Data(ctx)
	}
	return nil, nil
}

func loadIds(app *App) {
	for id, p := range app.Pages {
		p.Id = id
		app.Pages[id] = p
	}
	for id, p := range app.Widgets {
		p.Id = id
		app.Widgets[id] = p
	}
}

func loadHandlers(app *App) {
	for _, p := range app.Pages {
		if len(p.Url) > 0 && p.Handler != nil {
			http.Handle(fmt.Sprintf("%s%s", app.Subdirectory, p.Url), p.Handler)
		}
	}

}

// deprecated
//func GetTemplate(app *App, id string, tpl Template) Template {
//        if app.Dev {
//                th := tpl.(*HtmlTemplate)
//                tpl, s := loadTemplate(app.TemplatePath, th.File, id, app.TemplateFuncs)
//                return NewHtmlTemplate(th.Name(), th.File, s, tpl)
//        }
//        return tpl
//}

func loadTemplates(app *App, path string) {
	for _, p := range app.Pages {
		if len(p.File) > 0 {
			tpl, s := loadTemplate(path, p.File, p.Id, app.TemplateFuncs)
			page := app.Pages[p.Id]
			page.Template = NewHtmlTemplate(p.Name, p.File, s, tpl)
			app.Pages[p.Id] = page
		}
	}
	for _, p := range app.Widgets {
		if len(p.File) > 0 {
			tpl, s := loadTemplate(path, p.File, p.Id, app.TemplateFuncs)
			widget := app.Widgets[p.Id]
			widget.Template = NewHtmlTemplate(p.Name, p.File, s, tpl)
			app.Widgets[p.Id] = widget
		}
	}
	return
}

func loadTemplate(path, file, id string, funcs template.FuncMap) (*template.Template, string) {
	b, err := ioutil.ReadFile(filepath.Join(path, file))
	if err != nil {
		panic(err)
	}
	s := string(b)
	tpl := template.Must(template.New(id).Funcs(funcs).Parse(s))
	return tpl, s
}

func WrapContent(buff *bytes.Buffer, key string) *bytes.Buffer {
	str := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", key, buff.String())
	buff.Reset()
	buff.WriteString(str)
	return buff
}

// adds to existing template without affecting it
// uses clone for that
func AddToTemplate(buff *bytes.Buffer, tpl *template.Template) (*template.Template, error) {
	tpl, err := tpl.Clone()
	if err != nil {
		return tpl, err
	}
	return tpl.Parse(buff.String())
}
