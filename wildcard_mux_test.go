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
			pattern: []string{"/car/*/audi/*"},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.RequestURI)
			}),
			paths:   []string{"/car/suv/audi", "/fruit/apple", "/car/suv/audi/a5/sedan"},
			results: []bool{false, false, true},
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
				http.HandlerFunc(dummyHandler1),
			},
			paths:   []string{"/one/two/four", "/one/three", "/one/two/four"},
			results: []string{"1", "2", "1"},
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
