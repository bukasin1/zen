package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	frameworkErrors "github.com/Danieljosh-uduma/zen/pkg/framework/errors"
	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/response"
	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/utils"
	"github.com/Danieljosh-uduma/zen/pkg/framework/internal/validator"
	"github.com/Danieljosh-uduma/zen/pkg/framework/logger"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	params map[string]string
	keys   map[string]interface{}

	rawBody []byte

	requestID string
	startTime time.Time

	responseCommitted atomic.Bool
	afterResponse     []func(c *Context)

	logger logger.Logger

	// user represents the authenticated user
	user any
}

func NewContext(w http.ResponseWriter, r *http.Request, logger logger.Logger) *Context {
	rw := &responseWriter{
		ResponseWriter: w,
		status:         0,
		// size:           0,
	}

	// ✅ generate request ID first
	requestID := utils.GenerateRequestID()

	// ✅ attach to standard Go context
	ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

	// ✅ replace request with new context-aware request
	r = r.WithContext(ctx)

	return &Context{
		Writer:  rw,
		Request: r,

		requestID: requestID,
		startTime: time.Now(),
		logger:    logger,
	}
}

// --------------- Response Writer Helpers ------------

// StatusCode returns the status code of the response.
func (c *Context) StatusCode() int {
	if rw, ok := c.Writer.(*responseWriter); ok {
		return rw.status
	}
	return 0
}

// ResponseSize returns the size of the response.
func (c *Context) ResponseSize() int {
	if rw, ok := c.Writer.(*responseWriter); ok {
		return rw.size
	}
	return 0
}

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

// AfterResponse registers a function hook to be called after the response is sent.
func (c *Context) AfterResponse(fn func(c *Context)) {
	c.afterResponse = append(c.afterResponse, fn)
}

// writeResponse is a helper function to write a response to the client.
func (c *Context) writeResponse(writeFn func() error) error {
	if !c.responseCommitted.CompareAndSwap(false, true) {
		return nil
	}

	// // mark as committed
	// c.responseCommitted = true

	err := writeFn()

	// run extended hooks AFTER write attempt
	c.runAfterResponseHooks()

	return err
}

// runAfterResponseHooks runs all registered after response hooks.
func (c *Context) runAfterResponseHooks() {
	for _, fn := range c.afterResponse {
		fn(c)
	}
}

// Redirect redirects the client to the given URL with the given status code.
func (c *Context) Redirect(status int, url string) {
	_ = c.writeResponse(func() error {
		http.Redirect(c.Writer, c.Request, url, status)
		return nil
	})
}

// NoContent writes a no content response to the client.
func (c *Context) NoContent() {
	_ = c.writeResponse(func() error {
		c.DelHeader("Content-Type")
		c.Status(http.StatusNoContent)
		return nil
	})
}

// Text writes a text response to the client.
func (c *Context) Text(status int, message string) error {
	return c.writeResponse(func() error {
		c.SetHeader("Content-Type", "text/plain")
		c.Status(status)
		_, err := c.Writer.Write([]byte(message))
		return err
	})
}

// JSON writes a JSON response to the client.
func (c *Context) JSON(status int, data any) error {
	err := c.writeResponse(func() error {
		c.SetHeader("Content-Type", "application/json")
		c.Status(status)
		return json.NewEncoder(c.Writer).Encode(data)
	})
	return err
}

// HTML writes an HTML response to the client.
func (c *Context) HTML(status int, html string) error {
	return c.writeResponse(func() error {
		c.SetHeader("Content-Type", "text/html")
		c.Status(status)
		_, err := c.Writer.Write([]byte(html))
		return err
	})
}

func (c *Context) Success(status int, data any) {
	resp := response.SuccessResponse{
		Success: true,
		Data:    data,
	}

	c.JSON(status, resp)
}

func (c *Context) SuccessWithMeta(status int, data any, meta any) {
	resp := response.SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}

	c.JSON(status, resp)
}

func (c *Context) Fail(status int, message string) {
	resp := response.ErrorResponse{
		Success: false,
		Error: response.ErrorDetail{
			Message: message,
		},
	}

	c.JSON(status, resp)
}

// Error writes an error response to the client.
func (c *Context) Error(status int, message string, code string, details any) {
	resp := response.ErrorResponse{
		Success: false,
		Error: response.ErrorDetail{
			Message: message,
			Code:    code,
			Details: details,
		},
	}

	c.JSON(status, resp)
}

// SuccessOK writes a success response with status code 200.
func (c *Context) SuccessOK(data any) {
	c.Success(http.StatusOK, data)
}

// SuccessCreated writes a success response with status code 201.
func (c *Context) SuccessCreated(data any) {
	c.Success(http.StatusCreated, data)
}

// BadRequest writes a bad request response with status code 400.
func (c *Context) BadRequest(message string) {
	c.Fail(http.StatusBadRequest, message)
}

// Unauthorized writes an unauthorized response with status code 401.
func (c *Context) Unauthorized(message string) {
	c.Fail(http.StatusUnauthorized, message)
}

// Forbidden writes a forbidden response with status code 403.
func (c *Context) Forbidden(message string) {
	c.Fail(http.StatusForbidden, message)
}

// NotFound writes a not found response with status code 404.
func (c *Context) NotFound(message string) {
	c.Fail(http.StatusNotFound, message)
}

// InternalError writes an internal server error response with status code 500.
func (c *Context) InternalError(message string) {
	c.Fail(http.StatusInternalServerError, message)
}

// ----------------- Response Writer Helpers End ----------

// ------------------ Request Helpers ---------------

// Param returns the value of the route parameter with the given key.
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Params returns all route parameters.
func (c *Context) Params() map[string]string {
	return c.params
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

// Body returns the raw request body as a byte slice.
//
// Currently overwrites any set limit from [MaxBodySize] middleware.
func (c *Context) Body() ([]byte, error) {
	if c.rawBody != nil {
		return c.rawBody, nil
	}

	if c.Request.Body == nil {
		return nil, nil
	}

	defer c.Request.Body.Close()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	c.rawBody = body

	// restore body for future reads
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

// BindJSON binds the request body to the given target.
// It returns an error if the request body is empty or invalid.
//
//	err := c.BindJSON(&payload)
//
//	if err != nil {
//		if appErr, ok := frameworkErrors.AsAppError(err); ok {
//			// structured handling
//			c.Error(appErr.Status, appErr.Message, appErr.Code, appErr.Details)
//			return
//		}
//
//		// fallback
//		c.Fail(400, err.Error())
//	}
func (c *Context) BindJSON(target any) error {
	if c.Request.Body == nil {
		return http.ErrBodyNotAllowed
	}

	defer c.Request.Body.Close()

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(target)
	if err != nil {
		c.LogError("bind json failed", logger.Fields{
			"error": err.Error(),
		})

		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			return frameworkErrors.New("Request body too large", http.StatusRequestEntityTooLarge)
		}

		return frameworkErrors.New("invalid JSON payload", http.StatusBadRequest)
	}

	// if decoder.More() {
	// 	return frameworkErrors.New("request body must contain only one JSON object", http.StatusBadRequest)
	// }

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return frameworkErrors.New("request body must contain only one JSON object", http.StatusBadRequest)
	}

	return nil
}

// MustBindJSON binds the request body to the given target.
// It panics if the request body is empty or invalid.
func (c *Context) MustBindJSON(target any) {
	if err := c.BindJSON(target); err != nil {
		panic(err)
	}
}

// BindAndValidate binds the request body to the given target and validates it.
// It returns an error if the request body is empty or invalid.
//
//	err := c.BindAndValidate(&payload)
//
//	if err != nil {
//		if appErr, ok := frameworkErrors.AsAppError(err); ok {
//			// structured handling
//			c.Error(appErr.Status, appErr.Message, appErr.Code, appErr.Details)
//			return
//		}
//
//		// fallback
//		c.Fail(400, err.Error())
//	}
func (c *Context) BindAndValidate(target any) error {
	if err := c.BindJSON(target); err != nil {
		return err
	}

	validationErrors := validator.ValidateStruct(target)
	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (c *Context) MustBindAndValidate(target any) {
	if err := c.BindAndValidate(target); err != nil {

		switch e := err.(type) {

		case validator.ValidationErrors:
			panic(
				frameworkErrors.WithDetails(
					frameworkErrors.WithCode(
						frameworkErrors.New("validation failed", 400),
						frameworkErrors.ErrValidation,
					),
					e,
				),
			)

		default:
			panic(err)
		}
	}
}

// ------------------- Request Helpers End ----------------------------

// ----------------- Context Helpers ----------------------------

// StdContext returns the request standard context
func (c *Context) StdContext() context.Context {
	return c.Request.Context()
}

// Set sets a value in the context.
func (c *Context) Set(key string, value interface{}) {
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
}

// Get returns a value from the context.
func (c *Context) Get(key string) (interface{}, bool) {
	if c.keys == nil {
		return nil, false
	}
	val, ok := c.keys[key]
	return val, ok
}

// MustGet returns a value from the context.
// It panics if the value is not found.
func (c *Context) MustGet(key string) interface{} {
	val, ok := c.Get(key)
	if !ok {
		panic("key not found: " + key)
	}
	return val
}

// RequestID returns the request ID.
func (c *Context) RequestID() string {
	return c.requestID
}

// StartTime returns the request start time.
func (c *Context) StartTime() time.Time {
	return c.startTime
}

// Duration returns the request duration.
func (c *Context) Duration() time.Duration {
	return time.Since(c.startTime)
}

func (c *Context) LogInfo(msg string, fields logger.Fields) {
	if c.logger == nil {
		return
	}
	c.logger.Info(msg, fields)
}

func (c *Context) LogError(msg string, fields logger.Fields) {
	if c.logger == nil {
		return
	}
	c.logger.Error(msg, fields)
}

func (c *Context) LogWarn(msg string, fields logger.Fields) {
	if c.logger == nil {
		return
	}
	c.logger.Warn(msg, fields)
}

func (c *Context) LogDebug(msg string, fields logger.Fields) {
	if c.logger == nil {
		return
	}
	c.logger.Debug(msg, fields)
}

// SetUser sets the authenticated user in the context.
// It should be called by an authentication middleware.
func (c *Context) SetUser(user any) {
	c.user = user
}

// User returns the authenticated user.
func (c *Context) User() any {
	return c.user
}

func (c *Context) UserOK() (any, bool) {
	if c.user == nil {
		return nil, false
	}
	return c.user, true
}

func (c *Context) MustUser() any {
	if c.user == nil {
		panic(frameworkErrors.New("user not found in context", 401))
	}
	return c.user
}

// ----------------- Context Helpers End ----------------------------
