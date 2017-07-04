package drops

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
)

var (
	subdirRegex = regexp.MustCompile("(href|src)=\"(.*?)\"")
)

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

func loadTemplates(app *App, path string) {
	for _, p := range app.Pages {
		if len(p.File) > 0 {
			tpl, s := loadTemplate(path, p.File, p.Id, app.TemplateFuncs, app.Subdirectory)
			page := app.Pages[p.Id]
			page.Template = NewHtmlTemplate(p.Name, p.File, s, tpl)
			app.Pages[p.Id] = page
		}
	}
	for _, p := range app.Widgets {
		if len(p.File) > 0 {
			tpl, s := loadTemplate(path, p.File, p.Id, app.TemplateFuncs, app.Subdirectory)
			widget := app.Widgets[p.Id]
			widget.Template = NewHtmlTemplate(p.Name, p.File, s, tpl)
			app.Widgets[p.Id] = widget
		}
	}
}

func loadTemplate(path, file, id string, funcs template.FuncMap, subdirectory string) (*template.Template, string) {
	b, err := ioutil.ReadFile(filepath.Join(path, file))
	if err != nil {
		panic(err)
	}
	s := string(b)
	// TODO test and finish this
	//if len(subdirectory) > 0 {
	//        s = subdirRegex.ReplaceAllString(s, fmt.Sprintf("$1=\"%s$2\"", subdirectory))
	//}
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

// calculates all subcontents for pages and widgets
// including child
// it enables hierarchical organisation of subcontent widgets
func loadSubcontents(app *App) {
	for _, p := range app.Widgets {
		p.Subcontent = loadSubcontent(app, p.Id)
		p.Subcontent = uniqueSubcontents(p.Subcontent)
		app.Widgets[p.Id] = p
	}
	for _, p := range app.Pages {
		for _, sc := range p.Subcontent {
			p.Subcontent = append(p.Subcontent, loadSubcontent(app, sc)...)
		}
		p.Subcontent = uniqueSubcontents(p.Subcontent)
		app.Pages[p.Id] = p
	}
}

func Subcontent(app *App, contentKey string) []string {
	var result []string
	result = loadSubcontent(app, contentKey)
	result = uniqueSubcontents(result)
	return result
}

func Subcontents(app *App, contentKeys []string) []string {
	if len(contentKeys) == 0 {
		return contentKeys
	}
	var result []string
	for _, contentKey := range contentKeys {
		result = append(result, loadSubcontent(app, contentKey)...)
	}
	result = uniqueSubcontents(result)
	return result
}

// loads subcontent from content and all child subcontents
// it works with recursive call to child subcontent
func loadSubcontent(app *App, contentKey string) []string {
	subcontent := app.Widgets[contentKey].Subcontent
	for _, sc := range subcontent {
		subcontent = append(subcontent, loadSubcontent(app, sc)...) // recursive call for child subcontents
	}
	return subcontent
}

func uniqueSubcontents(subcontents []string) []string {
	var result []string
	umap := make(map[string]struct{})
	for _, sc := range subcontents {
		if _, ok := umap[sc]; ok {
			continue
		}
		umap[sc] = struct{}{}
		result = append(result, sc)
	}
	return result
}
