package commands

import (
	"FFG-Bot/internal/global"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "addorderelement",
		Description: "Ajouter un élément à une commande",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionInteger,
				Name:         "order_id",
				Description:  "Numéro de la commande",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "ressource",
				Description:  "Ressource à ajouter à la commande",
				Required:     true,
				Autocomplete: false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "quantite",
				Description: "Quantité de la ressource à ajouter",
				Required:    true,
			},
		},
	}, addOrderElementHandler)
}

func addOrderElementHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Récupération des options de la commande
	var orderID int64
	var ressource string
	var quantite int64

	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "order_id":
			orderID = opt.IntValue()
		case "ressource":
			ressource = opt.StringValue()
		case "quantite":
			quantite = opt.IntValue()
		}
	}

	// Ajout de l'élément à la base de données
	db, err := global.ConnectToDatabase()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Erreur de connexion à la base de données : " + err.Error(),
			},
		})
		return
	}
	defer db.Close()

	log.Printf("Ajout de l'élément à la commande %d : %s (Quantité: %d)", orderID, ressource, quantite)

	_, err = db.Exec("INSERT INTO order_elements (order_id, ressource, number) VALUES (?, ?, ?)",
		orderID, ressource, quantite)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Erreur lors de l'ajout de l'élément à la commande : " + err.Error(),
			},
		})
		return
	}

	response := "Élément ajouté à la commande " + strconv.FormatInt(orderID, 10) + ": " + ressource + " (Quantité: " + strconv.FormatInt(quantite, 10) + ")"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Autocomplete pour les options de la commande
func addOrderElementAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	choices := global.GetOrderIDsFromDB()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
