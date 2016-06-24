package i18n

import "testing"

const (
	en LangId = iota
	hr
)

var (
	bundle Bundle = map[string][]string{
		"msg1": []string{"msg1", "poruka1"},
		"msg3": []string{"", "poruka3"},
		"msg4": []string{"msg4 %s", "poruka4 %s"},
	}
)

func TestLoad(t *testing.T) {
	Load(bundle)
	translation := T("msg1")
	if translation != "msg1" {
		t.Fatalf("Wrong translation", "msg1", translation)
	}
	translation = L(hr)("msg1")
	if translation != "poruka1" {
		t.Fatalf("Wrong translation. expected: %s got %s", "poruka1", translation)
	}
	translation = T("msg2")
	if translation != "msg2" {
		t.Fatalf("Wrong translation. expected: %s got %s", "msg2", translation)
	}
	translation = L(hr)("msg2")
	if translation != "msg2" {
		t.Fatalf("Wrong translation. expected: %s got %s", "msg2", translation)
	}
	translation = T("msg3")
	if translation != "msg3" {
		t.Fatalf("Wrong translation expected: %s got %s", "msg3", translation)
	}
	translation = T("msg4", "f")
	if translation != "msg4 f" {
		t.Fatalf("Wrong translation. expected: %s got %s", "msg4 f", translation)
	}
	translation = L(hr)("msg4", "f")
	if translation != "poruka4 f" {
		t.Fatalf("Wrong translation. expected: %s got %s", "poruka4 f", translation)
	}

}
