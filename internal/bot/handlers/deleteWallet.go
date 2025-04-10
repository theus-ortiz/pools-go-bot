package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

func ExcluirCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		SendCarteirasEmbed(s, m.ChannelID, userPools)
		return
	}

	userPools.Positions = filtered
	if err := pools.SaveUserPools(userPools); err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao salvar alterações.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "🗑️ Carteira excluída com sucesso!")
}