package main

import (
	"FFG-Bot/bot"
	"FFG-Bot/config"
	"FFG-Bot/json"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	// Charger les salons configurés
	err := json.LoadSettings()
	if err != nil {
		log.Println("⚠️ Impossible de charger le fichier de paramètres, un nouveau sera créé.")
	}

	bot.Start(cfg.Token)
}
