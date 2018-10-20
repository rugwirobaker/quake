package common

import (
	"io"
	"net/http"
)

//LogResponseWriter wraps the builtin net/http ResponseWriter
// for response logging purposes.
type LogResponseWriter struct {
	http.ResponseWriter
	code int
}

//GzipResponseWriter wraps the builtin net/http ResponseWriter
//to add Gzip compression functionality
type GzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

//NewLogResponseWriter creates a new loggingResponse type.
func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{w, http.StatusOK}
}

//WriteHeader initializes loggingResponseWriter values.
func (lrw *LogResponseWriter) WriteHeader(code int) {
	lrw.code = code
	lrw.ResponseWriter.WriteHeader(code)
}
