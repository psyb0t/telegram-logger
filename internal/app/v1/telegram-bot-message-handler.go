package v1

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/psyb0t/glogger"
)

type telegramBotCommand string

const (
	telegramBotStartCommand telegramBotCommand = "/start"
	telegramBotStopCommand  telegramBotCommand = "/stop"
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
				log.Debug(fmt.Sprintf("telegram message received: chat id: %d - username: %s - message: %s",
					update.Message.Chat.ID, update.Message.From.UserName, update.Message.Text))

				switch telegramBotCommand(update.Message.Text) {
				case telegramBotStartCommand:
					err := a.telegramBotHandleStartCommand(update.Message.Chat.ID)
					if err != nil {
						return err
					}
				case telegramBotStopCommand:
					err := a.telegramBotHandleStopCommand(update.Message.Chat.ID)
					if err != nil {
						return err
					}
				default:
				}
			}
		}
	}
}

func (a *app) telegramBotHandleStartCommand(chatID int64) error {
	/*
		log := glogger.New(glogger.Caller{
			Service: os.Getenv(serviceNameEnvVarName),
			Package:  packageName,
			Receiver: "app",
			Function: "telegramBotHandleStartCommand",
		})

		return err
	*/

	return nil
}

func (a *app) telegramBotHandleStopCommand(chatID int64) error {
	/*
		log := glogger.New(glogger.Caller{
			Service: os.Getenv(serviceNameEnvVarName),
			Package:  packageName,
			Receiver: "app",
			Function: "telegramBotHandleStopCommand",
		})

		return err
	*/

	return nil
}
