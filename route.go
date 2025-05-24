package vayu

import "strings"

type route struct {
	pattern string
	handler HandlerFunc
}

type Router struct {
	routes map[string][]route
}

func (r *Router) matchRoute(method, path string) (HandlerFunc, map[string]string) {
	routes := r.routes[method]
	for _, route := range routes {
		patternParts := splitPath(route.pattern)
		pathParts := splitPath(path)

		if len(patternParts) != len(pathParts) {
			continue
		}

		params := make(map[string]string)
		match := true

		for i := range patternParts {
			if strings.HasPrefix(patternParts[i], ":") {
				params[patternParts[i][1:]] = pathParts[i]
			} else if patternParts[i] != pathParts[i] {
				match = false
				break
			}
		}

		if match {
			return route.handler, params
		}
	}
	return nil, nil
}
