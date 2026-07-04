package main

import (
	"path/filepath"
	"runtime"

	"github.com/bukasin1/zen"
)

var _, filename, _, _ = runtime.Caller(0)

// Gets the directory of the current source file
var currentDir = filepath.Dir(filename)

var (
	uploadDir    = currentDir + "/uploads"
	staticDir    = currentDir + "/public"
	templatesDir = currentDir + "/templates"
)

func main() {
	app := zen.New()

	server := &Server{}

	// Make static paths wildcard with /* to capture every file within the static directory
	app.Static("/static/*", staticDir)
	app.Static("/uploads/*", uploadDir)

	app.Route("/").
		Get(server.index)

	app.Route("/upload").
		Post(server.upload)

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
