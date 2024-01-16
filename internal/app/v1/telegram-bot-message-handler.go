package v1

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5" //nolint:depguard
	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types" //nolint:depguard
)

type telegramBotCommand string

const (
	telegramBotStartCommand telegramBotCommand = "/start"
	telegramBotStopCommand  telegramBotCommand = "/stop"
	telegramBotGetAllUsers  telegramBotCommand = "/getAllUsers"
	telegramBotAddUser      telegramBotCommand = "/addUser"
)

// telegramBotMessageHandler is responsible for handling incoming messages
// from Telegram. It listens to a channel of updates and processes each
// message as they come in. If the message is a command (e.g. "/start" or "/stop"),
// it invokes the corresponding command handler function. If the message is
// not a command, it does nothing.
// This function should be run in a separate goroutine.
//
//nolint:funlen,gocognit,cyclop
func (a *app) telegramBotMessageHandler() error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "telegramBotMessageHandler",
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.telegramBotAPI.GetUpdatesChan(u)

	for {
		select {
		case <-a.ctx.Done():
			return a.ctx.Err() //nolint:wrapcheck
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			chatID := update.Message.Chat.ID
			data := update.Message.Text

			parts := strings.Fields(data)

			var command string
			var arguments []string //nolint:wsl

			if len(parts) > 0 {
				command = parts[0]

				if len(parts) > 1 {
					arguments = parts[1:]
				}
			}

			log.Debug(fmt.Sprintf("telegram message received: chat id: %d - username: %s - message: %s",
				chatID, update.Message.From.UserName, update.Message.Text))

			switch telegramBotCommand(command) {
			case telegramBotStartCommand:
				err := a.telegramBotStartCommandHandler(chatID)
				if err != nil {
					log.Err(err).Error("an error occurred when handling start the command")
				}
			case telegramBotStopCommand:
				err := a.telegramBotStopCommandHandler(chatID)
				if err != nil {
					log.Err(err).Error("an error occurred when handling the stop command")
				}
			case telegramBotGetAllUsers:
				err := a.telegramBotGetAllUsersCommandHandler(chatID)
				if err != nil {
					log.Err(err).Error("an error occurred when handling the get all users command")
				}
			case telegramBotAddUser:
				err := a.telegramBotAddUserCommandHandler(chatID, arguments)
				if err != nil {
					log.Err(err).Error("an error occurred when handling the add user command")
				}
			default:
			}
		}
	}
}

func (a *app) telegramBotSendMessage(user types.User, msg string) error {
	m := tgbotapi.NewMessage(user.TelegramChatID, msg)
	_, err := a.telegramBotAPI.Send(m)

	return err //nolint:wrapcheck
}

func (a *app) telegramBotUserIsSuperUser(chatID int64) bool {
	return chatID == a.config.TelegramBot.SuperuserChatID
}
