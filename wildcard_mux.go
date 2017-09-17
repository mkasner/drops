package drops

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	wildardRegex = regexp.MustCompile(`^(\d+)\*$`)
)

// WildcardMux is meant for routing paths to correct handlers
// but with one crucial feature: It can contain wildcard (*) in the route.
// With that you can have dynamic parts of url's
type WildcardMux struct {
	m        []muxEntry
	NotFound http.Handler
}

type muxEntry struct {
	pattern  []string
	score    int // how many non empty patterns are there
	segments int // how many segments does this pattern cover
	h        http.Handler
}

func NewWildcardMux() *WildcardMux {
	return &WildcardMux{}
}

func (t *WildcardMux) Handle(pattern string, handler http.Handler) {

	// initialize if not initialized
	if t.m == nil {
		t.m = make([]muxEntry, 0)
	}
	split := strings.Split(pattern, "/")
	entry := muxEntry{
		pattern: make([]string, len(split)),
	}

	for i, segment := range split {
		if i == 0 && len(segment) == 0 {
			// ignore first empty segment
			continue
		}
		if segment == "*" {
			// leave path segment empty
			continue
		}
		// add path segment to position
		entry.pattern[i] = segment
		entry.score++
	}
	if len(entry.pattern) > 0 {
		entry.h = handler
		t.m = append(t.m, entry)
	}
}

func (t *WildcardMux) handler(path string) http.Handler {
	split := strings.Split(path, "/")
	var candidate muxEntry
	for _, entry := range t.m {
		if len(split) < len(entry.pattern) {
			// discard paths which are smaller than patterns
			continue
		}
		// if path is not same size as pattern, discard it
		if len(split) != len(entry.pattern) {
			continue
		}

		var matchScore int
		for i, segment := range entry.pattern {
			if len(segment) == 0 {
				// ignore empty pattern
				continue
			}
			if segment == split[i] {
				matchScore++
			}
		}
		if matchScore == entry.score {
			// if last element is wildcard, try to find exact match
			//if entry.pattern[len(entry.pattern)-1] == "" && candidate.h == nil {
			if len(entry.pattern) == len(split) && candidate.h == nil {
				candidate = entry
				continue
			}
			// if another wildcard route is found, skip it, use the first one
			if entry.pattern[len(entry.pattern)-1] == "" && candidate.h == nil {
				continue
			}
			return entry.h
		}
	}
	if candidate.h != nil {
		return candidate.h
	}
	return nil

}
func (t *WildcardMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := t.handler(r.URL.Path)
	if h == nil {
		if t.NotFound != nil {
			// custom notFound handler
			t.NotFound.ServeHTTP(w, r)
		} else {
			http.NotFoundHandler().ServeHTTP(w, r)
		}
		return
	}
	h.ServeHTTP(w, r)
}

// used for repeating for segments
func segmentFromWildcard(pattern string) int {
	submatch := wildardRegex.FindStringSubmatch(pattern)
	if len(submatch) == 2 {
		segmentCount, _ := strconv.Atoi(submatch[1])
		return segmentCount
	}
	return 0
}
