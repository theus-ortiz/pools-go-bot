package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/theus-ortiz/pools-bot/internal/bot"
	"github.com/theus-ortiz/pools-bot/internal/config"
	"github.com/theus-ortiz/pools-bot/internal/graphql"
)

func main() {
	cfg := config.Load()

	graphql.ConferSubgraph("0xb8566807fBBa74c6c35FD6C2Af2094C4Ba4104ED", "polygon")

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
