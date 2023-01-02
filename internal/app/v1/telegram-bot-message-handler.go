package v1

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/types"
)

type telegramBotCommand string

const (
	telegramBotStartCommand telegramBotCommand = "/start"
	telegramBotStopCommand  telegramBotCommand = "/stop"
	telegramBotGetAllUsers  telegramBotCommand = "/getAllUsers"
)

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
			return a.ctx.Err()
		case update := <-updates:
			if update.Message != nil {
				chatID := update.Message.Chat.ID

				log.Debug(fmt.Sprintf("telegram message received: chat id: %d - username: %s - message: %s",
					chatID, update.Message.From.UserName, update.Message.Text))

				switch telegramBotCommand(update.Message.Text) {
				case telegramBotStartCommand:
					err := a.telegramBotStartCommandHandler(chatID)
					if err != nil {
						log.Error("an error occurred when handling start the command", err)
					}
				case telegramBotStopCommand:
					err := a.telegramBotStopCommandHandler(chatID)
					if err != nil {
						log.Error("an error occurred when handling the stop command", err)
					}
				case telegramBotGetAllUsers:
					err := a.telegramBotGetAllUsersCommandHandler(chatID)
					if err != nil {
						log.Error("an error occurred when handling the get all users command", err)
					}
				default:
				}
			}
		}
	}
}

func (a *app) telegramBotSendMessage(user types.User, msg string) error {
	m := tgbotapi.NewMessage(user.TelegramChatID, msg)
	_, err := a.telegramBotAPI.Send(m)

	return err
}

func (a *app) telegramBotUserIsSuperUser(chatID int64) bool {
	return chatID == a.config.TelegramBot.SuperuserChatID
}
