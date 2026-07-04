package zencore

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	// DefaultMultipartMemory is the default maximum memory for multipart forms.
	DefaultMultipartMemory = 32 << 20 // 32 MB
)

// ParseMultipartForm parses the multipart form.
// If maxMemory is less than or equal to 0, it uses [DefaultMultipartMemory].
func (c *Context) ParseMultipartForm(
	maxMemory int64,
) error {

	if c.Request.MultipartForm != nil {
		return nil
	}

	if maxMemory <= 0 {
		maxMemory = DefaultMultipartMemory
	}

	err := c.Request.ParseMultipartForm(maxMemory)
	if err != nil {
		return err
	}

	return nil
}

// FormValue returns the value of the first form key.
func (c *Context) FormValue(
	key string,
) string {
	return c.Request.FormValue(key)
}

// MultipartFormValues returns the values of the first form key.
func (c *Context) MultipartFormValues(
	key string,
) []string {

	if c.Request.MultipartForm == nil {
		return nil
	}

	return c.Request.MultipartForm.Value[key]
}

// FormFile returns the first file header for the given key.
func (c *Context) FormFile(
	key string,
) (
	multipart.File,
	*multipart.FileHeader,
	error,
) {
	return c.Request.FormFile(key)
}

// MultipartFiles returns the files for the given key.
func (c *Context) MultipartFiles(
	key string,
) ([]*multipart.FileHeader, error) {

	if c.Request.MultipartForm == nil {
		return nil, errors.New(
			"multipart form not parsed",
		)
	}

	files := c.Request.MultipartForm.File[key]

	return files, nil
}

// RemoveMultipartFiles removes the multipart files.
func (c *Context) RemoveMultipartFiles() error {

	if c.Request.MultipartForm == nil {
		return nil
	}

	return c.Request.MultipartForm.RemoveAll()
}

func (c *Context) SaveUploadedFile(
	fileHeader *multipart.FileHeader,
	dst string,
) (int64, error) {

	src, err := fileHeader.Open()
	if err != nil {
		return 0, err
	}

	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer out.Close()

	return io.Copy(out, src)
}

// ServeFile sends the file to the client.
func (c *Context) ServeFile(name string) {
	http.ServeFile(c.Writer, c.Request, name)
}
