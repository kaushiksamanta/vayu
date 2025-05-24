package vayu

func WithMiddleware(handler HandlerFunc, middlewares ...HandlerFunc) HandlerFunc {
	return func(c *Context, next NextFunc) {
		full := append(middlewares, handler)

		var exec func(int)
		exec = func(i int) {
			if i < len(full) {
				full[i](c, func() { exec(i + 1) })
			} else {
				next()
			}
		}
		exec(0)
	}
}
