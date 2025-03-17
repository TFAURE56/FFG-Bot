package commands

import (
	"FFG-Bot/json"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RegisterAddstockpileCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "addstockpile",
		Description: "Ajoute un stockpile avec un nom, un hexagone et un code",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "nom",
				Description: "Nom du stockpile",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "hexa",
				Description: "Hexagone dans lequel se site le stockpile",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "Code d'accès du stockpile",
				Required:    true,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("❌ Impossible de créer la commande %s: %v", cmd.Name, err)
	} else {
		log.Printf("✅ Commande %s enregistrée avec succès", cmd.Name)
	}

	s.AddHandler(addStockpileHandler)
}

func addStockpileHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "addstockpile" {
		return
	}

	options := i.ApplicationCommandData().Options
	var nom, hexa, code string

	for _, option := range options {
		switch option.Name {
		case "nom":
			nom = option.StringValue()
		case "hexa":
			hexa = option.StringValue()
		case "code":
			code = option.StringValue()
		}
	}

	// Cooldown de 49 heures en timestamp UNIX
	cooldown := time.Now().Unix() + (49 * 3600)

	// Vérification des valeurs obligatoires
	if nom == "" || hexa == "" || code == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur : Tous les paramètres sont obligatoires.",
			},
		})
		return
	}

	// Ajout du stockpile dans le fichier JSON
	err := json.AddStockpile(i.GuildID, nom, hexa, code, cooldown)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de l'enregistrement du stockpile.",
			},
		})
		log.Printf("Erreur d'ajout du stockpile : %v", err)
		return
	}

	// Réponse de confirmation
	cooldownStr := time.Unix(cooldown, 0).Format("02/01/2006 15:04:05")
	response := fmt.Sprintf("📦 Stockpile **%s** ajouté à **%s** avec le code `%s`.\n⏳ Cooldown jusqu'à **%s**.", nom, hexa, code, cooldownStr)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
