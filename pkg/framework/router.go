package framework

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var multiSlashRegex = regexp.MustCompile(`/+`)

func normalizeRoutePath(path string) string {
	if path == "" {
		return "/"
	}

	// trim out leading and trailing spaces
	path = strings.TrimSpace(path)

	// collapse repeated slashes
	path = multiSlashRegex.ReplaceAllString(path, "/")

	// preserve trailing slash info
	hasTrailingSlash := len(path) > 1 && strings.HasSuffix(path, "/")

	// trim outer slashes and ensure leading slash
	path = "/" + strings.Trim(path, "/")

	// restore trailing slash if needed
	if hasTrailingSlash {
		path += "/"
	}

	return path
}

type parsedPath struct {
	path             string
	parts            []string
	hasTrailingSlash bool
}

func parsePath(path string) parsedPath {
	path = normalizeRoutePath(path)

	hasTrailingSlash := len(path) > 1 && strings.HasSuffix(path, "/")

	trimmed := strings.Trim(path, "/")

	var parts []string
	if trimmed != "" {
		parts = strings.Split(trimmed, "/")
	}

	return parsedPath{
		path:             path,
		parts:            parts,
		hasTrailingSlash: hasTrailingSlash,
	}
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
	path        string
	wildcardKey string
	paramKeys   []string
}

type cacheKey struct {
	method string
	path   string
}

type routeCache struct {
	routeNode *node
	params    map[string]string
	redirect  *redirectInfo
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
	cache      map[cacheKey]routeCache
	cacheMu    sync.RWMutex
	routeCount int

	// TODO: cleanup(remove old routes registering)
	routes       map[string][]route
	staticRoutes []staticRoute
}

type redirectInfo struct {
	redirectPath string // path to redirect to
	code         int    // http status code
}

func NewRouter() *Router {
	return &Router{
		routeTrees: make(map[string]*node),
		cache:      make(map[cacheKey]routeCache),
	}
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	// invalidate route cache on new route registration
	r.cacheMu.Lock()
	r.cache = make(map[cacheKey]routeCache)
	r.cacheMu.Unlock()
	// normalize path to handle multiple slashes and trailing slashes
	validateRoutePath(path)
	parsed := parsePath(path)

	path = parsed.path
	pathParts := parsed.parts

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
		panic(fmt.Sprintf("duplicate route registration for %s %s. handler already registered with %s", method, path, currentMethodNode.path))
	}
	currentMethodNode.handler = handler
	currentMethodNode.path = path
	currentMethodNode.paramKeys = paramKeys
	r.routeCount++
}

// processess wildcard node.
// merge wildcard key and remaining path to params.
func processWildcardKeyParams(wildcardNode *node, parsed parsedPath, wildcardIndex int, params *map[string]string) *node {
	remainingParts := parsed.parts[wildcardIndex:]

	wildcardValue := strings.Join(remainingParts, "/")

	if parsed.hasTrailingSlash && wildcardValue != "" {
		wildcardValue += "/"
	}

	key := wildcardNode.wildcardKey
	if key == "" {
		key = "*"
	}

	if *params == nil {
		*params = make(map[string]string)
	}

	(*params)[key] = wildcardValue
	return wildcardNode
}

func localRedirect(w http.ResponseWriter, r *http.Request, redirect *redirectInfo) {
	if q := r.URL.RawQuery; q != "" {
		redirect.redirectPath += "?" + q
	}
	w.Header().Set("Location", redirect.redirectPath)
	w.WriteHeader(redirect.code)
}

func matchRouteTree(methodNode *node, path string) (*node, map[string]string, bool, *redirectInfo) {
	if methodNode == nil {
		return nil, nil, false, nil
	}

	parsed := parsePath(path)

	pathParts := parsed.parts
	hasTrailingSlash := parsed.hasTrailingSlash

	var params map[string]string

	currentMethodNode := methodNode
	var paramValues []string

	var wildcardNode *node
	var wildcardIndex int

	for i, part := range pathParts {
		// 1.check if child exists for static part first (takes priority)
		if child, ok := currentMethodNode.children[part]; ok {
			currentMethodNode = child
			if currentMethodNode.wildcardChild != nil {
				// store the wildcard child for later use if needed
				wildcardNode = currentMethodNode.wildcardChild
				wildcardIndex = i + 1
			}
			continue
		}

		// 2. fallback to param child
		if currentMethodNode.paramChild != nil && part != "" {
			paramValues = append(paramValues, part)
			currentMethodNode = currentMethodNode.paramChild
			if currentMethodNode.wildcardChild != nil {
				// store the wildcard child for later use if needed
				wildcardNode = currentMethodNode.wildcardChild
				wildcardIndex = i + 1
			}
			continue
		}

		// 3. fallback to wildcard child
		if wildcardNode != nil {
			currentMethodNode = processWildcardKeyParams(wildcardNode, parsed, wildcardIndex, &params)
			break
		}

		return nil, nil, false, nil
	}

	if hasTrailingSlash && wildcardNode == nil {
		return nil, nil, false, nil
	}

	if currentMethodNode.wildcardChild != nil && currentMethodNode.handler == nil && !hasTrailingSlash {
		return currentMethodNode.wildcardChild,
			params,
			true,
			&redirectInfo{
				redirectPath: parsed.path + "/",
				code:         http.StatusTemporaryRedirect,
			}
	}

	// check if any wildcard node was encountered during traversal and the current node doesn't have a handler
	// this is the case for routes like /users and /users/*
	if currentMethodNode.handler == nil || hasTrailingSlash {
		if wildcardNode != nil && wildcardNode != currentMethodNode {
			currentMethodNode = processWildcardKeyParams(wildcardNode, parsed, wildcardIndex, &params)
		}
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

	return currentMethodNode, params, true, nil
}

func (r *Router) FindRoute(method, path string) (*node, map[string]string, bool, *redirectInfo) {
	// normalize path to handle multiple slashes and trailing slashes
	normalizedPath := normalizeRoutePath(path)

	pathCacheKey := cacheKey{
		method: method,
		path:   normalizedPath,
	}

	// check cache first
	r.cacheMu.RLock()
	cachedRoute, ok := r.cache[pathCacheKey]
	r.cacheMu.RUnlock()

	if ok {
		// if normalized path is different from the original path, redirect
		// This handles cases where the original path has multiple slashes
		if len(normalizedPath) != len(path) && cachedRoute.redirect == nil {
			return cachedRoute.routeNode, cloneParams(cachedRoute.params), true, &redirectInfo{
				redirectPath: normalizedPath,
				code:         http.StatusMovedPermanently,
			}
		}

		return cachedRoute.routeNode, cloneParams(cachedRoute.params), true, cachedRoute.redirect
	}

	methodNode, ok := r.routeTrees[method]
	if !ok {
		return nil, nil, false, nil
	}

	routeNode, params, ok, redirect := matchRouteTree(methodNode, normalizedPath)
	if !ok || routeNode == nil {
		return nil, nil, false, nil
	}

	// store in cache
	r.cacheMu.Lock()
	r.cache[pathCacheKey] = routeCache{
		routeNode: routeNode,
		params:    cloneParams(params),
		redirect:  redirect,
	}
	r.cacheMu.Unlock()

	// if normalized path is different from the original path, redirect
	// This handles cases where the original path has multiple slashes
	if len(normalizedPath) != len(path) && redirect == nil {
		return routeNode, params, true, &redirectInfo{
			redirectPath: normalizedPath,
			code:         http.StatusMovedPermanently,
		}
	}

	return routeNode, params, true, redirect
}

func (r *Router) ServeHTTP(ctx *Context) {
	if routeNode, params, ok, redirect := r.FindRoute(ctx.Request.Method, ctx.Request.URL.Path); ok {
		if redirect != nil {
			localRedirect(ctx.Writer, ctx.Request, redirect)
			return
		}
		ctx.params = params
		ctx.fullPath = routeNode.path
		routeNode.handler(ctx)
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
