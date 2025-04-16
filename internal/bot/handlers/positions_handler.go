package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

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

func loadUserPools(userID string) (pools.UserPools, error) {
	filePath := fmt.Sprintf("data/pools/%s.json", userID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return pools.UserPools{}, err
	}
	var userPools pools.UserPools
	if err := json.Unmarshal(data, &userPools); err != nil {
		return pools.UserPools{}, err
	}
	if userPools.Owner != userID {
		return pools.UserPools{}, fmt.Errorf("usuário não autorizado")
	}
	return userPools, nil
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
			"🧐 **Status:** %s\n"+
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

func buildDetailedField(p graphql.PositionDetailed, userID string) *discordgo.MessageEmbed {
	poolData := pools.BuildPoolDataFromPosition(p)
	summary := pools.FormatPoolSummary(poolData)

	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("📘 Posição: %s/%s | Fee: %s",
			p.Pool.Token0.Symbol,
			p.Pool.Token1.Symbol,
			p.Pool.FeeTier,
		),
		Description: fmt.Sprintf("👤 Usuário: <@%s>", userID),
		Color:       detailedEmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("📊 Resumo da Posição ID: `%s`", p.ID),
				Value:  summary,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Detalhes obtidos via Subgraph da Uniswap v3",
		},
	}
}

func buildOutOfRangeField(p graphql.PositionDetailed, userID string) *discordgo.MessageEmbed {
	poolData := pools.BuildPoolDataFromPosition(p)
	summary := pools.FormatPoolSummary(poolData)

	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("🛑 FORA DO INTERVALO: %s/%s | Fee: %s",
			p.Pool.Token0.Symbol,
			p.Pool.Token1.Symbol,
			p.Pool.FeeTier,
		),
		Description: fmt.Sprintf("👤 Usuário: <@%s>\n⚠️ Esta posição **não está ativa** no momento. O preço atual está fora da faixa definida.", userID),
		Color:       0xf44336,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("📊 Resumo da Posição ID: `%s`", p.ID),
				Value:  summary,
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Detalhes obtidos via Subgraph da Uniswap v3",
		},
	}
}

func ListPositionsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	os.MkdirAll(".cache", os.ModePerm)

	userPools, err := loadUserPools(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ "+err.Error())
		return
	}

	var openFields []*discordgo.MessageEmbedField
	var closedFields []*discordgo.MessageEmbedField

	for _, position := range userPools.Positions {
		rawData, ok := getCachedRaw(position.Address, position.Network)
		if !ok {
			rawData = graphql.QueryPositions(position.Address, position.Network)
			if len(rawData) > 0 && rawData[0] == '{' {
				setCacheRaw(position.Address, position.Network, rawData, 5*time.Minute)
			}
		}

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
	userPools, err := loadUserPools(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ "+err.Error())
		return
	}

	for _, position := range userPools.Positions {
		raw, ok := getCachedRaw(position.Address, position.Network)
		if !ok {
			raw = graphql.QueryPositions(position.Address, position.Network)
			if len(raw) > 0 && raw[0] == '{' {
				setCacheRaw(position.Address, position.Network, raw, 5*time.Minute)
			}
		}

		if len(raw) == 0 || raw[0] != '{' {
			log.Printf("⚠️ Resposta inválida do subgraph para %s na %s", position.Address, position.Network)
			continue
		}

		var res graphql.SubgraphResponse
		if err := json.Unmarshal([]byte(raw), &res); err != nil {
			log.Printf("⚠️ Erro ao fazer unmarshal para %s: %v", position.Address, err)
			continue
		}

		for _, pos := range res.Data.Positions {
			if isClosed(pos.Liquidity) {
				log.Printf("🔒 Posição %s está fechada (liquidez: %s)", pos.ID, pos.Liquidity)
				continue
			}

			detailRaw, ok := getCachedRaw(pos.ID, position.Network)
			if !ok {
				detailRaw = graphql.QueryPositionByID(pos.ID, position.Network)
				if len(detailRaw) > 0 && detailRaw[0] == '{' {
					setCacheRaw(pos.ID, position.Network, detailRaw, 5*time.Minute)
				}
			}

			if len(detailRaw) == 0 || detailRaw[0] != '{' {
				log.Printf("⚠️ Falha ao buscar detalhes da posição %s", pos.ID)
				continue
			}

			var detail graphql.SubgraphPositionByIDResponse
			if err := json.Unmarshal([]byte(detailRaw), &detail); err != nil {
				log.Printf("⚠️ Erro ao parsear detalhes da posição %s: %v", pos.ID, err)
				continue
			}

			if detail.Data.Position.ID == "" {
				log.Printf("⚠️ Dados da posição %s vieram vazios", pos.ID)
				continue
			}

			posDetail := detail.Data.Position

			// Acesse os valores de TickIdx diretamente.
			tick, err := strconv.Atoi(posDetail.Pool.Tick)
			if err != nil {
				log.Printf("❌ Erro ao converter Tick (ID: %s): %v", posDetail.ID, err)
				continue
			}

			tickLower, err1 := strconv.Atoi(posDetail.TickLower.TickIdx)
			tickUpper, err2 := strconv.Atoi(posDetail.TickUpper.TickIdx)

			if err1 != nil || err2 != nil {
				log.Printf("❌ Erro ao converter os ticks (ID: %s): %v %v", posDetail.ID, err1, err2)
				continue
			}

			var embed *discordgo.MessageEmbed

			if tick < tickLower || tick > tickUpper {
				log.Printf("🔺 FORA DO INTERVALO: ID %s (Tick atual: %d | Faixa: %d ~ %d)", posDetail.ID, tick, tickLower, tickUpper)
				embed = buildOutOfRangeField(posDetail, m.Author.ID)
			} else {
				log.Printf("✅ DENTRO DO INTERVALO: ID %s (Tick atual: %d | Faixa: %d ~ %d)", posDetail.ID, tick, tickLower, tickUpper)
				embed = buildDetailedField(posDetail, m.Author.ID)
			}

			if embed != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
			}
		}
	}
}
