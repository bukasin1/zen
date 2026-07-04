package main

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/bukasin1/zen"
)

// Server contains the application's handlers.
type Server struct{}

// index serves the upload page.
func (s *Server) index(c *zen.Context) {
	// c.ServeFile(filepath.Join(staticDir, "index.html"))

	tmpl := template.Must(template.ParseFiles(filepath.Join(templatesDir, "index.html")))

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.ErrorInternalServer("Unable to read upload directory.")
		return
	}

	names := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	data := struct {
		Files []string
	}{
		Files: names,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.ErrorInternalServer("Unable to render page.")
	}
}

// upload handles multipart file uploads.
func (s *Server) upload(c *zen.Context) {
	_, file, err := c.FormFile("file")
	if err != nil {
		c.ErrorBadRequest("No file was uploaded.")
		return
	}

	dst := filepath.Join(uploadDir, file.Filename)
	// dst := filepath.Join("./uploadd", file.Filename)

	if _, err := c.SaveUploadedFile(file, dst); err != nil {
		c.ErrorInternalServer("Failed to save uploaded file.")
		return
	}

	c.SuccessCreated(map[string]any{
		"message":     "File uploaded successfully.",
		"file":        file.Filename,
		"size":        file.Size,
		"contentType": file.Header.Get("Content-Type"),
		"dstPath":     dst,
		"url":         "/uploads/" + file.Filename,
	})
}
