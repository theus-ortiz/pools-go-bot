package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/theus-ortiz/pools-bot/internal/config"
)

// QueryPositions consulta um subgraph com uma carteira e rede específicas.
// Retorna a resposta JSON como string ou uma mensagem de erro formatada.
func QueryPositions(walletAddress string, network string) string {
	cfg := config.Load()

	// Verifica se a rede é suportada e pega a URL do subgraph correspondente
	url, ok := subgraphURLs[network]
	if !ok {
		return fmt.Sprintf("❌ Rede '%s' não suportada.", network)
	}

	// Lê o conteúdo da query GraphQL do arquivo
	query, err := ReadQuery("positions.graphql")
	if err != nil {
		return fmt.Sprintf("❌ Erro ao carregar query: %v", err)
	}

	// Substitui o caractere '?' pelo endereço da carteira (ajuste para seu template de query)
	query = ReplaceQueryVariable(query, "{{WALLET_ADDRESS}}", walletAddress)

	// Cria o payload da requisição GraphQL
	payload := map[string]interface{}{
		"query":         query,
		"operationName": "Subgraphs",
		"variables":     map[string]interface{}{},
	}

	// Converte o payload para JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("❌ Erro ao serializar payload: %v", err)
	}

	// Cria uma nova requisição HTTP POST com o JSON
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Sprintf("❌ Erro ao criar requisição: %v", err)
	}

	// Define os headers da requisição
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.GraphAPIKey)

	// Envia a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("❌ Erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	// Lê a resposta da requisição
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("❌ Erro ao ler resposta: %v", err)
	}

	// Retorna o corpo da resposta como string
	return string(body)
}
