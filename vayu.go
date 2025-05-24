// Package vayu provides a lightweight and flexible web framework for Go applications.
// It's designed to be minimal yet powerful, taking inspiration from Express.js.
package vayu

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

// App represents a vayu web application.
// It contains the router, middleware stack, and serves HTTP requests.
type App struct {
	router          *Router
	middleware      []HandlerFunc
	NotFoundHandler HandlerFunc
}

// NextFunc represents the next middleware or handler function to be called.
type NextFunc func()

// HandlerFunc represents a vayu request handler function.
// It receives a context and a next function, allowing middleware chaining.
type HandlerFunc func(*Context, NextFunc)

// New creates a new vayu application instance.
func New() *App {
	app := &App{
		router: &Router{
			routes: make(map[string][]route),
		},
	}

	// Default 404 handler
	app.NotFoundHandler = func(c *Context, _ NextFunc) {
		c.Writer.WriteHeader(StatusNotFound)
		_, err := c.Writer.Write([]byte("404 Not Found"))
		if err != nil {
			// Log error but can't do much else in a 404 handler
			log.Printf("Error writing 404 response: %v", err)
		}
	}

	return app
}

// Use adds a middleware function to the global middleware stack.
// Middleware functions are executed in the order they are added.
func (a *App) Use(mw HandlerFunc) *App {
	a.middleware = append(a.middleware, mw)
	return a
}

// addRoute registers a route with the given HTTP method, path, and handler.
func (a *App) addRoute(method, path string, handler HandlerFunc) *App {
	if a.router.routes[method] == nil {
		a.router.routes[method] = []route{}
	}
	a.router.routes[method] = append(a.router.routes[method], route{
		pattern: path,
		handler: handler,
	})
	return a
}

// GET registers a route for the GET HTTP method.
func (a *App) GET(path string, handler HandlerFunc) *App {
	return a.addRoute("GET", path, handler)
}

// POST registers a route for the POST HTTP method.
func (a *App) POST(path string, handler HandlerFunc) *App {
	return a.addRoute("POST", path, handler)
}

// PUT registers a route for the PUT HTTP method.
func (a *App) PUT(path string, handler HandlerFunc) *App {
	return a.addRoute("PUT", path, handler)
}

// DELETE registers a route for the DELETE HTTP method.
func (a *App) DELETE(path string, handler HandlerFunc) *App {
	return a.addRoute("DELETE", path, handler)
}

// PATCH registers a route for the PATCH HTTP method.
func (a *App) PATCH(path string, handler HandlerFunc) *App {
	return a.addRoute("PATCH", path, handler)
}

// OPTIONS registers a route for the OPTIONS HTTP method.
func (a *App) OPTIONS(path string, handler HandlerFunc) *App {
	return a.addRoute("OPTIONS", path, handler)
}

// HEAD registers a route for the HEAD HTTP method.
func (a *App) HEAD(path string, handler HandlerFunc) *App {
	return a.addRoute("HEAD", path, handler)
}

// ServeHTTP implements the http.Handler interface.
// This is the entry point for handling HTTP requests.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a new context with a timeout of 30 seconds
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Create the request context
	ctx := &Context{
		Writer:  NewResponseWriter(w),
		Request: r.WithContext(ctxWithTimeout),
		Params:  map[string]string{},
		Ctx:     ctxWithTimeout,
	}

	// Find matching route
	handler, params := a.router.matchRoute(r.Method, r.URL.Path)
	if handler == nil {
		// Use custom NotFoundHandler if defined
		if a.NotFoundHandler != nil {
			a.NotFoundHandler(ctx, func() {})
		} else {
			http.NotFound(w, r)
		}
		return
	}
	ctx.Params = params

	// Build middleware + handler chain
	mws := append(a.middleware, handler)

	var exec func(int)
	exec = func(i int) {
		// Check if request context is done
		select {
		case <-ctx.Ctx.Done():
			// Just write a gateway timeout if the context deadline was exceeded
			ctx.Writer.WriteHeader(StatusGatewayTimeout)
			return
		default:
			if i < len(mws) && !ctx.Stopped {
				mws[i](ctx, func() { exec(i + 1) })
			}
		}
	}

	exec(0)
}

// Static serves static files from the given directory under the specified route prefix.
// For example, Static("/assets", "./public") will serve files from ./public as /assets/filename.
func (a *App) Static(routePrefix string, dir string) *App {
	fs := http.FileServer(http.Dir(dir))
	// Strip the prefix before serving
	handler := http.StripPrefix(routePrefix, fs)
	a.Use(func(c *Context, next NextFunc) {
		if strings.HasPrefix(c.Request.URL.Path, routePrefix) {
			handler.ServeHTTP(c.Writer, c.Request)
			return
		}
		next()
	})
	return a
}

// Listen starts the HTTP server on the given address.
// For example, Listen(":8080") will start the server on port 8080.
func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a)
}

// ListenTLS starts the HTTPS server using the given certificate and key files.
func (a *App) ListenTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, a)
}

// SetNotFoundHandler sets a custom handler for 404 Not Found responses.
func (a *App) SetNotFoundHandler(handler HandlerFunc) *App {
	a.NotFoundHandler = handler
	return a
}

// splitPath splits a URL path into components for routing.
// It trims leading and trailing slashes and returns path segments.
func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}
