package handlers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "💧 Pools Bot — Central de Comandos",
		Description: "Confira abaixo os comandos disponíveis para interagir com o bot:",
		Color:       0x8B008B, // Roxo escuro
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.State.User.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Pools Bot 💧",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "📌 Comandos Gerais",
				Value: strings.Join([]string{
					"`!ping` — Responde com 'Pong!'",
					"`!oi`, `!ola`, `!eae`, `!hey` — Saudação personalizada",
					"`!help`, `!ajuda` — Mostra esta mensagem de ajuda",
				}, "\n"),
				Inline: false,
			},
			{
				Name: "👛 Comandos de Carteira",
				Value: strings.Join([]string{
					"`!addwallet` — Adiciona uma nova carteira",
					"`!carteiras` — Lista suas carteiras cadastradas",
					"`!excluir` — Exclui uma carteira cadastrada",
				}, "\n"),
				Inline: false,
			},
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
