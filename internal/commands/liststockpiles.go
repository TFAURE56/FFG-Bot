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

	// Changer le fuseau horaire
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Printf("�� Erreur lors du chargement du fuseau horaire : %v", err)
		return
	}
	time.Local = loc

	// Récupération des stockpiles depuis la base de données

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

	log.Printf("Connexion à la base de données réussie : %v", db)

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
		cooldownTime, err := time.Parse("2006-01-02 15:04:05", cooldownStr)
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
		// Calcul du temps restant
		cooldownStr := getTimeRemaining(sp.Cooldown)

		response += fmt.Sprintf("# 📦 **%s**\n🗺️ Hexagone: %s\n🔑 Code: ||`%s`||\n⏳ Temps restant: %s\n\n",
			sp.Nom, sp.Hexa, sp.Code, cooldownStr)
	}

	response += ("\n\n**Pour réinitialiser un stockpile, utilisez la commande `/resetstockpile <nom>`**")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Retourne le temps restant avant expiration du cooldown
func getTimeRemaining(expiration int64) string {
	now := time.Now().Unix()
	remaining := expiration - now

	if remaining <= 0 {
		return "✅ Disponible !"
	}

	hours := remaining / 3600
	minutes := (remaining % 3600) / 60
	seconds := remaining % 60

	return fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)
}
