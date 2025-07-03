package main

import (
	"FFG-Bot/internal/commands"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	token := os.Getenv("DISCORD_BOT_TOKEN")
	clientID := os.Getenv("ClientID")
	guildID := os.Getenv("GuildID")

	// Vérification des variables d'environnement
	if token == "" || clientID == "" || guildID == "" {
		log.Fatal("Variables d'environnement manquantes. Vérifiez DISCORD_BOT_TOKEN, ClientID et GuildID.")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la session Discord: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		commands.Start(dg)
	}()

	//go func() {
	//	defer wg.Done()
	//	logs.Start(dg)
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//	experiences.Start(dg)
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//	routines.Start(dg)
	//}()

	// Attente d'une interruption pour arrêter le bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	dg.Close()
	wg.Wait()
	log.Println("Bot arrêté.")
}
