package bot

import (
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func NewDiscordBot(token string) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	dg.AddHandler(MessageRouter)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		return nil, err
	}

	return &Bot{Session: dg}, nil
}

func (b *Bot) Close() {
	b.Session.Close()
}
