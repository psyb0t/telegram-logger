package v1

import (
	"context"
	"os"
	"sync"

	"github.com/psyb0t/glogger"
)

const (
	serviceNameEnvVarName = "SERVICENAME"
	packageName           = "v1"
)

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
	log := glogger.New(glogger.Caller{
		Service:  os.Getenv(serviceNameEnvVarName),
		Package:  packageName,
		Function: "Run",
	})

	cfg, err := newConfig()
	if err != nil {
		log.Err(err).Error("error when initializing config")

		return err
	}

	glogger.SetLogLevel(glogger.StrToLogLevel(cfg.Logger.Level))
	glogger.SetLogFormat(glogger.StrToLogFormat(cfg.Logger.Format))

	ctx, cancelFunc := context.WithCancel(parentCtx)
	defer cancelFunc()

	log.Debug("initializing app")
	a, err := newApp(ctx, cfg)
	if err != nil {
		log.Err(err).Error("error when initializing app")

		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancelFunc()

		log.Info("starting app")
		if err := a.start(); err != nil {
			log.Err(err).Error("app stopped with error")

			return
		}

		log.Info("app stopped")
	}()

	<-ctx.Done()
	a.stop()
	wg.Wait()

	return nil
}
