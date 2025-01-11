// configure.Cendit
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cendit.io/garage/config"
	"cendit.io/garage/logger"

	"cendit.io/gate/http/rest/api"
)

const (
	allowConnectionsAfterShutdown = 5
)

func main() {
	cfg := config.Environment()

	sDeps, err := config.ResolveDeps(cfg)
	if err != nil {
		logger.GetLogger().Sugar().Fatalf(`failed to set up dependencies :: %s`, err.Error())
	}

	a := &api.API{
		Variables: cfg,
		Deps:      sDeps,
	}

	go func() {
		// monitor.InitSentry(cfg)
		logger.GetLogger().Sugar().Fatal(a.Serve())
	}()

	// wait for goroutine above
	time.Sleep(1 * time.Second)

	logger.GetLogger().Info(fmt.Sprintf(`
		______  _______  _______  _______  _______  _______  _______  _______  _______

		REST API :::::: http://localhost:%d

		GRAPHQL API :::::: http://localhost:%d/pivot/cendit/graphql

		GraphQL Playground :::::: http://localhost:%d/graphql?witness=playground
	`, cfg.Port, cfg.Port, cfg.Port))

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	logger.GetLogger().Sugar().Infof("Request to shutdown server")
	logger.GetLogger().Sugar().Infof("%d seconds grace to close active connections", allowConnectionsAfterShutdown)

	// close all connections including the http connection
	go func() {
		logger.GetLogger().Info("Closing all open connections...")
		a.Deps.DB.Close()
		logger.GetLogger().Sugar().Fatal(a.Shutdown(allowConnectionsAfterShutdown * time.Second))
	}()

	// start a 5 seconds grace timer
	t := time.NewTimer(5 * time.Second)
	if e := <-t.C; !e.IsZero() {
		os.Exit(1)
	}

}
