package commands

import (
	"FFG-Bot/internal/global"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(
		&discordgo.ApplicationCommand{
			Name:        "liststockpiles",
			Description: "Liste tous les stockpiles enregistrés pour ce serveur",
		},
		listStockpilesHandler,
	)
}

func listStockpilesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "liststockpiles" {
		return
	}

	// Chargement du fuseau horaire pour la France (Europe/Paris)
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Printf("❌ Erreur lors du chargement du fuseau horaire : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors du chargement du fuseau horaire.",
			},
		})
		return
	}

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

	type Stockpile struct {
		Nom      string
		Hexa     string
		Code     string
		Cooldown int64
	}
	var stockpiles []Stockpile

	rows, err := db.Query("SELECT name, hexa, code, cooldown FROM stockpiles")
	if err != nil {
		log.Printf("❌ Erreur lors de la récupération des stockpiles : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la récupération des stockpiles.",
			},
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sp Stockpile
		var cooldownStr string
		if err := rows.Scan(&sp.Nom, &sp.Hexa, &sp.Code, &cooldownStr); err != nil {
			log.Printf("❌ Erreur lors du scan d'un stockpile : %v", err)
			continue
		}

		// Parse toujours en heure de Paris
		cooldownTime, err := time.ParseInLocation("2006-01-02 15:04:05", cooldownStr, time.UTC)
		if err != nil {
			// Fallback RFC3339 si jamais tu as du ISO en base
			cooldownTime, err = time.ParseInLocation(time.RFC3339, cooldownStr, loc)
		}
		if err != nil {
			log.Printf("❌ Erreur lors du parsing du cooldown : %v", err)
			sp.Cooldown = 0
		} else {
			sp.Cooldown = cooldownTime.Unix()
		}
		stockpiles = append(stockpiles, sp)
	}
	if err := rows.Err(); err != nil {
		log.Printf("❌ Erreur lors de l'itération des stockpiles : %v", err)
	}

	defer db.Close()

	var response string
	for _, sp := range stockpiles {
		expirationParis := time.Unix(sp.Cooldown, 0).In(loc).Format("02/01/2006 15:04:05")
		cooldownStr := getTimeRemaining(sp.Cooldown)

		response += fmt.Sprintf(
			"# 📦 **%s**\n🗺️ Hexagone: %s\n🔑 Code: ||`%s`||\n⏳ Temps restant: %s\n🕒 Expire le: %s (heure de Paris)\n\n",
			sp.Nom, sp.Hexa, sp.Code, cooldownStr, expirationParis)
	}

	response += ("\n**Pour réinitialiser un stockpile, utilisez la commande `/resetstockpile <nom>`**")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Retourne le temps restant avant expiration du cooldown en heure de Paris
func getTimeRemaining(expiration int64) string {
	now := time.Now().UTC().Unix()

	remaining := expiration - now
	if remaining <= 0 {
		return "Stockpile potentiellement expiré."
	}

	hours := remaining / 3600
	minutes := (remaining % 3600) / 60
	seconds := remaining % 60

	return fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)
}
