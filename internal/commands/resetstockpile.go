package commands

import (
	"FFG-Bot/internal/global"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "resetstockpile",
		Description: "Remet le cooldown d'un stockpile à 48h",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "nom",
				Description:  "Nom du stockpile à reset",
				Required:     true,
				Autocomplete: true,
			},
		},
	}, resetStockpileHandler)
}

// Handler pour la commande /resetstockpile
func resetStockpileHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "resetstockpile" {
		return
	}

	nom := i.ApplicationCommandData().Options[0].StringValue()

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
	defer db.Close()

	// Vérifier si le stockpile existe
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM stockpiles WHERE name = ?", nom).Scan(&count)
	if err != nil || count == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Le stockpile **%s** n'existe pas.", nom),
			},
		})
		return
	}

	// Mettre à jour le cooldown à maintenant + 48h
	newCooldown := time.Now().Add(48 * time.Hour)
	_, err = db.Exec("UPDATE stockpiles SET cooldown = ? WHERE name = ?", newCooldown, nom)
	if err != nil {
		log.Printf("❌ Erreur lors du reset du cooldown : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la mise à jour du stockpile.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🔄 Le cooldown du stockpile **%s** a été remis à **48 heures**.", nom),
		},
	})
}

// Reset un stockpile depuis un bouton (MessageComponent)
func ResetStockpileByButton(s *discordgo.Session, i *discordgo.InteractionCreate, name, hexa string) {
	db, err := global.ConnectToDatabase()
	if err != nil {
		log.Printf("Erreur de connexion à la base de données: %v", err)
		return
	}
	defer db.Close()

	// Reset le cooldown à 48h
	newCooldown := time.Now().Add(48 * time.Hour)
	_, err = db.Exec("UPDATE stockpiles SET cooldown = ?, alerted = 0 WHERE name = ? AND hexa = ?", newCooldown, name, hexa)
	if err != nil {
		log.Printf("Erreur lors du reset du stockpile %s (hexa: %s): %v", name, hexa, err)
	} else {
		log.Printf("Stockpile %s (hexa: %s) reset par bouton.", name, hexa)
	}

	// Réponse à l'utilisateur (déjà faite dans main.go, donc optionnel ici)
}
