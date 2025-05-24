package vayu

import "strings"

type Group struct {
	prefix string
	app    *App
}

func (a *App) Group(prefix string) *Group {
	return &Group{prefix: prefix, app: a}
}

func (g *Group) Use(mw HandlerFunc) {
	g.app.Use(func(c *Context, next NextFunc) {
		if strings.HasPrefix(c.Request.URL.Path, g.prefix) {
			mw(c, next)
		} else {
			next()
		}
	})
}

func (g *Group) GET(path string, handler HandlerFunc) {
	g.app.GET(g.prefix+path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc) {
	g.app.POST(g.prefix+path, handler)
}
