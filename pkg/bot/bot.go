package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"taskbot/pkg/router"
	"taskbot/pkg/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RouteManager interface {
	ManageCommand(command, arg, user string, chatID int64) (interface{}, error)
}

func NewRouteManager() (RouteManager, error) {
	storage, err := storage.NewStorage()
	if err != nil {
		return nil, err
	}
	return router.NewRouter(storage), nil
}

type BotCfg struct {
	Token string
	Link  string
	Port  string
}

func Start(ctx context.Context, token, link, port string) error {

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Printf("init bot error: [%s]\n", err.Error())
		return err
	}

	bot.Debug = false

	log.Printf("Authorized on account: [%s]\n", bot.Self.UserName)

	wh := tgbotapi.NewWebhook(link)

	_, err = bot.SetWebhook(wh)
	if err != nil {
		log.Printf("set webhook error: [%s]\n", err.Error())
		return err
	}

	updates := bot.ListenForWebhook("/")

	rm, err := NewRouteManager()
	if err != nil {
		log.Printf("init router error: [%s]\n", err.Error())
		return err
	}

	server := http.Server{
		Addr: ":" + port,
	}
	log.Printf("start listen: [%s]\n", port)
	go server.ListenAndServe()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			fmt.Println("=============")
			if !update.Message.IsCommand() {
				txt := "выберите команду из списка команд"
				log.Printf("user: [%s] request: [%s]\nresp: [%s]\n", update.Message.From.UserName, update.Message.Text, txt)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

			} else {

				txt := strings.TrimPrefix(update.Message.Text, "/")
				var command, arg string

				if update.Message.CommandArguments() != "" {
					command = update.Message.Command()
					arg = update.Message.CommandArguments()
				} else if strings.Contains(txt, "_") {

					data := strings.Split(txt, "_")
					command = data[0]
					arg = strings.Join(data[1:], "_")
				} else {
					command = update.Message.Command()
				}

				response, err := rm.ManageCommand(
					command,
					arg,
					update.Message.From.UserName,
					update.Message.Chat.ID,
				)
				if err != nil {
					log.Printf("error happen: command [%s], arg [%s], user [%s], chat id [%d]\n error text: %s\n",
						command,
						arg,
						update.Message.From.UserName,
						update.Message.Chat.ID,
						err.Error())
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "упс, что-то случилось, попробуйте воспользоваться ботом позже")
					bot.Send(msg)
					continue
				}

				var resultResponse string
				if message, ok := response.(string); ok {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
					bot.Send(msg)
					resultResponse = message

				} else if data, ok := response.(map[int64]string); ok {
					for chatID, message := range data {
						msg := tgbotapi.NewMessage(chatID, message)
						bot.Send(msg)
						resultResponse += fmt.Sprintf("{chat id: %v; message: %s}\n", chatID, message)
					}
				}

				log.Printf("user: [%s]; request: [%s]; response: [%v]\n", update.Message.From.UserName, update.Message.Text, resultResponse)
			}
		}
	}

	server.Shutdown(context.Background())

	return nil
}
