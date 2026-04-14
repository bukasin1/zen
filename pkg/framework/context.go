package framework

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

func (c *Context) Text(status int, message string) {
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(message))
}

func (c *Context) JSON(status int, data any) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	return json.NewEncoder(c.Writer).Encode(data)
}
