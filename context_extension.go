package vayu

import (
	"context"
	"time"
)

// Set stores a value in the request context with the given key.
func (c *Context) Set(key string, value interface{}) {
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = value
}

// Get retrieves a value from the request context with the given key.
func (c *Context) Get(key string) (interface{}, bool) {
	if c.store == nil {
		return nil, false
	}
	val, ok := c.store[key]
	return val, ok
}

// HTML sends an HTML response with the given status code and content.
func (c *Context) HTML(code int, html string) (int, error) {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(code)
	return c.Writer.Write([]byte(html))
}

// Status sets the HTTP response status code.
func (c *Context) Status(code int) *Context {
	c.Writer.WriteHeader(code)
	return c
}

// WithTimeout creates a new context with the given timeout.
func (c *Context) WithTimeout(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(c.Ctx, timeout)
	// Store the cancel function
	go func() {
		<-ctx.Done()
		cancel() // Call cancel when the context is done
	}()
	c.Ctx = ctx
}
