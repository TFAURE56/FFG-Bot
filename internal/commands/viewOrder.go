package commands

import (
	"FFG-Bot/internal/global"
	"log"

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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Détails de la commande n° " + orderDetails.ID,
					Description: orderDetails.Comment,
					Color:       0x00FF00, // Vert
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Livraison pour : ", Value: orderDetails.EndDate, Inline: true},
					},
				},
			},
			Content: "✅ Détails de la commande récupérés avec succès.",
		},
	})
}

func viewOrderAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {

	choices := global.GetOrderIDsFromDB()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// OrderDetails structure to hold order information
type OrderDetails struct {
	ID         string
	Comment    string
	EndDate    string
	Orderer    string
	Ressources []OrderElement
}

// OrderElement structure to hold individual order elements
type OrderElement struct {
	Ressource string
	Number    int64
}
