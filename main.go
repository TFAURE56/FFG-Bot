package main

import (
	"FFG-Bot/bot"
	"FFG-Bot/config"
)

func main() {
	cfg := config.LoadConfig()
	bot.Start(cfg.Token)
}
