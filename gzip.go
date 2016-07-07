package drops

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} { return gzip.NewWriter(nil) },
}

type gzipWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func (g gzipWriter) Write(data []byte) (int, error) {
	return g.gw.Write(data)
}

func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if !acceptsGzip(r) {
			h.ServeHTTP(w, r)
			return
		}
		gzw := gzipWriterPool.Get().(*gzip.Writer)

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{gzw, w}, r)
		gzw.Reset(w)
		gzipWriterPool.Put(gzw)
		//gzw.Close()
	})
}

func acceptsGzip(r *http.Request) bool {
	encoding := r.Header.Get("Accept-Encoding")
	if len(encoding) == 0 {
		return false
	}

	if strings.Contains(encoding, "gzip") {
		return true
	}
	return false
}
