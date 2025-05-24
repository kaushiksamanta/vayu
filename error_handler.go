package vayu

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

// ErrorHandler represents a function that handles errors in middleware or handlers
type ErrorHandler func(c *Context, err error)

// SilentMode controls whether panic logs are output
// Set to true to suppress panic logs (useful in tests)
// This can be controlled via linker flags: -ldflags="-X 'github.com/kaushiksamanta/vayu.SilentMode=true'"
var SilentMode string

// isSilentMode returns true if SilentMode is set to "true" or if we're running tests
func isSilentMode() bool {
	return SilentMode == "true" || (SilentMode == "" && strings.Contains(os.Args[0], "test"))
}

// LogPanic logs panic information if not in silent mode
func LogPanic(err error, stackTrace []byte) {
	if !isSilentMode() {
		log.Printf("Panic recovered: %v\n%s", err, stackTrace)
	}
}

// DefaultErrorHandler is the default error handler
var DefaultErrorHandler = func(c *Context, err error) {
	log.Printf("Error: %v", err)
	if !c.Writer.Written() {
		c.InternalServerError("An unexpected error occurred")
	}
}

// WithErrorHandling wraps a handler with error handling
func WithErrorHandling(handler HandlerFunc, errorHandler ErrorHandler) HandlerFunc {
	if errorHandler == nil {
		errorHandler = DefaultErrorHandler
	}

	return func(c *Context, next NextFunc) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				LogPanic(err, stackTrace)
				errorHandler(c, err)
			}
		}()

		handler(c, next)
	}
}

// ErrorHandlerMiddleware returns middleware that handles errors
func ErrorHandlerMiddleware(errorHandler ErrorHandler) HandlerFunc {
	if errorHandler == nil {
		errorHandler = DefaultErrorHandler
	}

	return func(c *Context, next NextFunc) {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := debug.Stack()
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				LogPanic(err, stackTrace)
				errorHandler(c, err)
			}
		}()

		next()
	}
}
