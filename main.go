package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/spuxy/barber/handlers"
	api "github.com/spuxy/barber/handlers/v1"
	"go.uber.org/zap"
)

var buildVersion string
var buildCommit string
var buildDate string

var serviceApi string = "barber-api"

func main() {
	logger, err := initLog()
	if err != nil {
		panic(err)
	}

	err = initConfig()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	if err := run(logger); err != nil {
		panic(err)
	}
}

func run(logger *zap.SugaredLogger) error {
	//==============================
	// RUN APPLICATION

	debugMux := handlers.DebugMux(logger, buildVersion, buildCommit, buildDate)
	debug := http.Server{
		Handler: debugMux,
		Addr:    viper.GetString("Server.DebugHost"),
	}

	apiMux := api.ApiMux()
	api := http.Server{
		Handler: apiMux,
		Addr:    viper.GetString("Server.APIHost"),
	}

	go func() {
		logger.Infow("startup", "debug server [started]", debug.Addr)
		if err := debug.ListenAndServe(); err != nil {
			logger.Errorw("shutdown", "debug server [shutdown]", err.Error())
		}
	}()

	// =========================================================================
	// Block goroutine
	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infow("startup", "api server [started]", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer logger.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("Server.ShutdownTimeout"))
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}

func initLog() (*zap.SugaredLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.InitialFields = map[string]interface{}{"service": serviceApi}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func initConfig() error {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in
	return viper.ReadInConfig()     // Find and read the config file
}
