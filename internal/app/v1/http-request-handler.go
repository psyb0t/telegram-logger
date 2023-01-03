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

func (a *app) getHTTPRequestHandler() fasthttp.RequestHandler {
	r := router.New()
	r.POST("/", a.rootHTTPHandler)

	return r.Handler
}

func (a *app) returnHTTPResponseString(ctx *fasthttp.RequestCtx, statusCode int, body string) {
	a.returnHTTPResponse(ctx, statusCode, contentTypeTextPlain, []byte(body))
}

func (a *app) returnHTTPResponseJSON(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
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

func (a *app) returnHTTPResponse(ctx *fasthttp.RequestCtx,
	statusCode int, contentType string, body []byte) {

	ctx.SetContentType(contentType)
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(body)
}
