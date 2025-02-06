package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fevse/todo_list/internal/app"
)

func Cli(app app.App) {
	var command string
	for {
		fmt.Scan(&command)
		switch command {
		case "list":
			list(app)
		case "task":
			task(app)
		case "create":
			create(app)
		case "delete":
			del(app)
		case "exit", "q", "quit", "close":
			fmt.Println("THE END")
			os.Exit(1)
		default:
			fmt.Println("wrong command, use e.g list, task, create, delete, exit")
		}
	}
}

func list(app app.App) {
	tasks, err := app.Storage.ShowList()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range tasks {
		fmt.Printf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
	}
}

func task(app app.App) {
	var id int
	fmt.Println("Show task")
	fmt.Print("ID: ")
	fmt.Scan(&id)
	t, err := app.Storage.ShowTask(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
}

func create(app app.App) {
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
	fmt.Println("New task has been successfully added")
}

func del(app app.App) {
	var id int
	fmt.Println("Delete task")
	fmt.Print("ID: ")
	fmt.Scan(&id)
	err := app.Storage.DeleteTask(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Task %d has been successfully deleted\n", id)
}
