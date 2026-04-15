package framework

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

// --------------- Response Writer Helpers ------------

// Status sets the status code of the response.
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// SetHeader sets the header with the given key and value.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// AddHeader adds the header with the given key and value.
func (c *Context) AddHeader(key, value string) {
	c.Writer.Header().Add(key, value)
}

// DelHeader deletes the header with the given key.
func (c *Context) DelHeader(key string) {
	c.Writer.Header().Del(key)
}

// Text writes a text response to the client.
func (c *Context) Text(status int, message string) error {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(status)
	_, err := c.Writer.Write([]byte(message))
	return err
}

// JSON writes a JSON response to the client.
func (c *Context) JSON(status int, data any) error {
	c.SetHeader("Content-Type", "application/json")
	c.Status(status)
	return json.NewEncoder(c.Writer).Encode(data)
}

// HTML writes an HTML response to the client.
func (c *Context) HTML(status int, html string) error {
	c.SetHeader("Content-Type", "text/html")
	c.Status(status)
	_, err := c.Writer.Write([]byte(html))
	return err
}

// Error writes an error response to the client.
func (c *Context) Error(status int, message string) error {
	return c.JSON(status, map[string]any{
		"error": message,
	})
}

// Redirect redirects the client to the given URL with the given status code.
func (c *Context) Redirect(status int, url string) {
	http.Redirect(c.Writer, c.Request, url, status)
}

// NoContent writes a no content response to the client.
func (c *Context) NoContent() {
	c.DelHeader("Content-Type")
	c.Status(http.StatusNoContent)
}

// ----------------- Response Writer Helpers End ----------

// ------------------ Request Helpers ---------------

// Param returns the value of the route parameter with the given key.
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Query returns the value of the query parameter with the given key.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// QueryDefault returns the value of the query parameter with the given key.
// If the query parameter is not found, it returns the fallback value.
func (c *Context) QueryDefault(key, fallback string) string {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	return value
}

// QuerySlice returns the values of the query parameter with the given key.
// If the query parameter is not found, it returns an empty slice.
func (c *Context) QuerySlice(key string) []string {
	return c.Request.URL.Query()[key]
}

// Queries returns all query parameters.
func (c *Context) Queries() map[string][]string {
	return c.Request.URL.Query()
}

// Header returns the value of the header with the given key.
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// HasHeader returns true if the header with the given key exists.
func (c *Context) HasHeader(key string) bool {
	return c.Header(key) != ""
}

// HeaderValues returns all values of the header with the given key.
func (c *Context) HeaderValues(key string) []string {
	return c.Request.Header.Values(key)
}

// ------------------- Request Helpers End ----------------------------
