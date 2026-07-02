package zencore

import "net/http"

// MountJSONDocs mounts a JSON documentation endpoint.
func (a *App) MountJSONDocs(path string, opts ...RouteDocOptions) {
	a.Route(path).
		Name("docs.json").
		Internal().
		Get(func(ctx *Context) {
			ctx.JSON(http.StatusOK, a.RouteDocs(opts...))
		})
}
