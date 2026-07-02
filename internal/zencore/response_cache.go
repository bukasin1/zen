package zencore

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
)

func GenerateETag(data []byte) string {

	hash := sha1.Sum(data)

	return `"` + hex.EncodeToString(hash[:]) + `"`
}

func (c *Context) SetETag(etag string) {
	c.SetHeader("ETag", etag)
}

func (c *Context) SetCacheControl(value string) {
	c.SetHeader("Cache-Control", value)
}

func (c *Context) IsETagMatch(etag string) bool {

	incoming := c.Header("If-None-Match")

	return incoming == etag
}

func (c *Context) NotModified() {
	_ = c.writeResponse(func() error {
		c.Status(http.StatusNotModified)
		return nil
	})

	// if err != nil {
	// 	c.app.logger.Error("failed to set not modified", logger.Fields{
	// 		"err":        err.Error(),
	// 		"request_id": c.RequestID(),
	// 	})
	// }
}
