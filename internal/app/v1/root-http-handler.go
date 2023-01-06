package v1

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
	"github.com/psyb0t/telegram-logger/pkg/types"
	"github.com/valyala/fasthttp"
)

const (
	headerNameXID = "X-ID"
)

// rootHTTPHandler handles HTTP requests to the root path. It gets the user
// associated with the request based on the value of the X-ID header, parses
// the JSON request body, builds a Telegram message string from the request,
// and sends the message to the user via the Telegram bot. It returns an HTTP
// response with the status code, header, and body serialized as JSON.
func (a *app) rootHTTPHandler(ctx *fasthttp.RequestCtx) {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "rootHTTPHandler",
	})

	id := string(ctx.Request.Header.Peek(headerNameXID))
	log.Data("id", id).Debug("getting user by id")
	user, err := a.db.GetUserRepositoryReader().Get(id)
	if err != nil {
		log.Err(err).Error("there was an error when getting the user by ID")

		if err == storage.ErrEmptyID || err == storage.ErrNotFound {
			a.returnHTTPResponseString(ctx, fasthttp.StatusUnauthorized, "")

			return
		}
	}

	log.Debug("parsing JSON request body")
	request := types.Request{}
	if err := json.Unmarshal(ctx.Request.Body(), &request); err != nil {
		log.Err(err).Error("could not parse JSON request body")

		a.returnHTTPResponseJSON(ctx, fasthttp.StatusBadRequest,
			types.Response{Error: err.Error()})

		return
	}

	log.Data("request", request).Debug("building telegram message string from request")
	telegramMessage, err := requestToTelegramMessageString(request)
	if err != nil {
		log.Err(err).Error("an error occurred when building telegram message string from request")

		a.returnHTTPResponseJSON(ctx, fasthttp.StatusInternalServerError,
			types.Response{Error: err.Error()})

		return
	}

	log.Data("telegramMessage", telegramMessage).
		Data("user", user).
		Debug("sending message to the user")

	if err := a.telegramBotSendMessage(user, telegramMessage); err != nil {
		log.Err(err).Error("there was an error when sending the message to the user")

		a.returnHTTPResponseJSON(ctx, fasthttp.StatusInternalServerError,
			types.Response{Error: err.Error()})

		return
	}

	response := types.Response{Message: "successfully sent log entry via Telegram"}
	a.returnHTTPResponseJSON(ctx, fasthttp.StatusOK, response)
}

// requestToTelegramMessageString builds a Telegram message string from a
// types.Request struct. It returns the message string and an error if
// there was an issue building the string.
func requestToTelegramMessageString(request types.Request) (string, error) {
	telegramMsg := ""
	if request.Caller != "" {
		telegramMsg += fmt.Sprintf("Caller: %s\n", request.Caller)
	}

	if request.Time != "" {
		telegramMsg += fmt.Sprintf("Time: %s\n", request.Time)
	}

	if request.Level != "" {
		telegramMsg += fmt.Sprintf("Level: %s\n", request.Level)
	}

	if request.RequestID != "" {
		telegramMsg += fmt.Sprintf("RequestID: %s\n", request.RequestID)
	}

	if request.TraceID != "" {
		telegramMsg += fmt.Sprintf("TraceID: %s\n", request.TraceID)
	}

	if request.SpanID != "" {
		telegramMsg += fmt.Sprintf("SpanID: %s\n", request.SpanID)
	}

	if request.Message != "" {
		telegramMsg += fmt.Sprintf("Message: %s\n", request.Message)
	}

	if request.Data != nil {
		serializedData, err := json.Marshal(request.Data)
		if err != nil {
			return "", err
		}

		telegramMsg += fmt.Sprintf("Data: %s\n", serializedData)
	}

	return telegramMsg, nil
}
