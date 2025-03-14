package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/fevse/todo_list/internal/app"
	"github.com/fevse/todo_list/internal/config"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(s *app.App, t config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(ShowTask(s)),
	}

	b, err := bot.New(t.TgBot.Token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func ShowTask(s *app.App) func(context.Context, *bot.Bot, *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		id, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			fmt.Println(err)
		}

		t, err := s.Storage.ShowTask(id)
		if err != nil {
			fmt.Println(err)
		}
		task := fmt.Sprintf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   task,
		})
	}
}

// func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	id, err := strconv.Atoi(update.Message.Text)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	t, err := storage.Storage.ShowTask(id)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	task := fmt.Sprintf("#%d %s - %s: %v\n", t.ID, t.Title, t.Status, t.Created.Format("02.01.2006 15:04:05"))
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   task,
// 	})
// }
