package handlers

import "github.com/bwmarrin/discordgo"

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "🏓 Pong!")
}