package drops

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestWildcardMux(t *testing.T) {

	testData := []struct {
		pattern []string
		handler http.Handler
		paths   []string
		results []bool
	}{
		{
			pattern: []string{"/one/*/three"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/one/two/three", "/one/twotwo/three", "/one/two/four"},
			results: []bool{true, true, false},
		},
		{
			pattern: []string{"/one/*/three", "/one/*/threethree"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/one/two/threethree", "/one/twotwo/three", "/one/twotwo/threethree"},
			results: []bool{true, true, true},
		},
		{
			pattern: []string{"/car/*/audi", "/car/*/bmw"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi", "/fruit/apple", "/car/suv/bmw/i8"},
			results: []bool{true, false, false},
		},
		{
			pattern: []string{"/car/*/audi/*", "/car/*/audi/*/*"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi", "/fruit/apple", "/car/suv/audi/a5", "/car/suv/audi/a5/sedan"},
			results: []bool{false, false, true, true},
		},
		{
			pattern: []string{"/car/*/audi/*/owner"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi/a5/owner", "/car/suv/audi/a5/owner/jim"},
			results: []bool{true, false},
		},
		{
			pattern: []string{"/car/*/audi/*/*/*/owner"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi/a5/owner/owner/owner", "/car/suv/audi/a5/owner/jim/repair"},
			results: []bool{true, false},
		},
		{
			pattern: []string{"/car/*/audi/*/*/*/owner", "/car/*/audi/*/*/owner/repair"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi/a5/owner/owner/owner", "/car/suv/audi/a5/jim/owner/repair"},
			results: []bool{true, true},
		},
		{
			pattern: []string{"/car/*/audi/*/*/*/owner", "/fruit/*/slice"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi/a5/owner/owner/owner", "/fruit/apple/slice"},
			results: []bool{true, true},
		},
	}

	for _, td := range testData {
		mux := &WildcardMux{}
		for _, p := range td.pattern {
			mux.Handle(p, td.handler)
		}
		for i, p := range td.paths {
			handler := mux.handler(p)
			if handler != nil && td.results[i] == false {
				t.Fatalf("Handler shouldn't be returned for path %s\n", p)
			}
			if handler == nil && td.results[i] == true {
				t.Fatalf("Handler should be returned for path %s\n", p)
			}

		}
	}

}

func dummyHandler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "1")
}
func dummyHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "2")
}

func dummyHandler3(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "3")
}
func dummyHandler4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "4")
}
func dummyHandler5(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "5")
}

// TestWildcardMuxMatch tests if correct handler is returing
func TestWildcardMuxMatch(t *testing.T) {

	testData := []struct {
		pattern []string
		handler []http.Handler
		paths   []string
		results []string
	}{
		{
			pattern: []string{"/one/*/three"},
			handler: []http.Handler{http.HandlerFunc(dummyHandler1)},
			paths:   []string{"/one/two/three", "/one/twotwo/three", "/one/two/four"},
			results: []string{"1", "1", ""},
		},
		{
			pattern: []string{"/*/*", "/*/three"},
			handler: []http.Handler{
				http.HandlerFunc(dummyHandler1),
				http.HandlerFunc(dummyHandler2),
			},
			paths:   []string{"/one/two/four", "/one/three", "/one/two"},
			results: []string{"", "2", "1"},
		},
		{
			pattern: []string{"/*/*/*", "/*/*", "/*/*/about", "/cars/audi/*/*"},
			handler: []http.Handler{
				http.HandlerFunc(dummyHandler1),
				http.HandlerFunc(dummyHandler2),
				http.HandlerFunc(dummyHandler3),
				http.HandlerFunc(dummyHandler4),
			},
			paths:   []string{"/fruit/green/apple", "/fruit/banana", "/fruit/red/apple", "/toys/ball/about", "/cars/audi/a5/3.0"},
			results: []string{"1", "2", "1", "3", "4"},
		},
		{
			pattern: []string{"/*/*/*", "/*/*", "/*/*/about", "/cars/audi/*/*", "/*"},
			handler: []http.Handler{
				http.HandlerFunc(dummyHandler1),
				http.HandlerFunc(dummyHandler2),
				http.HandlerFunc(dummyHandler3),
				http.HandlerFunc(dummyHandler4),
				http.HandlerFunc(dummyHandler5),
			},
			paths:   []string{"/fruit/green/apple", "/fruit/banana", "/fruit/red/apple", "/toys/ball/about", "/cars/audi/a5/3.0", "/basket", "/"},
			results: []string{"1", "2", "1", "3", "4", "5", "5"},
		},
	}

	for _, td := range testData {
		mux := &WildcardMux{}
		for i, p := range td.pattern {
			mux.Handle(p, td.handler[i])
		}
		for i, p := range td.paths {
			handler := mux.handler(p)
			if handler != nil && td.results[i] == "" {
				t.Fatalf("Handler shouldn't be returned for path %s\n", p)
			}
			if handler == nil && td.results[i] != "" {
				t.Fatalf("Handler should be returned for path %s\n", p)
			}
			if td.results[i] != "" {
				bw := newBufferWriter()
				handler.ServeHTTP(bw, &http.Request{})
				if bw.buff.String() != td.results[i] {
					t.Fatalf("Wrong handler returned. Path: %s Expected: %s  Got: %s\n", p, td.results[i], bw.buff.String())
				}

			}
		}
	}

}

func TestSegmentFromWildcard(t *testing.T) {
	testData := []struct {
		pattern string
		result  int
	}{
		{pattern: "12*", result: 12},
		{pattern: "12*1", result: 0},
		{pattern: "a12*", result: 0},
		{pattern: "2*", result: 2},
		{pattern: "2**", result: 0},
	}
	for _, td := range testData {
		result := segmentFromWildcard(td.pattern)
		if result != td.result {
			t.Fatalf("Error. Pattern: %s Expected: %d   Got: %d", td.pattern, td.result, result)
		}
	}

}

func newBufferWriter() *bufferWriter {
	return &bufferWriter{buff: bytes.NewBuffer(nil), header: make(http.Header)}
}

type bufferWriter struct {
	buff   *bytes.Buffer
	header http.Header
}

func (t *bufferWriter) Header() http.Header {
	return t.header
}

func (t *bufferWriter) WriteHeader(h int) {
}

func (t *bufferWriter) Write(d []byte) (int, error) {
	return t.buff.Write(d)
}
