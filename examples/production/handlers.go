package main

import (
	"path/filepath"

	"github.com/bukasin1/zen"
)

// Server contains the application's handlers.
type Server struct{}

// Home serves the landing page.
func (s *Server) Home(c *zen.Context) {
	pathToFile := filepath.Join(staticDir, "index.html")
	c.ServeFile(pathToFile)
}

// Hello returns a simple JSON response.
func (s *Server) Hello(c *zen.Context) {
	c.SuccessOK(map[string]any{
		"message": "Welcome to the Zen production example.",
		"version": "v0.1.0",
		"status":  "running",
	})
}
