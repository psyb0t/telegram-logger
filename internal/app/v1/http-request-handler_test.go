package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestApp_getHTTPRequestHandler(t *testing.T) {
	// Initialize an instance of the app struct
	a := &app{}

	// Test the getHTTPRequestHandler method
	reqHandler := a.getHTTPRequestHandler()

	// Test that the returned request handler is of the correct type
	assert.IsType(t, fasthttp.RequestHandler(nil), reqHandler)
}

func TestApp_returnHTTPResponseString(t *testing.T) {
	// Initialize an instance of the app struct
	a := &app{}

	// Initialize a fasthttp request context
	ctx := &fasthttp.RequestCtx{}

	// Test data
	testCases := []struct {
		statusCode int
		body       string
		expected   string
	}{
		{fasthttp.StatusOK, "OK", "OK"},
		{fasthttp.StatusNotFound, "Not Found", "Not Found"},
		{fasthttp.StatusInternalServerError, "Internal Server Error", "Internal Server Error"},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call the returnHTTPResponseString method
		a.returnHTTPResponseString(ctx, tc.statusCode, tc.body)

		// Test that the status code and body are set correctly
		assert.Equal(t, tc.statusCode, ctx.Response.StatusCode())
		assert.Equal(t, tc.expected, string(ctx.Response.Body()))
	}
}

func TestApp_returnHTTPResponseJSON(t *testing.T) {
	// Initialize an instance of the app struct
	a := &app{}

	// Initialize a fasthttp request context
	ctx := &fasthttp.RequestCtx{}

	// Test data
	testCases := []struct {
		statusCode int
		data       interface{}
		expected   string
	}{
		{fasthttp.StatusOK, map[string]string{"message": "OK"}, `{"message":"OK"}`},
		{fasthttp.StatusNotFound, map[string]string{"error": "Not Found"}, `{"error":"Not Found"}`},
		{fasthttp.StatusInternalServerError, map[string]string{
			"error": "Internal Server Error",
		}, `{"error":"Internal Server Error"}`},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call the returnHTTPResponseJSON method
		a.returnHTTPResponseJSON(ctx, tc.statusCode, tc.data)
		// Test that the status code and body are set correctly
		assert.Equal(t, tc.statusCode, ctx.Response.StatusCode())
		assert.Equal(t, tc.expected, string(ctx.Response.Body()))

		// Test that the content type is set to "application/json"
		assert.Equal(t, contentTypeApplicationJSON, string(ctx.Response.Header.ContentType()))
	}
}

func TestApp_returnHTTPResponse(t *testing.T) {
	// Initialize an instance of the app struct
	a := &app{}

	// Initialize a fasthttp request context
	ctx := &fasthttp.RequestCtx{}

	// Test data
	testCases := []struct {
		statusCode  int
		contentType string
		body        []byte
		expected    string
	}{
		{fasthttp.StatusOK, contentTypeTextPlain, []byte("OK"), "OK"},
		{fasthttp.StatusNotFound, contentTypeTextPlain, []byte("Not Found"), "Not Found"},
		{fasthttp.StatusInternalServerError, contentTypeApplicationJSON, []byte(
			`{"error":"Internal Server Error"}`), `{"error":"Internal Server Error"}`},
	}

	// Loop over the test cases
	for _, tc := range testCases {
		// Call the returnHTTPResponse method
		a.returnHTTPResponse(ctx, tc.statusCode, tc.contentType, tc.body)

		// Test that the status code, content type, and body are set correctly
		assert.Equal(t, tc.statusCode, ctx.Response.StatusCode())
		assert.Equal(t, tc.contentType, string(ctx.Response.Header.ContentType()))
		assert.Equal(t, tc.expected, string(ctx.Response.Body()))
	}
}
