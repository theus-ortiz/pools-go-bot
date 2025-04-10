package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

// CarteirasCommand lida com o comando !carteiras
func CarteirasCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	userPools, err := pools.LoadUserPools(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao carregar suas carteiras.")
		return
	}

	SendCarteirasEmbed(s, m.ChannelID, userPools)
}

// SendCarteirasEmbed envia um embed com as carteiras do usuário
func SendCarteirasEmbed(s *discordgo.Session, channelID string, userPools *pools.UserPools) {
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