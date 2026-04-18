package framework

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type node struct {
	children      map[string]*node // static
	paramChild    *node            // :id
	wildcardChild *node            // * or *filepath

	handler  HandlerFunc
	paramKey string
}

type route struct {
	pattern string
	handler HandlerFunc
}

type staticRoute struct {
	prefix  string
	handler http.Handler
}

type Router struct {
	routes       map[string][]route
	routeTrees   map[string]*node
	staticRoutes []staticRoute
}

func NewRouter() *Router {
	return &Router{
		routes:     make(map[string][]route),
		routeTrees: make(map[string]*node),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	path = strings.Trim(path, "/")
	pathParts := strings.Split(path, "/")

	if _, ok := r.routes[method]; !ok {
		r.routes[method] = []route{}
		r.routeTrees[method] = &node{}
	}

	currentMethodNode := r.routeTrees[method]
	for i, part := range pathParts {
		// wildcard
		if strings.HasPrefix(part, "*") {
			// enforce last segment
			if i != len(pathParts)-1 {
				panic("wildcard must be the last segment")
			}

			if currentMethodNode.wildcardChild == nil {
				currentMethodNode.wildcardChild = &node{
					paramKey: strings.TrimPrefix(part, "*"),
				}
			}

			currentMethodNode = currentMethodNode.wildcardChild
			break
		}

		// handle param (:)
		if strings.HasPrefix(part, ":") {
			// check if param child exists
			if currentMethodNode.paramChild == nil {
				currentMethodNode.paramChild = &node{
					paramKey: part[1:],
				}
			}
			currentMethodNode = currentMethodNode.paramChild
			continue
		}

		// handle static
		// initialize children map if nil
		if currentMethodNode.children == nil {
			currentMethodNode.children = make(map[string]*node)
		}
		// check if child exists
		if _, ok := currentMethodNode.children[part]; !ok {
			currentMethodNode.children[part] = &node{}
		}
		currentMethodNode = currentMethodNode.children[part]
	}

	// set handler
	currentMethodNode.handler = handler

	r.routes[method] = append(r.routes[method], route{
		pattern: path,
		handler: handler,
	})
}

func (r *Router) HandleStatic(prefix string, handler http.Handler) {
	r.staticRoutes = append(r.staticRoutes, staticRoute{
		prefix:  prefix,
		handler: handler,
	})
}

func matchRoute(pattern, path string) (bool, map[string]string) {
	// check if match route path matches pattern and return a map of params
	pattern = strings.Trim(pattern, "/")
	path = strings.Trim(path, "/")

	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)

	for i := range patternParts {
		part := patternParts[i]
		value := pathParts[i]

		// param part
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = value
			continue
		}

		// static part
		if part != value {
			return false, nil
		}
	}

	return true, params
}

func matchRouteTree(methodNode *node, path string) (HandlerFunc, map[string]string, bool) {
	path = strings.Trim(path, "/")
	pathParts := strings.Split(path, "/")

	params := make(map[string]string)

	currentMethodNode := methodNode
	for i, part := range pathParts {
		// 1.check if child exists for static part first (takes priority)
		if child, ok := currentMethodNode.children[part]; ok {
			currentMethodNode = child
			continue
		}

		// 2. fallback to param child
		if currentMethodNode.paramChild != nil {
			params[currentMethodNode.paramChild.paramKey] = part
			currentMethodNode = currentMethodNode.paramChild
			continue
		}

		// 3. fallback to wildcard child
		if currentMethodNode.wildcardChild != nil {
			remainingPath := strings.Join(pathParts[i:], "/")
			paramKey := currentMethodNode.wildcardChild.paramKey
			if paramKey == "" {
				paramKey = "*"
			}
			params[paramKey] = remainingPath
			currentMethodNode = currentMethodNode.wildcardChild
			break
		}

		return nil, nil, false
	}

	return currentMethodNode.handler, params, true
}

func (r *Router) FindRoute(method, path string) (HandlerFunc, map[string]string, bool) {
	// methodRoutes, ok := r.routes[method]
	// if !ok {
	// 	return nil, nil, false
	// }

	// for _, route := range methodRoutes {
	// 	matched, params := matchRoute(route.pattern, path)
	// 	if matched {
	// 		return route.handler, params, true
	// 	}
	// }

	// return nil, nil, false

	methodNode, ok := r.routeTrees[method]
	if !ok {
		return nil, nil, false
	}

	return matchRouteTree(methodNode, path)
}

func matchStaticPrefix(path, prefix string) bool {
	path = "/" + strings.Trim(path, "/")
	prefix = "/" + strings.Trim(prefix, "/")

	// root fallback matches everything
	if prefix == "/" {
		return true
	}

	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, params, ok := r.FindRoute(req.Method, req.URL.Path); ok {
		ctx := NewContext(w, req)
		ctx.params = params
		handler(ctx)
		return
	}

	for _, static := range r.staticRoutes {
		if matchStaticPrefix(req.URL.Path, static.prefix) {
			static.handler.ServeHTTP(w, req)
			return
		}
	}

	ctx := NewContext(w, req)
	ctx.Error(http.StatusNotFound, "404 page not found!")

	// http.NotFound(w, req)
}
