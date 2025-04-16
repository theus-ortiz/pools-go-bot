package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type rawCache struct {
	Data      string    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	cacheStore = make(map[string]rawCache)
	cacheMutex sync.RWMutex
)

func getCacheKey(address, network string) string {
	return address + "::" + network
}

func getCacheFilePath(address, network string) string {
	fileName := address + "_" + network + ".json"
	return filepath.Join(".cache", fileName)
}

func getCachedRaw(address, network string) (string, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	key := getCacheKey(address, network)

	// Primeiro tenta memória
	if c, exists := cacheStore[key]; exists && time.Now().Before(c.ExpiresAt) {
		return c.Data, true
	}

	// Depois tenta arquivo
	path := getCacheFilePath(address, network)
	fileData, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	var rc rawCache
	if err := json.Unmarshal(fileData, &rc); err != nil || time.Now().After(rc.ExpiresAt) {
		return "", false
	}

	// Atualiza memória a partir do disco
	cacheStore[key] = rc
	return rc.Data, true
}

func setCacheRaw(address, network, data string, ttl time.Duration) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	key := getCacheKey(address, network)
	exp := time.Now().Add(ttl)

	rc := rawCache{
		Data:      data,
		ExpiresAt: exp,
	}

	// Salva em memória
	cacheStore[key] = rc

	// Salva em disco
	_ = os.MkdirAll(".cache", 0755)
	path := getCacheFilePath(address, network)
	_ = os.WriteFile(path, mustJSON(rc), 0644)
}

func mustJSON(v interface{}) []byte {
	data, _ := json.MarshalIndent(v, "", "  ")
	return data
}
