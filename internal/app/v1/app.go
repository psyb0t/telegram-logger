package v1

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/valyala/fasthttp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	logger "github.com/psyb0t/glogger"
)

type app struct {
	ctx            context.Context
	cancelFunc     context.CancelFunc
	config         config
	httpServer     fasthttp.Server
	telegramBotAPI *tgbotapi.BotAPI
}

func newApp(parentCtx context.Context, cfg config) (*app, error) {
	log := logger.New(logger.Caller{
		Service:  os.Getenv("SERVICENAME"),
		Package:  packageName,
		Function: "newApp",
	})

	ctx, cancelFunc := context.WithCancel(parentCtx)

	a := &app{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		config:     cfg,
	}

	log.Info("setting up the telegram bot connection")

	var err error
	a.telegramBotAPI, err = tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Error("an error occurred when setting up the telegram bot connection", err)

		return nil, err
	}

	// bot.Debug = true

	log.Debug(fmt.Sprintf("Authorized on telegram account %s", a.telegramBotAPI.Self.UserName))

	log.Info("setting up HTTP server")
	a.httpServer = fasthttp.Server{
		Handler:               a.getHTTPRequestHandler(),
		GetOnly:               false,
		NoDefaultServerHeader: true,
		NoDefaultDate:         true,
		MaxRequestBodySize:    50 * 1024 * 1024, // 50mb
		ReadBufferSize:        1024 * 1024,      // 1mb
	}

	return a, nil
}

func (a *app) start() error {
	log := logger.New(logger.Caller{
		Service:  os.Getenv("SERVICENAME"),
		Package:  packageName,
		Receiver: "app",
		Function: "start",
	})

	defer a.cleanup()

	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		defer wg.Done()

		log.Info("Starting HTTP server on " + a.config.ListenAddress)
		defer log.Info("HTTP server stopped")

		err := a.httpServer.ListenAndServe(a.config.ListenAddress)
		if errCh != nil {
			errCh <- err
			close(errCh)
			errCh = nil
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.telegramBotMessageHandler()
		if errCh != nil {
			errCh <- err
			close(errCh)
			errCh = nil
		}
	}()

	select {
	case <-a.ctx.Done():
		log.Info("context is done")

		log.Info("gracefully shutting down the HTTP sever")
		if err := a.httpServer.Shutdown(); err != nil {
			log.Info("HTTP server graceful shutdown failed:", err)
		}

		err = a.ctx.Err()
	case err = <-errCh:
		if err != nil {
			log.Info("HTTP server encountered an error:", err)
		}
	}

	wg.Wait()

	return err
}

func (a *app) stop() {
	a.cancelFunc()
}

func (a *app) cleanup() {
	log := logger.New(logger.Caller{
		Service:  os.Getenv("SERVICENAME"),
		Package:  packageName,
		Receiver: "app",
		Function: "cleanup",
	})

	log.Info("cleanup started")
	defer log.Info("cleanup complete")
}
