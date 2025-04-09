package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Greetings(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "👋 Olá!",
		Description: fmt.Sprintf("Fala aí, **%s**! Tudo certo?", m.Author.Username),
		Color:       0x8B008B,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: m.Author.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Pools Bot 💧",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
