package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"math"
// 	"strconv"

// 	"github.com/bwmarrin/discordgo"
// 	"github.com/theus-ortiz/pools-bot/internal/graphql"
// 	"github.com/theus-ortiz/pools-bot/internal/pools"
// )

// func PriceFromTick(tick int, token0Decimals int, token1Decimals int) float64 {
// 	base := 1.0001
// 	decimalsFactor := math.Pow10(token1Decimals - token0Decimals)
// 	price := math.Pow(base, float64(tick)) * decimalsFactor
// 	return price
// }

// func HandleUserPools(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	filePath := fmt.Sprintf("data/pools/%s.json", m.Author.ID)
// 	data, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		log.Println("Erro ao ler arquivo JSON:", err)
// 		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao acessar os dados.")
// 		return
// 	}

// 	var userPools pools.UserPools
// 	err = json.Unmarshal(data, &userPools)
// 	if err != nil {
// 		log.Println("Erro ao interpretar JSON:", err)
// 		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao interpretar os dados.")
// 		return
// 	}

// 	if userPools.Owner != m.Author.ID {
// 		s.ChannelMessageSend(m.ChannelID, "❌ Você não possui posições cadastradas.")
// 		return
// 	}

// 	var embeds []*discordgo.MessageEmbed

// 	for _, position := range userPools.Positions {
// 		raw := graphql.PositionByID(position.Address, position.Network)
// 		log.Println("Raw GraphQL response:", raw)

// 		var result struct {
// 			Data struct {
// 				Position struct {
// 					ID                  string `json:"id"`
// 					Liquidity           string `json:"liquidity"`
// 					DepositedToken0     string `json:"depositedToken0"`
// 					DepositedToken1     string `json:"depositedToken1"`
// 					WithdrawnToken0     string `json:"withdrawnToken0"`
// 					WithdrawnToken1     string `json:"withdrawnToken1"`
// 					CollectedFeesToken0 string `json:"collectedFeesToken0"`
// 					CollectedFeesToken1 string `json:"collectedFeesToken1"`
// 					TickLower           struct {
// 						TickIdx string `json:"tickIdx"`
// 					} `json:"tickLower"`
// 					TickUpper struct {
// 						TickIdx string `json:"tickIdx"`
// 					} `json:"tickUpper"`
// 					Pool struct {
// 						ID        string `json:"id"`
// 						FeeTier   string `json:"feeTier"`
// 						Tick      string `json:"tick"`
// 						SqrtPrice string `json:"sqrtPrice"`
// 						Token0    struct {
// 							Symbol   string `json:"symbol"`
// 							Decimals string `json:"decimals"`
// 						} `json:"token0"`
// 						Token1 struct {
// 							Symbol   string `json:"symbol"`
// 							Decimals string `json:"decimals"`
// 						} `json:"token1"`
// 					} `json:"pool"`
// 				} `json:"position"`
// 			} `json:"data"`
// 		}

// 		err := json.Unmarshal([]byte(raw), &result)
// 		if err != nil || result.Data.Position.ID == "" {
// 			log.Println("Erro ao interpretar retorno da posição:", err)
// 			continue
// 		}

// 		p := result.Data.Position
// 		liq, _ := strconv.ParseFloat(p.Liquidity, 64)
// 		if liq == 0 {
// 			continue // pula posições fechadas
// 		}

// 		// Parse
// 		dep0, _ := strconv.ParseFloat(p.DepositedToken0, 64)
// 		dep1, _ := strconv.ParseFloat(p.DepositedToken1, 64)
// 		wit0, _ := strconv.ParseFloat(p.WithdrawnToken0, 64)
// 		wit1, _ := strconv.ParseFloat(p.WithdrawnToken1, 64)
// 		fee0, _ := strconv.ParseFloat(p.CollectedFeesToken0, 64)
// 		fee1, _ := strconv.ParseFloat(p.CollectedFeesToken1, 64)

// 		tickLower, _ := strconv.Atoi(p.TickLower.TickIdx)
// 		tickUpper, _ := strconv.Atoi(p.TickUpper.TickIdx)
// 		token0Decimals, _ := strconv.Atoi(p.Pool.Token0.Decimals)
// 		token1Decimals, _ := strconv.Atoi(p.Pool.Token1.Decimals)

// 		priceLower := PriceFromTick(tickLower, token0Decimals, token1Decimals)
// 		priceUpper := PriceFromTick(tickUpper, token0Decimals, token1Decimals)
// 		currentPrice := PriceFromTick(parseInt(p.Pool.Tick), token0Decimals, token1Decimals)

// 		// Investimento
// 		invested0 := dep0 - wit0
// 		invested1 := dep1 - wit1

// 		// Tarifas acumuladas (em USD, usando preço atual)
// 		feesUSD := fee0 + (fee1 * currentPrice)

// 		embed := &discordgo.MessageEmbed{
// 			Title: fmt.Sprintf("🔹 Pool %s/%s", p.Pool.Token0.Symbol, p.Pool.Token1.Symbol),
// 			Fields: []*discordgo.MessageEmbedField{
// 				{
// 					Name:  "📍 Posição",
// 					Value: fmt.Sprintf("ID: `%s`\nRede: **%s**", p.ID, position.Network),
// 				},
// 				{
// 					Name: "🧮 Liquidez e Investimento",
// 					Value: fmt.Sprintf(
// 						"`%.2f` %s | `%.4f` %s",
// 						invested0, p.Pool.Token0.Symbol,
// 						invested1, p.Pool.Token1.Symbol,
// 					),
// 				},
// 				{
// 					Name:  "💸 Tarifas Acumuladas",
// 					Value: fmt.Sprintf("≈ `$%.2f`", feesUSD),
// 				},
// 				{
// 					Name:  "📊 Faixa de Preço",
// 					Value: fmt.Sprintf("`%.4f` → `%.4f` (%s/%s)", priceLower, priceUpper, p.Pool.Token0.Symbol, p.Pool.Token1.Symbol),
// 				},
// 			},
// 			Color: 0x00ffcc,
// 		}

// 		embeds = append(embeds, embed)
// 	}

// 	if len(embeds) == 0 {
// 		s.ChannelMessageSend(m.ChannelID, "📭 Nenhuma posição **aberta** foi encontrada.")
// 		return
// 	}

// 	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
// 		Content: "**📊 Resumo das Pools (via Subgraph Uniswap v3)**",
// 		Embeds:  embeds,
// 	})
// }


