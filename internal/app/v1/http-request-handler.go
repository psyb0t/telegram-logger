package v1

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func (a *app) getHTTPRequestHandler() fasthttp.RequestHandler {
	r := router.New()
	r.GET("/", a.rootHTTPHandler)

	requestHandler := r.Handler

	/*
		// last in first out
		requestHandler = middleware.LogResponse(requestHandler)
		requestHandler = middleware.LogRequest(requestHandler)

		if a.config.Tracer.Enabled {
			requestHandler = middleware.TraceID(requestHandler)
		}

		requestHandler = middleware.RequestID(requestHandler)
		requestHandler = middleware.RealIP(requestHandler)
		requestHandler = corsHandler.Handler(requestHandler)
	*/

	return requestHandler
}
