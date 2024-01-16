package v1

import (
	"os"
	"strconv"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types" //nolint:depguard
)

//nolint:funlen
func (a *app) telegramBotAddUserCommandHandler(chatID int64, arguments []string) error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotAddUserCommandHandler",
	})

	log.Debug("handling command")

	requestUser := types.User{TelegramChatID: chatID}

	errMsg := ""
	defer func() {
		if errMsg != "" {
			if err := a.telegramBotSendMessage(requestUser, errMsg); err != nil {
				log.Err(err).Error("error when sending telegram error message")
			}
		}
	}()

	log.Data("chatID", chatID).Debug("checking if the user is superadmin")
	if !a.telegramBotUserIsSuperUser(chatID) {
		err := ErrUnauthorizedToUseTelegramBotCommand
		errMsg = err.Error()
		log.Data("chatID", chatID).Err(err).Error(errMsg)

		return err
	}

	if len(arguments) < 1 {
		err := ErrInsufficientArguments
		errMsg = err.Error()

		log.Data("data", map[string]interface{}{
			"chatID":    chatID,
			"arguments": arguments,
		}).Err(err).Error(errMsg)

		return err
	}

	otherChatID, err := strconv.ParseInt(arguments[0], 10, 64)
	if err != nil {
		errMsg = "could not parse chat ID"

		log.Data("data", map[string]interface{}{
			"chatID":    chatID,
			"arguments": arguments,
		}).Err(err).Error(errMsg)

		return err //nolint:wrapcheck
	}

	newUser := types.User{
		ID:             generateUserID(),
		TelegramChatID: otherChatID,
	}

	if err := a.createUser(newUser); err != nil {
		errMsg = "error when creating user" //nolint:goconst
		log.Err(err).Error(errMsg)

		return err
	}

	return nil
}
