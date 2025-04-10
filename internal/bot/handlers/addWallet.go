package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

func AddWalletCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) != 3 {
		s.ChannelMessageSend(m.ChannelID, "❌ Uso incorreto. Use: `!addPool <endereço> <rede>`")
		return
	}

	address := args[1]
	network := args[2]
	userID := m.Author.ID

	err := pools.AddPool(userID, address, network)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao adicionar a pool: "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "✅ Pool adicionada com sucesso!")
}