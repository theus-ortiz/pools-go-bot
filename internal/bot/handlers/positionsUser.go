package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/theus-ortiz/pools-bot/internal/graphql"
	"github.com/theus-ortiz/pools-bot/internal/pools"
)

func HandleUserPools(s *discordgo.Session, m *discordgo.MessageCreate) {
	filePath := fmt.Sprintf("data/pools/%s.json", m.Author.ID)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Erro ao ler arquivo JSON:", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao acessar os dados.")
		return
	}

	var userPools pools.UserPools
	err = json.Unmarshal(data, &userPools)
	if err != nil {
		log.Println("Erro ao fazer parse do JSON:", err)
		s.ChannelMessageSend(m.ChannelID, "❌ Erro ao interpretar os dados.")
		return
	}

	// Verifica se o autor do comando tem dados salvos
	if userPools.Owner != m.Author.ID {
		s.ChannelMessageSend(m.ChannelID, "❌ Você não possui posições cadastradas.")
		return
	}

	var resposta strings.Builder
	resposta.WriteString("📊 **Resumo das Pools**\n\n")

	for _, position := range userPools.Positions {
		resposta.WriteString(fmt.Sprintf("🔸 **%s** na rede **%s**\n", position.Address, position.Network))
		info := graphql.ConferSubgraph(position.Address, position.Network)
		resposta.WriteString(info + "\n")
	}

	s.ChannelMessageSend(m.ChannelID, resposta.String())
}
