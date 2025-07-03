package commands

import (
	"fmt"
	"log"
	"time"

	"FFG-Bot/internal/global"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
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
				Description: "Hexagone du stockpile",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "Code d'accès au stockpile",
				Required:    true,
			},
		},
	}, addStockpileHandler)
}

// Ajout d'un stockpile en base de donnée
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

	// Ajout du stockpile dans la base de données Mariadb
	db, err := global.ConnectToDatabase()
	log.Printf("Connexion à la base de données : %v", db)
	if err != nil {
		log.Printf("❌ Erreur de connexion à la base de données : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur de connexion à la base de données.",
			},
		})
		return
	}
	// Calcul de la date d'expiration du cooldown (+48h)
	cooldownExpiration := time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04:05")

	_, err = db.Exec(
		"INSERT INTO stockpiles (name, hexa, code, cooldown) VALUES (?, ?, ?, ?)",
		nom, hexa, code, cooldownExpiration,
	)
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
	defer db.Close()

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
