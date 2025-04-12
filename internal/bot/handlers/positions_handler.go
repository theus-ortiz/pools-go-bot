package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/graphql"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

const embedColor = 0x00bcd4
const detailedEmbedColor = 0x4caf50

func toFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func isClosed(liquidity string) bool {
	liq := toFloat(liquidity)
	return math.Abs(liq) < 0.000001
}

func buildPositionField(p graphql.PositionRaw, address, network string) *discordgo.MessageEmbedField {
	dep0 := toFloat(p.DepositedToken0)
	dep1 := toFloat(p.DepositedToken1)
	wit0 := toFloat(p.WithdrawnToken0)
	wit1 := toFloat(p.WithdrawnToken1)
	usdIn := toFloat(p.AmountDepositedUSD)
	usdOut := toFloat(p.AmountWithdrawnUSD)

	status := "Aberta ✅"
	if isClosed(p.Liquidity) {
		status = "Fechada 🔒"
	}

	return &discordgo.MessageEmbedField{
		Name: fmt.Sprintf("🔸 Posição ID: `%s`", p.ID),
		Value: fmt.Sprintf(
			"🧠 **Status:** %s\n"+
				"📍 **Endereço:** `%s`\n"+
				"🌐 **Rede:** `%s`\n"+
				"💧 **Token0:** depositado `%.2f` | retirado `%.2f`\n"+
				"💧 **Token1:** depositado `%.4f` | retirado `%.4f`\n"+
				"💵 **USD:** depositado `≈ $%.2f` | retirado `≈ $%.2f`",
			status, address, network, dep0, wit0, dep1, wit1, usdIn, usdOut,
		),
		Inline: false,
	}
}

func buildDetailedField(p graphql.PositionDetailed) *discordgo.MessageEmbedField {
	liquidity := toFloat(p.Liquidity)

	return &discordgo.MessageEmbedField{
		Name: fmt.Sprintf("📘 Posição Detalhada ID: `%s`", p.ID),
		Value: fmt.Sprintf(
			"💧 **Liquidez:** `%.2f`\n"+
				"📊 **Ticks:** [%s → %s] | **Atual:** %s\n"+
				"💹 **Fee Tier:** `%s`\n"+
				"🪙 **Pool:** `%s/%s`\n"+
				"💰 **Depósitos:** token0 `%.4f`, token1 `%.4f`\n"+
				"💸 **Retiradas:** token0 `%.4f`, token1 `%.4f`\n"+
				"🏦 **Taxas Coletadas:** token0 `%.4f`, token1 `%.4f`\n",
			liquidity,
			p.TickLower.TickIdx, p.TickUpper.TickIdx, p.Pool.Tick,
			p.Pool.FeeTier,
			p.Pool.Token0.Symbol, p.Pool.Token1.Symbol,
			toFloat(p.DepositedToken0), toFloat(p.DepositedToken1),
			toFloat(p.WithdrawnToken0), toFloat(p.WithdrawnToken1),
			toFloat(p.CollectedFeesToken0), toFloat(p.CollectedFeesToken1),
		),
		Inline: false,
	}
}
// ListPositionsCommand envia ao usuário do Discord o resumo das posições de liquidez armazenadas.
func ListPositionsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	filePath := fmt.Sprintf("data/pools/%s.json", m.Author.ID)

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Erro ao ler arquivo JSON:", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Não encontrei posições salvas para sua conta.")
		return
	}

	var userPools pools.UserPools
	if err := json.Unmarshal(data, &userPools); err != nil {
		log.Println("Erro ao fazer parse do JSON:", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao interpretar os dados do seu perfil.")
		return
	}

	if userPools.Owner != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "❌ Você não possui posições cadastradas.")
		return
	}

	var openFields []*discordgo.MessageEmbedField
	var closedFields []*discordgo.MessageEmbedField

	for _, position := range userPools.Positions {
		rawData := graphql.QueryPositions(position.Address, position.Network)

		if len(rawData) == 0 || rawData[0] != '{' {
			openFields = append(openFields, &discordgo.MessageEmbedField{
				Name:   "⚠️ Erro na consulta",
				Value:  fmt.Sprintf("Falha ao consultar `%s` na rede `%s`", position.Address, position.Network),
				Inline: false,
			})
			continue
		}

		var subgraphResp graphql.SubgraphResponse
		if err := json.Unmarshal([]byte(rawData), &subgraphResp); err != nil || len(subgraphResp.Data.Positions) == 0 {
			openFields = append(openFields, &discordgo.MessageEmbedField{
				Name:   "⚠️ Posição não encontrada",
				Value:  fmt.Sprintf("Erro ao consultar endereço `%s` na rede `%s`", position.Address, position.Network),
				Inline: false,
			})
			continue
		}

		for _, p := range subgraphResp.Data.Positions {
			field := buildPositionField(p, position.Address, position.Network)
			if isClosed(p.Liquidity) {
				closedFields = append(closedFields, field)
			} else {
				openFields = append(openFields, field)
			}
		}
	}

	if len(openFields)+len(closedFields) == 0 {
		s.ChannelMessageSend(m.ChannelID, "⚠️ Nenhuma posição foi encontrada nas carteiras salvas.")
		return
	}

	var embeds []*discordgo.MessageEmbed

	if len(openFields) > 0 {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       "🔓 Posições Abertas",
			Description: fmt.Sprintf("👤 Usuário: <@%s>", m.Author.ID),
			Color:       embedColor,
			Fields:      openFields,
		})
	}

	if len(closedFields) > 0 {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       "🔒 Posições Fechadas",
			Description: fmt.Sprintf("👤 Usuário: <@%s>", m.Author.ID),
			Color:       0x9e9e9e,
			Fields:      closedFields,
		})
	}

	for _, embed := range embeds {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Dados obtidos via Subgraph da Uniswap v3",
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

func ListDetailedPositionsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	filePath := fmt.Sprintf("data/pools/%s.json", m.Author.ID)

	data, err := os.ReadFile(filePath)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Não encontrei posições salvas para sua conta.")
		return
	}

	var userPools pools.UserPools
	if err := json.Unmarshal(data, &userPools); err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao interpretar os dados do seu perfil.")
		return
	}

	if userPools.Owner != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "❌ Você não possui posições cadastradas.")
		return
	}

	for _, position := range userPools.Positions {
		raw := graphql.QueryPositions(position.Address, position.Network)

		if len(raw) == 0 || raw[0] != '{' {
			continue
		}

		var res graphql.SubgraphResponse
		if err := json.Unmarshal([]byte(raw), &res); err != nil {
			continue
		}

		for _, pos := range res.Data.Positions {
			if isClosed(pos.Liquidity) {
				continue
			}

			detailRaw := graphql.QueryPositionByID(pos.ID, position.Network)
			if len(detailRaw) == 0 || detailRaw[0] != '{' {
				continue
			}

			var detail graphql.SubgraphPositionByIDResponse
			if err := json.Unmarshal([]byte(detailRaw), &detail); err != nil {
				continue
			}

			if detail.Data.Position.ID == "" {
				continue
			}

			field := buildDetailedField(detail.Data.Position)

			// Criar um embed individual por pool
			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("📘 Posição: %s/%s | Fee: %s",
					detail.Data.Position.Pool.Token0.Symbol,
					detail.Data.Position.Pool.Token1.Symbol,
					detail.Data.Position.Pool.FeeTier,
				),
				Description: fmt.Sprintf("👤 Usuário: <@%s>", m.Author.ID),
				Color:       detailedEmbedColor,
				Fields:      []*discordgo.MessageEmbedField{field},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Detalhes obtidos via Subgraph da Uniswap v3",
				},
			}

			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}