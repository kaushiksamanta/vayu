package vayu

import (
	"fmt"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context, next NextFunc) {
		start := time.Now()
		next()
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s (%s)\n", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, duration)
	}
}
