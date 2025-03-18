package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/bot"
	"github.com/fevse/todo_list/internal/config"
	httpserver "github.com/fevse/todo_list/internal/server/http"
	"github.com/fevse/todo_list/internal/storage"
)

func main() {
	flag.Parse()
	fmt.Println("***TODO LIST***")
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.New(conf)
	err = storage.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	err = storage.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	app := app.New(storage)
	httpserver := httpserver.NewServer(app, conf.HTTPServer.Host, conf.HTTPServer.Port)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpserver.Stop(ctx); err != nil {
			log.Print(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := httpserver.Start(ctx)
		if err != nil {
			fmt.Println(err)
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		bot.Start(app, conf.TgBot.Token)
	}()

	fmt.Println("This is fine")
	<-ctx.Done()
	wg.Wait()
}
