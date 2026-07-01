package framework

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter

	writer io.Writer
}

func (g *gzipResponseWriter) Write(
	b []byte,
) (int, error) {
	return g.writer.Write(b)
}

func GzipCompression() Middleware {
	return func(next HandlerFunc) HandlerFunc {

		return func(c *Context) {

			if !strings.Contains(
				c.Header("Accept-Encoding"),
				"gzip",
			) {
				next(c)
				return
			}

			if c.responseCommitted.Load() {
				next(c)
				return
			}

			gzipWriter := gzip.NewWriter(c.Writer)

			defer gzipWriter.Close()

			c.SetHeader("Content-Encoding", "gzip")
			c.AddHeader("Vary", "Accept-Encoding")

			originalWriter := c.Writer

			c.Writer = &gzipResponseWriter{
				ResponseWriter: originalWriter,
				writer:         gzipWriter,
			}

			next(c)

			c.Writer = originalWriter
		}
	}
}
