package drops

import (
	"html/template"
	"os"
	"testing"
)

func TestTemplateClone(t *testing.T) {
	tpl1 := `{{define "tpl1"}}
  {{.MainText}}
	{{block "subtpl1" .}}{{end}}
	{{end}}`
	tpl2 := `{{define "subtpl1"}}
  {{.SubText}}
	{{end}}`
	t1, err := template.New("").Parse(tpl1)
	if err != nil {
		t.Fatal(err)
	}
	t1, err = t1.Parse(tpl2)
	if err != nil {
		t.Fatal(err)
	}
	t2, err := t1.Lookup("tpl1").Clone()
	if err != nil {
		t.Fatal(err)
	}
	t1.Lookup("tpl1").Execute(os.Stdout, struct {
		MainText string
		SubText  string
	}{
		MainText: "mainText",
		SubText:  "subText",
	})
	t2.Lookup("tpl1").Execute(os.Stdout, struct {
		MainText string
		SubText  string
	}{
		MainText: "mainText2",
		SubText:  "subText2",
	})

}
