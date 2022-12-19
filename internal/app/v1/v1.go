package v1

import (
	"context"
	"os"
	"sync"

	logger "github.com/psyb0t/glogger"
)

const packageName = "v1"

// Run starts the app in a separate goroutine and waits for it to stop.
// It initializes a logger and reads in a configuration file.
// It then creates a context with a cancel function, which is used to stop the app when the context is cancelled.
// The app is started in a separate goroutine, and the function waits for the context to be done before stopping the app and waiting for the goroutine to finish.
// If the app stopped with an error, that error is returned by the function.
//
// parentCtx: the parent context to use for the app's context.
//
// Returns an error if the app stopped with an error or if there was an error when initializing the config.
func Run(parentCtx context.Context) error {
	log := logger.New(logger.Caller{
		Service:  os.Getenv("SERVICENAME"),
		Package:  packageName,
		Function: "Run",
	})

	cfg, err := newConfig()
	if err != nil {
		log.Error("error when initializing config", err)

		return err
	}

	logger.SetLogLevel(logger.StrToLogLevel(cfg.LogLevel))

	ctx, cancelFunc := context.WithCancel(parentCtx)
	defer cancelFunc()

	a, err := newApp(ctx, cfg)
	if err != nil {
		log.Error("error when initializing app", err)

		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancelFunc()

		log.Info("starting app")
		if err := a.start(); err != nil {
			log.Error("app stopped with error", err)

			return
		}

		log.Info("app stopped")
	}()

	<-ctx.Done()
	a.stop()
	wg.Wait()

	return nil
}
