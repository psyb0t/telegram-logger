package v1

import (
	"os"

	"github.com/psyb0t/glogger"
)

func (a *app) telegramBotStopCommandHandler(chatID int64) error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotStopCommandHandler",
	})

	log.Debug("handling command")

	log.Debug("finding user by telegram chat ID", chatID)
	user, err := a.db.GetUserRepositoryReader().FindByTelegramChatID(chatID)
	if err != nil {
		log.Error("an error occurred when trying to find a user by Telegram chat ID", err)

		return err
	}

	log.Debug("deleting user", user)
	if err = a.db.GetUserRepositoryWriter().Delete(user.ID); err != nil {
		log.Error("an error occurred when trying to delete user", err)

		return err
	}

	return err
}
