package v1

import (
	"encoding/json"
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

func (a *app) telegramBotGetAllUsersCommandHandler(chatID int64) error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotGetAllUsersCommandHandler",
	})

	log.Debug("handling command")

	requestUser := types.User{TelegramChatID: chatID}

	// define errMsg which is used to send a generic message to
	// the sender of the command via telegram when an error occurs
	// in the deferred function
	errMsg := ""
	defer func() {
		if errMsg != "" {
			if err := a.telegramBotSendMessage(requestUser, errMsg); err != nil {
				log.Error("error when sending telegram error message", err)
			}
		}
	}()

	log.Debug("checking if the user is superadmin", chatID)
	if !a.telegramBotUserIsSuperUser(chatID) {
		err := ErrUnauthorizedToUseTelegramBotCommand
		errMsg = err.Error()
		log.Error(err, chatID)

		return err
	}

	log.Debug("getting all users")
	allUsers, err := a.db.GetUserRepositoryReader().GetAll()
	if err != nil {
		errMsg = "an error occurred when trying to get all users"
		log.Error(errMsg, err)

		return err
	}

	log.Debug("serializing the resulting users", allUsers)
	allUsersJSON, err := json.MarshalIndent(allUsers, "", " ")
	if err != nil {
		errMsg = "error when serializing the resulting users"
		log.Error(errMsg, err)

		return err
	}

	log.Debug("sending the response to the user")
	if err := a.telegramBotSendMessage(requestUser, string(allUsersJSON)); err != nil {
		errMsg = "could not send telegram message"
		log.Error(errMsg, err)

		return err
	}

	return nil
}
