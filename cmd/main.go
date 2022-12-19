package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	logger "github.com/psyb0t/glogger"
	app "github.com/psyb0t/telegram-logger/internal/app/v1"
)

const (
	serviceNameEnvVarName = "SERVICENAME"
	serviceName           = "telegram-logger"
)

func main() {
	os.Setenv(serviceNameEnvVarName, serviceName)
	defer os.Unsetenv(serviceNameEnvVarName)

	log := logger.New(logger.Caller{
		Service:  os.Getenv("SERVICENAME"),
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
			log.Info("app encountered an error:", err)
		}
	}

	cancelFunc()

	log.Info("waiting for wait group...")
	wg.Wait()
}
