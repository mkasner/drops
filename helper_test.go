package drops

import (
	"testing"
)

func TestLoadSubcontents(t *testing.T) {
	app := &App{
		Pages: map[string]Page{
			"page1": {Subcontent: []string{"widget1"}},
			"page2": {Subcontent: []string{"widget2"}},
			"page3": {Subcontent: []string{"widget2", "widget1"}},
		},
		Widgets: map[string]Widget{
			"widget1": {},
			"widget2": {Subcontent: []string{"widget1"}},
		},
	}

	testData := []struct {
		page            string
		subcontentCount int
	}{
		{page: "page1", subcontentCount: 1},
		{page: "page2", subcontentCount: 2},
		{page: "page3", subcontentCount: 2},
	}

	loadIds(app)
	loadSubcontents(app)
	for _, td := range testData {
		if len(app.Pages[td.page].Subcontent) != td.subcontentCount {
			t.Fatal("Page subcontent not loaded", td.page, "expected", td.subcontentCount, "got", len(app.Pages[td.page].Subcontent))
		}
	}
}

func TestSubcontents(t *testing.T) {
	app := &App{
		Pages: map[string]Page{
			"page1": {Subcontent: []string{"widget1"}},
			"page2": {Subcontent: []string{"widget2"}},
			"page3": {Subcontent: []string{"widget2", "widget1"}},
		},
		Widgets: map[string]Widget{
			"widget1": {},
			"widget2": {Subcontent: []string{"widget1"}},
		},
	}

	testData := []struct {
		widgets         []string
		subcontentCount int
	}{
		{widgets: []string{"widget1"}, subcontentCount: 1},
		{widgets: []string{"widget1", "widget2"}, subcontentCount: 2},
		{widgets: []string{"widget2"}, subcontentCount: 2},
	}

	loadIds(app)
	for _, td := range testData {
		subcontents := Subcontents(app, td.widgets)
		if len(subcontents) != td.subcontentCount {
			t.Fatal("Subcontents not loaded", td.widgets, "expected", td.subcontentCount, "got", len(subcontents))
		}
	}
}
