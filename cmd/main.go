package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"taskbot/config"
	"taskbot/pkg/bot"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("get config error: [%s]\n", err.Error())
	}

	ctx, finish := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		finish()
	}()

	err = bot.Start(ctx, cfg.Token, cfg.Link, cfg.Port)
	if err != nil {
		log.Fatalf("start bot failed: [%s]\n", err.Error())
	}

	log.Println("bot stopped")

}
