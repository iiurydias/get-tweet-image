package main

import (
	"context"
	"get-tweet-image/app"
	"github.com/micro/go-micro/config"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := runApplication("./config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func runApplication(configFileName string) error {
	configParams, err := getConfigParams(configFileName)
	if err != nil {
		return err
	}
	application, err := app.LoadApp(*configParams)
	if err != nil {
		return errors.Wrap(err, "failed to load application")
	}
	ctx := gracefullyShutdown()
	application.Run()
	defer application.Close()
	select {
	case <-ctx.Done():
		return errors.New("gracefully shutdown")
	}
}

func getConfigParams(configFileName string) (*app.Config, error) {
	var configParams app.Config
	if err := config.LoadFile(configFileName); err != nil {
		return nil, errors.Wrap(err, "failed to load config file")
	}
	if err := config.Scan(&configParams); err != nil {
		return nil, errors.Wrap(err, "failed to scan file")
	}
	return &configParams, nil
}

func gracefullyShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-quit
		cancel()
	}()
	return ctx
}
