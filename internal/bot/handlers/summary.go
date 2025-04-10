package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/debank"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

func ResumoCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	userPools, err := pools.LoadUserPools(userID)
	if err != nil || len(userPools.Positions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "📭 Você não tem carteiras cadastradas. Use `!addPool`.")
		return
	}

	var summary strings.Builder
	summary.WriteString("📊 **Resumo das Pools (via DeBank)**\n\n")

	for _, pos := range userPools.Positions {
		summary.WriteString(fmt.Sprintf("🔸 **%s** (`%s`)\n", strings.ToUpper(pos.Network), pos.Address))

		protocols, err := debank.FetchProtocols(pos.Address)
		if err != nil || len(protocols) == 0 {
			summary.WriteString("  ⚠️ Nenhuma pool encontrada ou erro na consulta.\n\n")
			continue
		}

		for _, p := range protocols {
			if p.Portfolio.NetUSD > 0 {
				summary.WriteString(fmt.Sprintf("  • %s: **$%.2f**\n", p.Name, p.Portfolio.NetUSD))
			}
		}

		summary.WriteString("\n")
	}

	s.ChannelMessageSend(m.ChannelID, summary.String())
}
