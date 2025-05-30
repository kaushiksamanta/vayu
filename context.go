// Package vayu provides a lightweight web framework for Go applications.
package vayu

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

// Context represents the context of an HTTP request.
// It encapsulates request and response objects, parameters, and control flow.
type Context struct {
	Writer  *ResponseWriter
	Request *http.Request
	Params  map[string]string
	Stopped bool
	Ctx     context.Context

	// Custom data store for request-scoped data with type information
	store map[string]any
}

// Query returns the value of the URL query parameter with the given key.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// JSON sends a JSON response with the given status code and object.
func (c *Context) JSON(code int, obj any) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer).Encode(obj)
}

// Send sends a text response with the given status code and content.
func (c *Context) Send(code int, text string) (int, error) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	return c.Writer.Write([]byte(text))
}

// Stop prevents further middleware from being invoked.
func (c *Context) Stop() {
	c.Stopped = true
}

// BindJSON binds the request body as JSON to the given struct.
func (c *Context) BindJSON(dest any) error {
	if c.Request.Body == nil {
		return io.EOF
	}
	defer c.Request.Body.Close()
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(dest)
}

// FormFile returns the uploaded file with the given form field name.
func (c *Context) FormFile(field string) (multipart.File, *multipart.FileHeader, error) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return nil, nil, err
	}
	return c.Request.FormFile(field)
}
