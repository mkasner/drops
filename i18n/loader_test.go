package i18n

import (
	"fmt"
	"testing"
)

func TestLoadFile(t *testing.T) {
	const (
		en LangId = iota
		de
		hr
	)
	files := []string{"sample_en.yaml", "sample.yaml"}
	err := LoadFile(files...)
	if err != nil {
		t.Fatal(err)
	}
	var translation, result string
	fmt.Println(instance)
	// test hr
	hrtr := L(hr)
	translation = hrtr("msg1")
	if translation != "poruka1" {
		t.Fatalf("Wrong translation. expected: %s got %s", "poruka1", translation)
	}
	translation = hrtr("msg2")
	if translation != "poruka2" {
		t.Fatalf("Wrong translation. expected: %s got %s", "poruka2", translation)
	}
	translation = hrtr("msg3")
	if translation != "poruka3" {
		t.Fatalf("Wrong translation. expected: %s got %s", "poruka3", translation)
	}
	// test de
	detr := L(de)
	translation = detr("msg1")
	result = "nachricht1"
	if translation != result {
		t.Fatalf("Wrong translation. expected: %s got %s", result, translation)
	}
	translation = detr("msg2")
	result = "nachricht2"
	if translation != result {
		t.Fatalf("Wrong translation. expected: %s got %s", result, translation)
	}
	translation = detr("msg3")
	if translation != "msg3" {
		t.Fatalf("Wrong translation. expected: %s got %s", "msg3", translation)
	}

}
