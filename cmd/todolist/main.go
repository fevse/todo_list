package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/bot"
	"github.com/fevse/todo_list/internal/config"
	"github.com/fevse/todo_list/internal/logger"
	httpserver "github.com/fevse/todo_list/internal/server/http"
	"github.com/fevse/todo_list/internal/storage"
)

func main() {
	flag.Parse()
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewLogger()

	storage := storage.New(conf)
	err = storage.Connect()
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	defer storage.Close()

	err = storage.Migrate()
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	app := app.New(storage, logger)
	httpserver := httpserver.NewServer(app, conf.HTTPServer.Host, conf.HTTPServer.Port)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpserver.Stop(ctx); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := httpserver.Start(ctx)
		if err != nil {
			logger.Logger.Error(err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		bot.Start(app, conf.TgBot.Token)
	}()

	logger.Logger.Info("TODO LIST successfully started")
	<-ctx.Done()
	wg.Wait()
}
