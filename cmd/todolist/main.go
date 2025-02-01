package main

import (
	"flag"
	"fmt"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/cli"
	"github.com/fevse/todo_list/internal/config"
	"github.com/fevse/todo_list/internal/storage"
)

func main() {
	flag.Parse()
	fmt.Println("***TODO LIST***")
	conf, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	storage := storage.New(conf)
	app := app.New(storage)
	app.Storage.Migrate()
	cli.Cli(*app)
}
