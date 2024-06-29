package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

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
	case "task":
		var id int64
		fmt.Println("Show task")
		fmt.Print("ID: ")
		fmt.Scan(&id)
		t, err := app.Storage.ShowTask(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
	case "create":
		var t, s string
		fmt.Println("Add new task")
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Title: ")
		t, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		t, _ = strings.CutSuffix(t, "\n")
		fmt.Print("Status: ")
		s, err = buf.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		s, _ = strings.CutSuffix(s, "\n")
		err = app.Storage.CreateTask(t, s)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		var id int64
		fmt.Println("Delete task")
		fmt.Print("ID: ")
		fmt.Scan(&id)
		err := app.Storage.DeleteTask(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("need command to execute")
	}

	fmt.Println("THE END")
}
