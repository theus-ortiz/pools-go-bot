package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/theus-ortiz/pools-bot/internal/config"
)

var subgraphURLs = map[string]string{
	"polygon": "https://gateway.thegraph.com/api/subgraphs/id/3hCPRGf4z88VC5rsBKU5AA9FBBq5nF3jbKJG7VZCbhjm",
	"base":    "https://gateway.thegraph.com/api/subgraphs/id/43Hwfi3dJSoGpyas9VwNoDAv55yjgGrPpNSmbQZArzMG",
}

func ConferSubgraph(walletAddress string, network string) string {
	cfg := config.Load()

	url, ok := subgraphURLs[network]
	if !ok {
		return fmt.Sprintf("❌ Rede '%s' não suportada.", network)
	}

	query, err := ReadQuery("id_pool.graphql")
	if err != nil {
		return fmt.Sprintf("❌ Erro ao carregar query: %v", err)
	}

	query = strings.Replace(query, "?", walletAddress, 1)

	payload := map[string]interface{}{
		"query":         query,
		"operationName": "Subgraphs",
		"variables":     map[string]interface{}{},
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