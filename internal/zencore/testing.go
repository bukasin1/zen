package zencore

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// PerformTestRequest performs an HTTP request to the given handler.
// This is used for testing purposes.
func PerformTestRequest(
	handler http.Handler,
	method string,
	path string,
	body []byte,
	headers map[string]string,
) *httptest.ResponseRecorder {

	req := httptest.NewRequest(
		method,
		path,
		bytes.NewBuffer(body),
	)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	return recorder
}

// PerformTestJSONRequest performs an HTTP JSON request to the given handler.
// This is used for testing purposes.
func PerformTestJSONRequest(
	handler http.Handler,
	method string,
	path string,
	body any,
	headers map[string]string,
) *httptest.ResponseRecorder {

	var payload []byte

	if body != nil {
		payload, _ = json.Marshal(body)
	}

	if headers == nil {
		headers = make(map[string]string)
	}

	headers["Content-Type"] = "application/json"

	return PerformTestRequest(
		handler,
		method,
		path,
		payload,
		headers,
	)
}

func DecodeJSONResponse(
	rec *httptest.ResponseRecorder,
	target any,
) error {

	return json.Unmarshal(
		rec.Body.Bytes(),
		target,
	)
}

func DecodeJSONResponseAs[T any](
	rec *httptest.ResponseRecorder,
) (T, error) {

	var target T

	err := json.Unmarshal(
		rec.Body.Bytes(),
		&target,
	)

	return target, err
}

func NewTestContext(
	method string,
	path string,
	body []byte,
) (*Context, *httptest.ResponseRecorder) {

	req := httptest.NewRequest(
		method,
		path,
		bytes.NewBuffer(body),
	)

	rec := httptest.NewRecorder()

	app := New()

	ctx := NewContext(rec, req, app)

	return ctx, rec
}

func HasHeader(
	rec *httptest.ResponseRecorder,
	key string,
	expected string,
) bool {

	return rec.Header().Get(key) == expected
}

func HasStatus(
	rec *httptest.ResponseRecorder,
	status int,
) bool {

	return rec.Code == status
}

func ResponseBody(
	rec *httptest.ResponseRecorder,
) string {
	return rec.Body.String()
}
