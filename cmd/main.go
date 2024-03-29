package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/psyb0t/glogger"
	app "github.com/psyb0t/telegram-logger/internal/app/v1"
)

const (
	serviceNameEnvVarName = "SERVICENAME"
	serviceName           = "telegram-logger"
)

// main is the entry point of the application. It sets the service name
// environment variable, creates a channel to receive os signals, creates
// a context and wait group and starts the app. If the app returns an error
// or an os signal is received, the context is cancelled and the wait group
// is waited to finish.
func main() {
	if err := os.Setenv(serviceNameEnvVarName, serviceName); err != nil {
		panic(err)
	}

	defer os.Unsetenv(serviceNameEnvVarName) //nolint:errcheck

	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  "main",
		Function: "main",
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancelFunc := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info("running app...")
		defer log.Info("app closed")

		errCh <- app.Run(ctx)
		close(errCh)
	}()

	select {
	case <-c:
		log.Info("received interrupt signal")
	case err := <-errCh:
		if err != nil {
			log.Err(err).Error("app encountered an error")
		}
	}

	cancelFunc()

	log.Info("waiting for wait group...")
	wg.Wait()
}
