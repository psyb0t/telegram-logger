package v1

import (
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

func (a *app) telegramBotStartCommandHandler(chatID int64) error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotStartCommandHandler",
	})

	log.Debug("handling command")

	user := types.User{
		ID:             generateUserID(),
		TelegramChatID: chatID,
	}

	log.Debug("creating user", user)

	// remove any other existing users with the given chatID

	if err := a.db.GetUserRepositoryWriter().Create(user); err != nil {
		log.Error("error when creating user", err)

		return err
	}

	if err := a.telegramBotSendWelcomeMessage(user); err != nil {
		log.Error("could not send telegram message", err)

		return err
	}

	return nil
}
