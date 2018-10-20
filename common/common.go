package common

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
	"time"
)

//DontCache sets an HTTP response header to prevent rogue caches
// from caching any responses we haven't explicitly marked as cacheable.
// Subsequent calls to w.Header().Set("Cache-Control" ...) will override this value
func DontCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

//Timer logs the response time for every request
func Timer(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%d %s", time.Since(start).Nanoseconds()/1e3, r.URL.Path)
	}
}

//Logger logs the request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLogResponseWriter(w)
		defer func() {
			addr := r.RemoteAddr
			if i := strings.LastIndex(addr, ":"); i != -1 {
				addr = addr[:i]
			}
			url := strings.TrimPrefix(r.URL.String(), "/")
			log.Printf("[%s] %s /%s -%d", addr, r.Method, url, lrw.code)
		}()
		next.ServeHTTP(lrw, r)
	})
}

//Gzip applies gzip compression to responses whose corresponding
//request's headers is gzip enanless
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(GzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
