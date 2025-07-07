package main

import (
	"FFG-Bot/internal/commands"
	"FFG-Bot/internal/global"
	"FFG-Bot/internal/routines"
	"fmt"
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

	// Connexion à la base de données
	db, err := global.ConnectToDatabase()
	if err != nil {
		log.Fatalf("❌ Erreur de connexion à la base de données : %v", err)
	}
	defer db.Close()

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
	go func() {
		defer wg.Done()
		routines.Start(dg, db)
	}()

	// Ajoutez ce handler lors de l'initialisation du bot Discord
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			if len(i.MessageComponentData().CustomID) > 16 && i.MessageComponentData().CustomID[:15] == "reset_stockpile" {
				customID := i.MessageComponentData().CustomID
				data := customID[16:]
				lastUnderscore := -1
				for i := len(data) - 1; i >= 0; i-- {
					if data[i] == '_' {
						lastUnderscore = i
						break
					}
				}
				var name, hexa string
				if lastUnderscore != -1 {
					name = data[:lastUnderscore]
					hexa = data[lastUnderscore+1:]
				} else {
					name = data
					hexa = ""
				}
				// Appel correct pour un bouton
				commands.ResetStockpileByButton(s, i, name, hexa)
				log.Printf("Reset du stockpile %s (hexa: %s) demandé par %s", name, hexa, i.Member.User.Username)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Le stockpile **%s** (hexa: %s) a été reset !", name, hexa),
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				//Supprimer le message de l'alerte cooldown
				if i.Message != nil {
					err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
					if err != nil {
						log.Printf("Erreur lors de la suppression du message d'alerte cooldown: %v", err)
					} else {
						log.Printf("Message d'alerte cooldown supprimé pour le stockpile %s (hexa: %s)", name, hexa)
					}
				}
			}
		}
	})

	// Lancer la routine cooldown avec la session Discord et la connexion DB
	routines.StartCooldown(dg, db)

	// Attente d'une interruption pour arrêter le bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	dg.Close()
	wg.Wait()
	log.Println("Bot arrêté.")
}
