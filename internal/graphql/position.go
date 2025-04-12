package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/theus-ortiz/pools-bot/internal/config"
)

func QueryPositionByID(positionID string, network string) string {
	cfg := config.Load()

	url, ok := subgraphURLs[network]
	if !ok {
		return fmt.Sprintf("❌ Rede '%s' não suportada.", network)
	}

	query, err := ReadQuery("position.graphql")
	if err != nil {
		return fmt.Sprintf("❌ Erro ao carregar query: %v", err)
	}

	// Substitui o placeholder na query
	query = ReplaceQueryVariable(query, "{{POSITION_ID}}", positionID)

	payload := map[string]interface{}{
		"query":     query,
		"variables": map[string]interface{}{},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("❌ Erro ao serializar payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Sprintf("❌ Erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.GraphAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("❌ Erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("❌ Erro ao ler resposta: %v", err)
	}

	return string(body)
}
