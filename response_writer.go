package vayu

import (
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that tracks
// whether a response has been written yet.
type ResponseWriter struct {
	http.ResponseWriter
	written bool
	status  int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		written:        false,
		status:         0,
	}
}

// WriteHeader sets the status code for the response and marks
// the response as written.
func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Write writes the data to the connection and marks the response
// as written.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

// Written returns whether the response has been written yet.
func (w *ResponseWriter) Written() bool {
	return w.written
}

// Status returns the status code of the response.
func (w *ResponseWriter) Status() int {
	return w.status
}
