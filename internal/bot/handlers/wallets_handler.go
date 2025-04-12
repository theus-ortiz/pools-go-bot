package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

// SendWalletsEmbed envia um embed com as carteiras do usuário
func SendWalletsEmbed(s *discordgo.Session, channelID string, userPools *pools.UserPools) {
	if len(userPools.Positions) == 0 {
		s.ChannelMessageSend(channelID, "📭 Você ainda não cadastrou nenhuma carteira.")
		return
	}

	fields := []*discordgo.MessageEmbedField{}
	for _, pos := range userPools.Positions {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "🔹 Endereço",
			Value:  fmt.Sprintf("`%s`\n🌐 Rede: **%s**", pos.Address, pos.Network),
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:  "💼 Suas carteiras cadastradas",
		Color:  0x00c7a9,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Use !addPool para adicionar mais carteiras.",
		},
	}

	s.ChannelMessageSendEmbed(channelID, embed)
}

func ListWalletsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	userPools, err := pools.LoadUserPools(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao carregar suas carteiras.")
		return
	}

	SendWalletsEmbed(s, m.ChannelID, userPools)
}

func AddWalletCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)

	if len(args) != 3 {
		s.ChannelMessageSend(m.ChannelID, "❌ Uso incorreto. Use: `!addPool <endereço> <rede>`")
		return
	}

	address := args[1]
	network := args[2]
	userID := m.Author.ID

	err := pools.AddPosition(userID, address, network)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao adicionar a pool: "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "✅ Pool adicionada com sucesso!")
}

func DeleteWalletCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "❌ Uso incorreto. Utilize: `!excluir <endereço> <rede>`")
		return
	}

	address := strings.ToLower(strings.TrimSpace(args[1]))
	network := strings.ToLower(strings.TrimSpace(args[2]))
	userID := m.Author.ID

	userPools, err := pools.LoadUserPools(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao carregar suas carteiras.")
		return
	}

	originalLen := len(userPools.Positions)
	filtered := []pools.Position{}

	for _, pos := range userPools.Positions {
		if strings.ToLower(pos.Address) == address && strings.ToLower(pos.Network) == network {
			continue // exclui essa
		}
		filtered = append(filtered, pos)
	}

	if len(filtered) == originalLen {
		// Nenhuma carteira removida – mostra embed
		s.ChannelMessageSend(m.ChannelID, "⚠️ Endereço ou rede não encontrados. Veja suas carteiras abaixo:")
		SendWalletsEmbed(s, m.ChannelID, userPools)
		return
	}

	userPools.Positions = filtered
	if err := pools.SaveUserPools(userPools); err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao salvar alterações.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "🗑️ Carteira excluída com sucesso!")
}


