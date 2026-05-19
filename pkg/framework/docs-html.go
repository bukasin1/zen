package framework

import (
	"html/template"
	"net/http"
)

const routeDocsHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>API Routes</title>

<style>
body {
	font-family: Arial, sans-serif;
	padding: 40px;
	background: #f5f5f5;
}

.route {
	background: white;
	padding: 20px;
	margin-bottom: 20px;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0,0,0,0.08);
}

.method {
	font-weight: bold;
}

.path {
	font-family: monospace;
	font-size: 16px;
}

.tag {
	display: inline-block;
	padding: 4px 8px;
	margin-right: 6px;
	background: #eee;
	border-radius: 4px;
	font-size: 12px;
}

.deprecated {
	color: red;
	font-weight: bold;
}
</style>
</head>

<body>

<h1>API Routes</h1>

{{range .}}

<div class="route">

<div class="method">
{{.Method}}
</div>

<div class="path">
{{.Path}}
</div>

{{if .Name}}
<div>
<strong>Name:</strong> {{.Name}}
</div>
{{end}}

{{if .Summary}}
<div>
<strong>Summary:</strong> {{.Summary}}
</div>
{{end}}

{{if .Description}}
<div>
<strong>Description:</strong> {{.Description}}
</div>
{{end}}

{{if .Version}}
<div>
<strong>Version:</strong> {{.Version}}
</div>
{{end}}

{{if .OperationID}}
<div>
<strong>Operation ID:</strong> {{.OperationID}}
</div>
{{end}}

{{if .Tags}}
<div>
{{range .Tags}}
<span class="tag">{{.}}</span>
{{end}}
</div>
{{end}}

{{if .Middlewares}}
<div>
<strong>Middlewares:</strong>

{{range .Middlewares}}
<span class="tag">{{.}}</span>
{{end}}
</div>
{{end}}

{{if .Deprecated}}
<div class="deprecated">
DEPRECATED
</div>
{{end}}

</div>

{{end}}

</body>
</html>
`

// MountHTMLDocs mounts a lightweight HTML docs page.
func (a *App) MountHTMLDocs(
	path string,
	options ...RouteDocOptions,
) {
	tmpl := template.Must(
		template.New("route-docs").
			Parse(routeDocsHTMLTemplate),
	)

	a.Route(path).
		Internal().
		Get(func(ctx *Context) {
			ctx.Writer.Header().Set(
				"Content-Type",
				"text/html; charset=utf-8",
			)

			err := tmpl.Execute(
				ctx.Writer,
				a.RouteDocs(options...),
			)

			if err != nil {
				http.Error(
					ctx.Writer,
					"failed to render docs",
					http.StatusInternalServerError,
				)

				return
			}
		})
}
