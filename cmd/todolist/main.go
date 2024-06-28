package main

import (
	"flag"
	"fmt"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/storage"
)

var command string

func init() {
	flag.StringVar(&command, "command", "", "write command to execute")
}

func main() {
	flag.Parse()
	fmt.Println("***TODO LIST***")

	storage := storage.New()
	app := app.New(storage)

	switch command {
	case "list":
		tasks, err := app.Storage.ShowList()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, t := range tasks {
			fmt.Printf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
		}
	case "create":
		var t, s string
		fmt.Println("Add new task")
		fmt.Print("Title: ")
		fmt.Scan(&t)
		fmt.Print("Status: ")
		fmt.Scan(&s)
		err := app.Storage.CreateTask(t, s)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("need command to execute")
	}

	fmt.Println("THE END")
}
