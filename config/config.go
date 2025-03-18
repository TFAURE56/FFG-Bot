package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token   string
	GuildID string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur de chargement du fichier .env")
	}
	return Config{
		Token: os.Getenv("DISCORD_TOKEN"),
	}
}
