package handlers

import "github.com/bwmarrin/discordgo"

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := "**📋 Lista de Comandos:**\n\n" +
		"**!ping** - Testa se o bot está online\n" +
		"**!oi** - Receba um cumprimento\n" +
		"**!help** - Mostra esta mensagem\n" +
		"**!info** - Informações sobre o bot\n" +
		"**!convite** - Link para convidar o bot"
	s.ChannelMessageSend(m.ChannelID, helpMessage)
}
