package commands

import (
	"FFG-Bot/internal/global"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "neworder",
		Description: "Crée une nouvelle commande de production",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "description",
				Description: "Description de la commande (Objectif de la production / Lieu où y stocker ...)",
				Required:    true,
			},
			// Option pour la date de livraison
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "date_livraison",
				Description: "Date de livraison de la commande (format AAAA-MM-JJ HH:MM)",
				Required:    true,
			},
		},
	},
		newOrderHandler,
	)
}

func newOrderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "neworder" {
		return
	}

	description := i.ApplicationCommandData().Options[0].StringValue()
	dateLivraison := i.ApplicationCommandData().Options[1].StringValue()
	orderer := i.Member.User.Username

	log.Printf("Nouvelle commande de production demandée par %s : %s, Date de livraison : %s", orderer, description, dateLivraison)

	// Enregistre la commande dans la base de données
	db, err := global.ConnectToDatabase()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur de connexion à la base de données.",
			},
		})
		return
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO orders (comment, orderer, end_date, working) VALUES (?, ?, ?, ?)", description, orderer, dateLivraison, "on")
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de l'enregistrement de la commande.",
			},
		})
		log.Printf("❌ Erreur lors de l'enregistrement de la commande : %v", err)
		return
	}

	// Récupérer l'ID de la commande créée
	var orderID int64
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&orderID)
	log.Printf("Nouvelle commande créée avec ID : %d", orderID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la récupération de l'ID de la commande.",
			},
		})
		return
	}

	response := "**ID de la commande :** " + strconv.FormatInt(orderID, 10) + "\n" +
		"**Description :** " + description + "\n" +
		"**Date de livraison :** " + dateLivraison
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral, // Message visible uniquement par l'utilisateur
		},
	})
}
