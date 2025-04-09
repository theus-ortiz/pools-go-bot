package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorar mensagens do próprio bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Processar apenas comandos que começam com "!"
	if strings.HasPrefix(m.Content, "!") {
		// Remover o "!" e converter para minúsculas
		command := strings.ToLower(strings.TrimSpace(m.Content[1:]))

		// Determinar ação baseada no comando
		switch command {
		case "ping":
			s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")

		case "oi", "ola", "eae", "hey":
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

		case "help", "ajuda", "comandos":
			helpMessage := "**📋 Lista de Comandos:**\n\n" +
				"**!ping** - Testa se o bot está online\n" +
				"**!oi** - Receba um cumprimento\n" +
				"**!help** - Mostra esta mensagem\n" +
				"**!info** - Informações sobre o bot\n" +
				"**!convite** - Link para convidar o bot"
			s.ChannelMessageSend(m.ChannelID, helpMessage)

		case "info", "sobre":
			infoMsg := "🤖 **Bot Básico**\n\n" +
				"Versão: 1.0\n" +
				"Criado por: Você\n" +
				"Descrição: Um simples bot Discord com comandos básicos"
			s.ChannelMessageSend(m.ChannelID, infoMsg)

		case "convite", "invite":
			s.ChannelMessageSend(m.ChannelID, "🔗 Link para convidar o bot: [coloque seu link aqui]")

		case "hora", "tempo":
			s.ChannelMessageSend(m.ChannelID, "⏰ Função de hora ainda não implementada!")

		default:
			s.ChannelMessageSend(m.ChannelID, "Comando não reconhecido. Digite **!help** para ver os comandos disponíveis.")
		}
	}
}
