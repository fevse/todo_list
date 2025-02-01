package main

import (
	"flag"
	"fmt"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/cli"
	"github.com/fevse/todo_list/internal/storage"
)

func main() {
	flag.Parse()
	fmt.Println("***TODO LIST***")

	storage := storage.New()
	app := app.New(storage)
	app.Storage.Migrate("db/migrations")
	cli.Cli(*app)
}
