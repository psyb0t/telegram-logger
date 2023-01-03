package v1

import (
	"fmt"
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
	"github.com/valyala/fasthttp"
)

const (
	headerNameXID = "X-ID"
)

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

	/*
		// define errMsg which is used to send a generic message to
		// the sender of the command via telegram when an error occurs
		// in the deferred function
		errMsg := ""
		defer func() {
			if errMsg != "" {
				if err := a.telegramBotSendMessage(user, errMsg); err != nil {
					log.Error("error when sending telegram error message", err)
				}
			}
		}()
	*/

	msg := `you got this shit
	%s`

	reqMsg := string(ctx.Request.Body())

	if err := a.telegramBotSendMessage(user, fmt.Sprintf(msg, reqMsg)); err != nil {
		log.Err(err).Error("there was an error when sending the message to the user")
	}

	a.returnHTTPResponseJSON(ctx, 200, nil)
}
