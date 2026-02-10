package commands

import (
	"FFG-Bot/internal/global"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "vieworderon",
		Description: "Voir les commandes en cours",
	}, viewOrderOnHandler)
}

// Handler principal de la commande
func viewOrderOnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Récupérer les commandes en cours (logique à implémenter)
	db, err := global.ConnectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur de connexion à la base de données.",
			},
		})
		return
	}
	defer db.Close()

	// récupérer les commandes en cours
	rows, err := db.Query("SELECT id, comment, end_date, orderer FROM orders WHERE working = 'on'")
	if err != nil {
		log.Println("Erreur lors de la récupération des commandes en cours:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la récupération des commandes en cours.",
			},
		})
		return
	}
	defer rows.Close()

	var response string
	for rows.Next() {
		var orderDetails OrderDetails
		err := rows.Scan(&orderDetails.ID, &orderDetails.Comment, &orderDetails.EndDate, &orderDetails.Orderer)
		if err != nil {
			log.Println("Erreur lors de la lecture des détails de la commande:", err)
			continue
		}
		response += formatOrderDetails(orderDetails) + "\n\n"
	}

	if response == "" {
		response = "Aucune commande en cours."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func formatOrderDetails(order OrderDetails) string {
	return fmt.Sprintf("# **Commande ID:** %s\n**Commentaire:** %s\n**Date de fin:** %s\n**Commandé par:** %s\n", order.ID, order.Comment, order.EndDate.Format("2006-01-02"), order.Orderer)
}
