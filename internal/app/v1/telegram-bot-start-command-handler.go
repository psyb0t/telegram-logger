package v1

import (
	"fmt"
	"os"

	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

const (
	telegramBotWelcomeMessageTpl = `Welcome!
Your ID is %s`
)

// telegramBotStartCommandHandler handles the telegramBotStartCommand command
// received from a user via a telegram bot. It generates a unique ID for the
// user and stores it in the database and it sends a welcome message to the
// user, containing the unique ID.
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

	log.Data("chatID", chatID).Debug("deleting all users matching the telegram chat id")
	err := a.db.GetUserRepositoryWriter().DeleteAllByTelegramChatID(chatID)
	if err != nil {
		errMsg = "error when cleaning up the database by telegram chat ID"
		log.Err(err).Error(errMsg)

		return err
	}

	log.Data("user", user).Debug("creating user")
	if err := a.db.GetUserRepositoryWriter().Create(user); err != nil {
		errMsg = "error when creating user"
		log.Err(err).Error(errMsg)

		return err
	}

	log.Data("user", user).Debug("sending welcome message to user")
	msg := fmt.Sprintf(telegramBotWelcomeMessageTpl, user.ID)
	if err := a.telegramBotSendMessage(user, msg); err != nil {
		errMsg = "could not send telegram welcome message"
		log.Err(err).Error(errMsg)

		return err
	}

	return nil
}
