package v1

import (
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

func (a *app) telegramBotStopCommandHandler(chatID int64) error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotStopCommandHandler",
	})

	log.Debug("handling command")

	// define errMsg which is used to send a generic message to
	// the sender of the command via telegram when an error occurs
	// in the deferred function
	errMsg := ""
	defer func() {
		if errMsg != "" {
			u := types.User{TelegramChatID: chatID}
			if err := a.telegramBotSendMessage(u, errMsg); err != nil {
				log.Error("error when sending telegram error message", err)
			}
		}
	}()

	log.Debug("deleting all users by Telegram chat ID", chatID)
	err := a.db.GetUserRepositoryWriter().DeleteAllByTelegramChatID(chatID)
	if err != nil {
		errMsg = "an error occurred when trying to delete users"
		log.Error(errMsg, err)

		return err
	}

	return nil
}
