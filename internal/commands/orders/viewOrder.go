package commands

import (
	"FFG-Bot/internal/global"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "vieworder",
		Description: "Voir les détails d'une commande",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionInteger,
				Name:         "order_id",
				Description:  "Numéro de la commande",
				Required:     true,
				Autocomplete: true,
			},
		},
	}, viewOrderHandler)
}

// Handler principal de la commande
func viewOrderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	orderID := i.ApplicationCommandData().Options[0].IntValue()

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

	log.Println("Récupération des détails de la commande ID:", orderID)

	var orderDetails OrderDetails
	err = db.QueryRow("SELECT id, comment, end_date, orderer FROM orders WHERE id = ?", orderID).Scan(&orderDetails.ID, &orderDetails.Comment, &orderDetails.EndDate, &orderDetails.Orderer)
	if err != nil {
		log.Println("Erreur lors de la récupération des détails de la commande:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la récupération des détails de la commande.",
			},
		})
		return
	}

	// Récupérer les elements associé à la commande
	rows, err := db.Query("SELECT ressource, number, slave FROM order_elements WHERE order_id = ?", orderID)
	if err != nil {
		log.Println("Erreur lors de la récupération des éléments de la commande:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la récupération des éléments de la commande.",
			},
		})
		return
	}
	defer rows.Close()
	var elements []OrderElements
	for rows.Next() {
		var element OrderElements
		err := rows.Scan(&element.Ressource, &element.Number, &element.Slave)
		if err != nil {
			log.Println("Erreur lors de la lecture des éléments de la commande:", err)
			continue
		}
		if !element.Slave.Valid {
			element.Slave.String = "Non assigné"
			element.Slave.Valid = true
		}

		elements = append(elements, element)
	}
	orderDetails.Ressources = elements

	// Transforme la date de livraison en format lisible
	orderDetails.EndDateString = orderDetails.EndDate.Format("02/01/2006")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "__Détails de la commande n° " + orderDetails.ID + "__",
					Description: "**Description: \n**" + orderDetails.Comment,
					Color:       0x00FF00, // Vert
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Livraison pour : ", Value: orderDetails.EndDateString, Inline: true},
						{Name: "Commandé par : ", Value: orderDetails.Orderer, Inline: true},
						// Champs pour les elements de la commande
						{Name: "Éléments de la commande :", Value: formatOrderElements(orderDetails.Ressources), Inline: false},
					},
				},
			},
			Content: "✅ Détails de la commande récupérés avec succès.",
		},
	})
}

func formatOrderElements(elements []OrderElements) string {
	if len(elements) == 0 {
		return "Aucun élément associé à cette commande.\nAjoutez-en avec **/addorderelement**."
	}

	// Formate les éléments de la commande en une chaîne de caractères

	var result string
	for _, element := range elements {
		slaveStr := "n/a"
		if element.Slave.Valid {
			slaveStr = element.Slave.String
		}

		result += fmt.Sprintf("⭕ %s\t: %dㅤ- %s\n", element.Ressource, element.Number, slaveStr)
	}
	return result

}

// OrderDetails structure to hold order information
type OrderDetails struct {
	ID            string
	Comment       string
	EndDate       time.Time
	EndDateString string
	Orderer       string
	Ressources    []OrderElements
}

// OrderElement structure to hold individual order elements
type OrderElements struct {
	Ressource string
	Number    int64
	Slave     sql.NullString
}
