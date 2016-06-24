package drops

import (
	"html/template"
	"strconv"
)

type Pagination struct {
	Max     int // top number of pages
	Current int // current page
	Size    int // elements on page
}
type ButtonType int

const (
	ButtonTypePage        ButtonType = iota // absolute pages
	ButtonTypeRelative                      // prev/next
	ButtonTypePlaceholder                   // ...
)

type PageButton struct {
	Type   ButtonType
	Text   string
	Page   int
	Url    string
	Active bool
	Html   template.HTML
}

const (
	PagePrev = "Prev"
	PageNext = "Next"
)

func Pages(pagination Pagination, maxPageButtons int) []PageButton {
	var result []PageButton
	// prev
	if pagination.Current > 1 {
		result = append(result, PageButton{Type: ButtonTypeRelative, Page: pagination.Current - 1, Text: PagePrev})
	}
	if pagination.Max <= maxPageButtons {
		for i := 1; i <= pagination.Max; i++ {
			p := PageButton{Page: i, Text: strconv.Itoa(i)}
			if pagination.Current == i {
				p.Active = true
			}
			result = append(result, p)
		}
	}
	if pagination.Max > maxPageButtons {
		first := PageButton{Page: 1, Text: strconv.Itoa(1)}
		if pagination.Current == 1 {
			first.Active = true
		}
		result = append(result, first)
		// ...
		if pagination.Current > 1+2 {
			result = append(result, PageButton{Type: ButtonTypePlaceholder})
		}
		if pagination.Current == pagination.Max && pagination.Current-3 > 1 {
			result = append(result, PageButton{Page: pagination.Current - 3, Text: strconv.Itoa(pagination.Current - 3)})
		}
		if (pagination.Current+1 == pagination.Max || pagination.Current == pagination.Max) && pagination.Current-2 > 1 {
			result = append(result, PageButton{Page: pagination.Current - 2, Text: strconv.Itoa(pagination.Current - 2)})
		}
		if pagination.Current-1 > 1 {
			result = append(result, PageButton{Page: pagination.Current - 1, Text: strconv.Itoa(pagination.Current - 1)})
		}
		if pagination.Current > 1 && pagination.Current < pagination.Max {
			result = append(result, PageButton{Page: pagination.Current, Text: strconv.Itoa(pagination.Current), Active: true})
		}
		if pagination.Current+1 < pagination.Max {
			result = append(result, PageButton{Page: pagination.Current + 1, Text: strconv.Itoa(pagination.Current + 1)})
		}
		if (pagination.Current-1 == 1 || pagination.Current == 1) && pagination.Current+2 < pagination.Max {
			result = append(result, PageButton{Page: pagination.Current + 2, Text: strconv.Itoa(pagination.Current + 2)})
		}
		if pagination.Current == 1 && pagination.Current+3 < pagination.Max {
			result = append(result, PageButton{Page: pagination.Current + 3, Text: strconv.Itoa(pagination.Current + 3)})
		}
		// ...
		if pagination.Current < pagination.Max-2 {
			result = append(result, PageButton{Type: ButtonTypePlaceholder})
		}
		last := PageButton{Page: pagination.Max, Text: strconv.Itoa(pagination.Max)}
		if pagination.Current == pagination.Max {
			last.Active = true
		}
		result = append(result, last)
	}
	// next
	if pagination.Current < pagination.Max {
		result = append(result, PageButton{Type: ButtonTypeRelative, Page: pagination.Current + 1, Text: PageNext})
	}
	return result
}
