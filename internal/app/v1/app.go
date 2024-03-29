package v1

import (
	"context"
	"fmt"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/psyb0t/glogger"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage"
	"github.com/psyb0t/telegram-logger/internal/pkg/storage/badgerdb"
	"github.com/valyala/fasthttp"
)

// app contains the context, cancel function, config, HTTP server,
// Telegram bot API, and database connection for the app.
type app struct {
	ctx            context.Context //nolint:containedctx
	cancelFunc     context.CancelFunc
	config         config
	httpServer     fasthttp.Server
	telegramBotAPI *tgbotapi.BotAPI
	db             storage.Storage
}

// newApp creates a new app struct and initializes the Telegram
// bot connection and database. It also sets up the HTTP server.
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
	a.telegramBotAPI, err = tgbotapi.NewBotAPI(cfg.TelegramBot.Token)
	if err != nil {
		log.Err(err).Error("an error occurred when setting up the telegram bot connection")

		return nil, err
	}

	// a.telegramBotAPI.Debug = true

	log.Debug(fmt.Sprintf("Authorized on telegram account %s", a.telegramBotAPI.Self.UserName))

	log.Info("setting up the database")
	if err := a.setupDatabase(); err != nil {
		log.Err(err).Error("an error occurred when setting up the database")

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

// sets up the database connection for the app based on the specified
// storage type in the config.
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

// start starts the app by opening the database connection and starting the
// HTTP server and Telegram bot message handler in separate goroutines.
// It waits for either the context to be cancelled or for one of the goroutines
// to return an error. If the context is cancelled, it sets the error to the
// context's error. If one of the goroutines returns an error, it sets the
// error to that error. It then stops the app and returns the error.
func (a *app) start() error {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "start",
	})

	log.Info("opening the database connection")
	if err := a.db.Open(a.config.Storage.BadgerDB.DSN); err != nil {
		log.Err(err).Error(ErrUnableToOpenDatabaseConnection.Error())

		return ErrUnableToOpenDatabaseConnection
	}

	var wg sync.WaitGroup

	httpServerErrCh := make(chan error, 1)
	wg.Add(1)
	go a.startHTTPServer(&wg, httpServerErrCh)

	telegramBotMessageHandlerErrCh := make(chan error, 1)
	wg.Add(1)
	go a.startTelegramBotMessageHandler(&wg, telegramBotMessageHandlerErrCh)

	var err error
	select {
	case <-a.ctx.Done():
		log.Info("context is done")
		err = a.ctx.Err()
	case err = <-httpServerErrCh:
		if err != nil {
			log.Err(err).Error("HTTP server encountered an error")
		}
	case err = <-telegramBotMessageHandlerErrCh:
		if err != nil {
			log.Err(err).Error("Telegram bot message handler encountered an error")
		}
	}

	a.cancelFunc()
	a.cleanup()

	log.Debug("waiting for wait group to de done")
	wg.Wait()

	return err
}

// startHTTPServer starts the HTTP server and waits for it to stop. If the
// server returns an error, it is passed on to the calling function via
// the provided error channel.
func (a *app) startHTTPServer(wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	defer close(errCh)

	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "startHTTPServer",
	})

	log.Info("Starting HTTP server on " + a.config.ListenAddress)
	defer log.Info("HTTP server stopped")

	errCh <- a.httpServer.ListenAndServe(a.config.ListenAddress)
}

// startTelegramBotMessageHandler starts the Telegram bot message handler
// and waits for it to stop. If the message handler returns an error, it is
// passed on to the calling function via the provided error channel.
func (a *app) startTelegramBotMessageHandler(wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	defer close(errCh)

	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "startHTTPServer",
	})

	log.Info("starting the Telegram bot message handler")
	defer log.Info("Telegram bot message handler stopped")

	errCh <- a.telegramBotMessageHandler()
}

func (a *app) stop() {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "stop",
	})

	log.Info("cancelling app context")
	a.cancelFunc()
}

/*
func (a *app) healthCheck() error {
	return nil
}
*/

// cleanup gracefully shuts down the HTTP server and closes the
// database connection.
func (a *app) cleanup() {
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Receiver: "app",
		Function: "cleanup",
	})

	log.Info("cleanup started")
	defer log.Info("cleanup complete")

	log.Info("gracefully shutting down the HTTP sever")
	if err := a.httpServer.Shutdown(); err != nil {
		log.Err(err).Error("HTTP server graceful shutdown failed")
	}

	log.Info("closing the database connection")
	if err := a.db.Close(); err != nil {
		log.Err(err).Error("error when closing the database connection")
	}
}
