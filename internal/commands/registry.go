package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	Registry      []*discordgo.ApplicationCommand
	Handlers      = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	ModalHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func Register(cmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	Registry = append(Registry, cmd)
	if handler != nil {
		Handlers[cmd.Name] = handler
	}
}

// Pour enregistrer un handler de modal
func RegisterModal(customID string, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	ModalHandlers[customID] = handler
}

func Start(dg *discordgo.Session) {

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if handler, ok := Handlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		}
		// Pour les modals
		HandleModal(s, i)
	})

	dg.AddHandler(resetStockpileAutocomplete)
	dg.AddHandler(viewOrderAutocomplete)
	dg.AddHandler(getOrderElementAutocomplete)
	//dg.AddHandler()

	err := dg.Open()
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la session Discord: %v", err)
	}

	fmt.Println("Bot en ligne. Appuyez sur CTRL+C pour quitter.")

	// Les commandes sont déjà enregistrées dans commands.Registry grâce aux init()
	cmds := Registry

	ClientID := os.Getenv("ClientID")
	GuildID := os.Getenv("GuildID")
	if ClientID == "" || GuildID == "" {
		log.Fatal("ClientID ou GuildID manquant dans les variables d'environnement.")
		return
	}

	// Lister les commandes dans cmds
	for _, cmd := range cmds {
		log.Printf("Commande enregistrée : %s (ID: %s)\n", cmd.Name, cmd.ID)
	}

	// Enregistrement des commandes
	_, err = dg.ApplicationCommandBulkOverwrite(ClientID, GuildID, cmds)
	if err != nil {
		log.Fatalf("Erreur lors de l'enregistrement des commandes : %v", err)
		return
	}

	log.Println("Module commandes démarré.")
}

func HandleModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionModalSubmit {
		if handler, ok := ModalHandlers[i.ModalSubmitData().CustomID]; ok {
			handler(s, i)
		}
	}
}
