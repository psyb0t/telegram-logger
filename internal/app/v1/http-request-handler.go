package v1

import (
	"encoding/json"
	"os"

	"github.com/fasthttp/router"
	"github.com/psyb0t/glogger"
	"github.com/valyala/fasthttp"
)

const (
	contentTypeTextPlain       = "text/plain"
	contentTypeApplicationJSON = "application/json"
)

// getHTTPRequestHandler returns an HTTP request handler for the app.
// It handles requests by delegating to the appropriate handler function
// based on the request method and path.
func (a *app) getHTTPRequestHandler() fasthttp.RequestHandler {
	r := router.New()
	r.POST("/", a.rootHTTPHandler)

	return r.Handler
}

// returnHTTPResponseString returns an HTTP response with the provided
// status code and body as a string.
func (a *app) returnHTTPResponseString(
	ctx *fasthttp.RequestCtx, statusCode int, body string,
) {
	a.returnHTTPResponse(ctx, statusCode, contentTypeTextPlain, []byte(body))
}

// returnHTTPResponseJSON returns an HTTP response with the provided
// status code, and body serialized as JSON.
func (a *app) returnHTTPResponseJSON(
	ctx *fasthttp.RequestCtx, statusCode int, data interface{},
) {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "returnHTTPResponseJSON",
	})

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Error("could json.Marshal data")

		a.returnHTTPResponseString(ctx, fasthttp.StatusInternalServerError,
			fasthttp.StatusMessage(fasthttp.StatusInternalServerError))

		return
	}

	a.returnHTTPResponse(ctx, statusCode, contentTypeApplicationJSON, jsonData)
}

// returnHTTPResponse returns an HTTP response with the provided status code,
// content type, and body.
func (a *app) returnHTTPResponse(ctx *fasthttp.RequestCtx,
	statusCode int, contentType string, body []byte,
) {
	ctx.SetContentType(contentType)
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(body)
}
