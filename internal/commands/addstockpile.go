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
				Type:        discordgo.ApplicationCommandOptionType(3), // discordgo.ApplicationCommandOptionString est égal à 3, mais on utilise la valeur brute pour éviter les problèmes de version
				Name:        "hexa",
				Description: "Hexagone du stockpile",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ville",
				Description: "Ville du stockpile",
				Required:    true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "type",
				Description:  "Type du stockpile",
				Required:     true,
				Autocomplete: true,
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
	var nom, hexa, ville, style, code string

	for _, option := range options {
		switch option.Name {
		case "nom":
			nom = option.StringValue()
		case "hexa":
			hexa = option.StringValue()
		case "ville":
			ville = option.StringValue()
		case "type":
			style = option.StringValue()
		case "code":
			code = option.StringValue()
		}

	}

	// Calcul de la date d'expiration du cooldown (+48h) en UTC
	cooldownExpiration := time.Now().UTC().Add(48 * time.Hour)

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

	// Vérifier si le stockpile existe déjà
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM stockpiles WHERE name = ?", nom).Scan(&count)
	if err != nil {
		log.Printf("❌ Erreur lors de la vérification de l'existence du stockpile : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la vérification du stockpile.",
			},
		})
		return
	}
	if count > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Le stockpile **%s** existe déjà.", nom),
			},
		})
		return
	}

	_, err = db.Exec(
		"INSERT INTO stockpiles (name, hexa, ville, style, code, cooldown) VALUES (?, ?, ?, ?, ?, ?)",
		nom, hexa, ville, style, code, cooldownExpiration.Format("2006-01-02 15:04:05"),
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
	loc, err := time.LoadLocation("Europe/Paris")

	// Définir l'emoji en fonction du type de stockpile
	var response string
	if style == "normal" {
		response = fmt.Sprintf("📦 Stockpile **%s** ajouté à **%s** dans la ville de **%s** avec le code `%s`.\n⏳ Cooldown jusqu'à **%s** (heure Paris).",
			nom, hexa, ville, code, time.Unix(int64(cooldownExpiration.Unix()), 0).In(loc).Format("02/01/2006 15:04:05"))
	} else if style == "bateau" {
		response = fmt.Sprintf("⚓ Stockpile **%s** ajouté à **%s** dans la ville de **%s** avec le code `%s`.\n⏳ Cooldown jusqu'à **%s** (heure Paris).",
			nom, hexa, ville, code, time.Unix(int64(cooldownExpiration.Unix()), 0).In(loc).Format("02/01/2006 15:04:05"))
	} else if style == "avion" {
		response = fmt.Sprintf("✈️ Stockpile **%s** ajouté à **%s** dans la ville de **%s** avec le code `%s`.\n⏳ Cooldown jusqu'à **%s** (heure Paris).",
			nom, hexa, ville, code, time.Unix(int64(cooldownExpiration.Unix()), 0).In(loc).Format("02/01/2006 15:04:05"))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
