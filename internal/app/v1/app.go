package v1

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage/badgerdb"
	"github.com/valyala/fasthttp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type app struct {
	ctx            context.Context
	cancelFunc     context.CancelFunc
	config         config
	httpServer     fasthttp.Server
	telegramBotAPI *tgbotapi.BotAPI
	db             storage.Storage
}

func newApp(parentCtx context.Context, cfg config) (*app, error) {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
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

	// a.telegramBotAPI.Debug = true

	log.Debug(fmt.Sprintf("Authorized on telegram account %s", a.telegramBotAPI.Self.UserName))

	log.Info("setting up the database")
	if err := a.setupDatabase(); err != nil {
		log.Error("an error occurred when setting up the database", err)

		return nil, err
	}

	log.Info("setting up HTTP server")
	a.httpServer = fasthttp.Server{
		Handler:               a.getHTTPRequestHandler(),
		GetOnly:               false,
		NoDefaultServerHeader: true,
		NoDefaultDate:         true,
		MaxRequestBodySize:    1 * 1024 * 1024, // 1mb
		ReadBufferSize:        1024 * 1024,     // 1mb
	}

	return a, nil
}

func (a *app) setupDatabase() error {
	var err error

	switch a.config.Storage.Type {
	case storageTypeBadgerDB:
		a.db, err = badgerdb.New(a.ctx)
	default:
		return ErrUnsupportedStorageType
	}

	return err
}

func (a *app) start() error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "start",
	})

	defer a.cleanup()
	defer a.cancelFunc()

	var err error
	var wg sync.WaitGroup

	log.Info("opening the database connection")
	if err = a.db.Open(a.config.Storage.BadgerDB.DSN); err != nil {
		return errors.Wrap(ErrUnableToOpenDatabaseConnection, err.Error())
	}

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

		log.Info("starting the Telegram bot message handler")
		defer log.Info("Telegram bot message handler stopped")

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
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "stop",
	})

	log.Debug("cancelling app context")
	a.cancelFunc()
}

/*
func (a *app) healthCheck() error {
	return nil
}
*/

func (a *app) cleanup() {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "cleanup",
	})

	log.Info("cleanup started")
	defer log.Info("cleanup complete")

	log.Info("closing the database connection")
	if err := a.db.Close(); err != nil {
		log.Error("error when closing the database connection", err)
	}
}
