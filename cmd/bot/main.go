package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/theus-ortiz/pools-bot/internal/bot"
	"github.com/theus-ortiz/pools-bot/internal/config"
)

func main() {
	cfg := config.Load()

	discordBot, err := bot.NewDiscordBot(cfg.DiscordToken)
	if err != nil {
		log.Fatal("Erro ao iniciar bot do Discord:", err)
	}

	fmt.Println("Bot está rodando!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discordBot.Close()
}
