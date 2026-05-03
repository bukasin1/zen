package framework

import (
	"fmt"
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
	cleanPath := "/" + strings.Trim(path, "/")

	// Collapse multiple slashes: "/users//1" -> "/users/1"
	cleanPath = multiSlashRegex.ReplaceAllString(cleanPath, "/")

	// Remove trailing slash (except root)
	if len(cleanPath) > 1 && strings.HasSuffix(cleanPath, "/") {
		cleanPath = strings.TrimSuffix(cleanPath, "/")
	}

	// add back trailing slash if the original path had one
	if path[len(path)-1] == '/' && cleanPath[len(cleanPath)-1] != '/' {
		return cleanPath + "/"
	}

	return cleanPath
}

func getPathParts(path string) (string, []string) {
	// TODO: test more, it is confusing sometimes?
	// path = strings.Trim(path, "/")
	var pathParts []string
	if path == "" {
		pathParts = []string{}
	} else {
		pathParts = strings.Split(path, "/")
	}
	return path, pathParts
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

// HandlerFunc is a function that handles an HTTP request.
// It takes a *Context as an argument and returns nothing.
type HandlerFunc func(*Context)

type node struct {
	segment       string
	children      map[string]*node // static
	paramChild    *node            // :id
	wildcardChild *node            // * or *filepath

	handler     HandlerFunc
	wildcardKey string
	paramKeys   []string
}

type cacheKey struct {
	method string
	path   string
}

type cachedRoute struct {
	handler  HandlerFunc
	params   map[string]string
	redirect *redirectInfo
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

type redirectInfo struct {
	redirectPath string // path to redirect to
	code         int    // http status code
}

func (r *redirectInfo) isNil() bool {
	if r == nil {
		return true
	}

	if r.redirectPath == "" && r.code == 0 {
		return true
	}

	return false
}

func NewRouter() *Router {
	return &Router{
		routeTrees: make(map[string]*node),
		cache:      make(map[cacheKey]cachedRoute),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	// invalidate route cache on new route registration
	r.cache = make(map[cacheKey]cachedRoute)
	// normalize path to handle multiple slashes and trailing slashes
	validateRoutePath(path)
	path, pathParts := getPathParts(normalizeRoutePath(path))

	if _, ok := r.routeTrees[method]; !ok {
		r.routeTrees[method] = &node{
			segment: method,
		}
	}

	currentMethodNode := r.routeTrees[method]
	var paramKeys []string

	for i, part := range pathParts {
		if part == "" {
			continue
		}
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
				segment:     strings.Split(path, "*")[0],
				wildcardKey: strings.TrimPrefix(part, "*"),
			}

			// // set handler to base wildcard parent node
			// if currentMethodNode.handler == nil {
			// 	currentMethodNode.handler = handler
			// }

			currentMethodNode = currentMethodNode.wildcardChild
			break
		}

		// handle param (:)
		if strings.HasPrefix(part, ":") {
			paramKeys = append(paramKeys, part[1:])

			// check if param child exists
			if currentMethodNode.paramChild == nil {
				currentMethodNode.paramChild = &node{
					segment: part,
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
			currentMethodNode.children[part] = &node{
				segment: part,
			}
		}
		currentMethodNode = currentMethodNode.children[part]
	}

	// set handler
	if currentMethodNode.handler != nil {
		panic(fmt.Sprintf("duplicate route registration for %s %s %v", method, path))
	}
	currentMethodNode.handler = handler
	currentMethodNode.paramKeys = paramKeys
}

// processess wildcard node.
// merge wildcard key and remaining path to params.
func processWildcard(wildcardNode *node, path string, params *map[string]string) *node {
	strippedPath := strings.Split(path, wildcardNode.segment)
	if len(strippedPath) < 2 {
		return wildcardNode
	}

	remainingPath := strippedPath[1]

	wildcardKey := wildcardNode.wildcardKey
	if wildcardKey == "" {
		wildcardKey = "*"
	}
	if *params == nil {
		*params = make(map[string]string)
	}
	(*params)[wildcardKey] = remainingPath
	return wildcardNode
}

func localRedirect(w http.ResponseWriter, r *http.Request, redirect *redirectInfo) {
	if q := r.URL.RawQuery; q != "" {
		redirect.redirectPath += "?" + q
	}
	w.Header().Set("Location", redirect.redirectPath)
	w.WriteHeader(redirect.code)
}

func matchRouteTree(methodNode *node, path string) (HandlerFunc, map[string]string, bool, *redirectInfo) {
	if methodNode == nil {
		return nil, nil, false, nil
	}

	_, pathParts := getPathParts(path)

	var params map[string]string

	currentMethodNode := methodNode
	var paramValues []string
	seenWildcard := currentMethodNode.wildcardChild

	for i, part := range pathParts {
		if part == "" {
			// If part is empty and we are not at the first part, it means we have a trailing slash
			// We should check if there is a wildcard child to handle this case
			if i > 0 && seenWildcard == nil {
				return nil, nil, false, nil
			}
			continue
		}
		// 1.check if child exists for static part first (takes priority)
		if child, ok := currentMethodNode.children[part]; ok {
			currentMethodNode = child
			if currentMethodNode.wildcardChild != nil {
				// store the wildcard child for later use if needed
				seenWildcard = currentMethodNode.wildcardChild
			}
			continue
		}

		// 2. fallback to param child
		if currentMethodNode.paramChild != nil && part != "" {
			paramValues = append(paramValues, part)
			currentMethodNode = currentMethodNode.paramChild
			if currentMethodNode.wildcardChild != nil {
				// store the wildcard child for later use if needed
				seenWildcard = currentMethodNode.wildcardChild
			}
			continue
		}

		// 3. fallback to wildcard child
		if seenWildcard != nil {
			currentMethodNode = processWildcard(seenWildcard, path, &params)
			break
		}

		return nil, nil, false, nil
	}

	var redirect *redirectInfo
	if currentMethodNode.wildcardChild != nil && currentMethodNode.handler == nil && path[len(path)-1] != '/' {
		redirect = &redirectInfo{
			redirectPath: path + "/",
			code:         http.StatusTemporaryRedirect,
		}
	}

	// check if any wildcard node was encountered during traversal and the current node doesn't have a handler
	// this is the case for routes like /users and /users/*
	if seenWildcard != nil && seenWildcard != currentMethodNode && (currentMethodNode.handler == nil || path[len(path)-1] == '/') {
		currentMethodNode = processWildcard(seenWildcard, path, &params)
	}

	if currentMethodNode.handler == nil {
		return nil, nil, false, nil
	}

	for i, key := range currentMethodNode.paramKeys {
		if i < len(paramValues) {
			if params == nil {
				params = make(map[string]string)
			}
			params[key] = paramValues[i]
		}
	}

	return currentMethodNode.handler, params, true, redirect
}

func (r *Router) FindRoute(method, path string) (HandlerFunc, map[string]string, bool, *redirectInfo) {
	// normalize path to handle multiple slashes and trailing slashes
	normalizedPath := normalizeRoutePath(path)

	pathCacheKey := cacheKey{
		method: method,
		path:   normalizedPath,
	}

	// check cache first
	if cachedRoute, ok := r.cache[pathCacheKey]; ok {
		return cachedRoute.handler, cloneParams(cachedRoute.params), true, cachedRoute.redirect
	}

	methodNode, ok := r.routeTrees[method]
	if !ok {
		return nil, nil, false, nil
	}

	handler, params, ok, redirect := matchRouteTree(methodNode, normalizedPath)
	if !ok || handler == nil {
		return nil, nil, false, nil
	}

	// store in cache
	r.cache[pathCacheKey] = cachedRoute{
		handler:  handler,
		params:   cloneParams(params),
		redirect: redirect,
	}

	return handler, params, true, redirect
}

func (r *Router) ServeHTTP(ctx *Context) {
	if handler, params, ok, redirect := r.FindRoute(ctx.Request.Method, ctx.Request.URL.Path); ok {
		if !redirect.isNil() {
			localRedirect(ctx.Writer, ctx.Request, redirect)
			return
		}
		ctx.params = params
		handler(ctx)
		return
	}

	ctx.Fail(http.StatusNotFound, "404 page not found!")

	// http.NotFound(w, req)
}

// --------------------------- Deprecated Functions ---------------------------------

// Deprecated: meant for removal
// TODO: (for clean up) remove this function
func (r *Router) ServeHTTPOld(ctx *Context) {
	if handler, params, ok := r.FindRouteOld(ctx.Request.Method, ctx.Request.URL.Path); ok {
		ctx.params = params
		handler(ctx)
		return
	}

	for _, static := range r.staticRoutes {
		if matchStaticPrefix(ctx.Request.URL.Path, static.prefix) {
			static.handler.ServeHTTP(ctx.Writer, ctx.Request)
			return
		}
	}

	ctx.Fail(http.StatusNotFound, "404 page not found!")
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

// Deprecated: use Static instead
// TODO: remove this function
func (r *Router) HandleStatic(prefix string, handler http.Handler) {
	r.staticRoutes = append(r.staticRoutes, staticRoute{
		prefix:  prefix,
		handler: handler,
	})
}

// Deprecated: use Handle instead
// TODO: remove this function
func (r *Router) HandleOld(method, path string, handler HandlerFunc) {
	path = normalizeRoutePath(path)

	if r.routes == nil {
		r.routes = make(map[string][]route)
	}

	if _, ok := r.routes[method]; !ok {
		r.routes[method] = []route{}
	}
	r.routes[method] = append(r.routes[method], route{
		pattern: path,
		handler: handler,
	})
}
