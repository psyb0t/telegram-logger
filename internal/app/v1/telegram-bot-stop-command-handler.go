package v1

import (
	"fmt"
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

const (
	telegramBotByeMessageTpl = `Bye!
Chat ID %d has been removed from the system.`
)

// telegramBotStopCommandHandler handles the telegramBotStopCommand command
// of the Telegram bot. When this function is called, it removes all users
// with the provided chat ID from the system and sends notification to the user
// via Telegram.
func (a *app) telegramBotStopCommandHandler(chatID int64) error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotStopCommandHandler",
	})

	log.Debug("handling command")

	user := types.User{TelegramChatID: chatID}

	// define errMsg which is used to send a generic message to
	// the sender of the command via telegram when an error occurs
	// in the deferred function
	errMsg := ""
	defer func() {
		if errMsg != "" {
			if err := a.telegramBotSendMessage(user, errMsg); err != nil {
				log.Err(err).Error("error when sending telegram error message")
			}
		}
	}()

	log.Data("chatID", chatID).Debug("deleting all users by Telegram chat ID")
	err := a.db.GetUserRepositoryWriter().DeleteAllByTelegramChatID(chatID)
	if err != nil {
		errMsg = "an error occurred when trying to delete users"
		log.Err(err).Error(errMsg)

		return err
	}

	msg := fmt.Sprintf(telegramBotByeMessageTpl, chatID)
	if err := a.telegramBotSendMessage(user, msg); err != nil {
		errMsg = "could not send telegram bye message"
		log.Err(err).Error(errMsg)

		return err
	}

	return nil
}
