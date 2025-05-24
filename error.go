package vayu

import "fmt"

func Recovery() HandlerFunc {
	return func(c *Context, next NextFunc) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("❌ Panic recovered: %v\n", r)
				err := c.JSON(500, map[string]string{"error": "Internal Server Error"})
				if err != nil {
					fmt.Printf("❌ Error sending JSON response: %v\n", err)
				}
			}
		}()
		next()
	}
}
