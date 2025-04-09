package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ .env não encontrado, usando variáveis do sistema.")
		}

		cfg = &Config{
			DiscordToken: mustGetEnv("DISCORD_TOKEN"),
		}
	})

	return cfg
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ Variável obrigatória %s não definida", key)
	}
	return val
}
