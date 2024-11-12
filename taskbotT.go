package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func startTaskBot(ctx context.Context, httpListenAddr string) error {

	bot, err := tgbotapi.NewBotAPI("_golangcourse_test")

	if err != nil {
		fmt.Println("error happen with init bot :", err)
		return err
	}

	bot.Debug = false

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh := tgbotapi.NewWebhook("http://127.0.0.1:8081")

	_, err = bot.SetWebhook(wh)
	if err != nil {
		fmt.Println("error happen with set webhook :", err)
		return err
	}

	updates := bot.ListenForWebhook("/")

	fmt.Println("start listen :", httpListenAddr)

	go http.ListenAndServe(httpListenAddr, nil)

	router := &Router{
		Data: NewStorage(),
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !strings.Contains(update.Message.Text, "/") {
			txt := "there is no command here"
			log.Printf("user: [%s] request: [%s]\nresp: [%s]\n", update.Message.From.UserName, update.Message.Text, txt)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}
		txt := strings.TrimPrefix(update.Message.Text, "/")
		var command, arg string

		if strings.Contains(txt, "_") {
			index := strings.IndexAny(txt, "_")
			command = txt[:index]
			arg = txt[index+1:]
		} else {
			index := strings.IndexAny(txt, " ")
			if index == -1 {
				command = txt
			} else {
				command = txt[:index]
				arg = txt[index+1:]
			}

		}

		response, err := router.RouteManage(
			command,
			arg,
			update.Message.From.UserName,
			update.Message.Chat.ID,
		)
		if err != nil {
			fmt.Printf("error happen: command [%s], arg [%s], user [%s], chat id [%v]\n error text: %v\n",
				command,
				arg,
				update.Message.From.UserName,
				update.Message.Chat.ID,
				err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "error happen please try use bot later")
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
				resultResponse += message + " "
			}
		}

		fmt.Printf("user: [%s] request: [%s]\nresponse: [%v]\n", update.Message.From.UserName, update.Message.Text, resultResponse)
		fmt.Println()
	}

	return nil
}

func main() {
	err := startTaskBot(context.Background(), ":8081")
	if err != nil {
		log.Fatalln(err)
	}
}

// это заглушка чтобы импорт сохранился
func __dummy() {
	tgbotapi.APIEndpoint = "_dummy"
}
