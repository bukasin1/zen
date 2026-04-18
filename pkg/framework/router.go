package framework

import (
	"net/http"
	"regexp"
	"strings"
)

var multiSlashRegex = regexp.MustCompile(`/+`)

func normalizeRoutePath(path string) string {
	if path == "" {
		return "/"
	}

	// Ensure leading slash
	path = "/" + strings.Trim(path, "/")

	// Collapse multiple slashes: "/users//1" -> "/users/1"
	path = multiSlashRegex.ReplaceAllString(path, "/")

	// Remove trailing slash (except root)
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}

// validateRoute checks if the route path is valid.
// Most importantly, it checks for valid wildcard patterns
func validateRoutePath(path string) {
	if strings.Contains(path, "**") {
		panic("invalid route: double wildcard not allowed")
	}
	if strings.Contains(path, "*/") && !strings.HasSuffix(path, "*") {
		panic("invalid route: wildcard must be last segment")
	}
}

func cloneParams(p map[string]string) map[string]string {
	copy := make(map[string]string)
	for k, v := range p {
		copy[k] = v
	}
	return copy
}

type HandlerFunc func(*Context)

type node struct {
	children      map[string]*node // static
	paramChild    *node            // :id
	wildcardChild *node            // * or *filepath

	handler  HandlerFunc
	paramKey string
}

type cacheKey struct {
	method string
	path   string
}

type cachedRoute struct {
	handler HandlerFunc
	params  map[string]string
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
	routeTrees map[string]*node
	cache      map[cacheKey]cachedRoute

	// TODO: cleanup(remove old routes registering)
	routes       map[string][]route
	staticRoutes []staticRoute
}

func NewRouter() *Router {
	return &Router{
		routeTrees: make(map[string]*node),
		cache:      make(map[cacheKey]cachedRoute),

		// TODO: cleanup(remove old routes registering)
		routes: make(map[string][]route),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	// normalize path to handle multiple slashes and trailing slashes
	path = normalizeRoutePath(path)
	validateRoutePath(path)

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

			if currentMethodNode.wildcardChild != nil {
				panic("conflicting wildcard route")
			}

			currentMethodNode.wildcardChild = &node{
				paramKey: strings.TrimPrefix(part, "*"),
			}

			// set handler to base wildcard parent node
			currentMethodNode.handler = handler

			currentMethodNode = currentMethodNode.wildcardChild
			break
		}

		// handle param (:)
		if strings.HasPrefix(part, ":") {
			// check if param child exists
			if currentMethodNode.paramChild != nil {
				panic("conflicting param route at same segment")
			}

			currentMethodNode.paramChild = &node{
				paramKey: part[1:],
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
	if currentMethodNode.handler != nil {
		panic("duplicate route registration")
	}
	currentMethodNode.handler = handler

	// TODO: cleanup(remove old routes registering)
	r.routes[method] = append(r.routes[method], route{
		pattern: path,
		handler: handler,
	})
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

	if currentMethodNode.handler == nil {
		return nil, nil, false
	}

	return currentMethodNode.handler, params, true
}

func (r *Router) FindRoute(method, path string) (HandlerFunc, map[string]string, bool) {
	// normalize path to handle multiple slashes and trailing slashes
	path = normalizeRoutePath(path)

	pathCacheKey := cacheKey{
		method: method,
		path:   path,
	}

	// check cache first
	if cachedRoute, ok := r.cache[pathCacheKey]; ok {
		return cachedRoute.handler, cachedRoute.params, true
	}

	methodNode, ok := r.routeTrees[method]
	if !ok {
		return nil, nil, false
	}

	handler, params, ok := matchRouteTree(methodNode, path)
	if !ok || handler == nil {
		return nil, nil, false
	}

	// store in cache
	r.cache[pathCacheKey] = cachedRoute{
		handler: handler,
		params:  cloneParams(params),
	}

	return handler, params, true
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, params, ok := r.FindRoute(req.Method, req.URL.Path); ok {
		ctx := NewContext(w, req)
		ctx.params = params
		handler(ctx)
		// log.Println(req.Method, req.URL.Path, req.Response, "in router")
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

// Deprecated: use Static instead
// TODO: remove this function
func matchStaticPrefix(path, prefix string) bool {
	path = "/" + strings.Trim(path, "/")
	prefix = "/" + strings.Trim(prefix, "/")

	// root fallback matches everything
	if prefix == "/" {
		return true
	}

	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

// Deprecated: use Static instead
// TODO: remove this function
func (r *Router) HandleStatic(prefix string, handler http.Handler) {
	r.staticRoutes = append(r.staticRoutes, staticRoute{
		prefix:  prefix,
		handler: handler,
	})
}

// Deprecated: use matchRouteTree instead
// TODO: remove this function
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

// Deprecated: use FindRoute instead
// TODO: remove this function
func (r *Router) FindRouteOld(method, path string) (HandlerFunc, map[string]string, bool) {
	// normalize path to handle multiple slashes and trailing slashes
	path = normalizeRoutePath(path)

	methodRoutes, ok := r.routes[method]
	if !ok {
		return nil, nil, false
	}

	for _, route := range methodRoutes {
		matched, params := matchRoute(route.pattern, path)
		if matched {
			return route.handler, params, true
		}
	}

	return nil, nil, false
}
