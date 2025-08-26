package commands

import (
	"FFG-Bot/internal/global"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
				Description: "Date de livraison de la commande (format JJ/MM/AAAA HH:MM)",
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
	orderer := i.Member.User.GlobalName

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

	// réorganiser la date de livraison
	dateLivraison = reorganizeDate(dateLivraison)

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors du chargement de la localisation.",
			},
		})
		return
	}

	log.Printf("Date de livraison avant conversion: %s", dateLivraison)

	dateLivraisonFR, err := time.ParseInLocation("2006-01-02 15:04:05", dateLivraison, loc)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la conversion de la date de livraison.",
			},
		})
		return
	}

	dateLivraisonUTC := dateLivraisonFR.UTC()

	log.Printf("Nouvelle commande de %s : %s, date de livraison : %s", orderer, description, dateLivraisonUTC)

	_, err = db.Exec("INSERT INTO orders (comment, orderer, end_date, working) VALUES (?, ?, ?, ?)", description, orderer, dateLivraisonUTC, "on")
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

func reorganizeDate(date string) string {
	// Supposons que la date soit au format "DD/MM/YYYY HH:MM"
	// On la réorganise au format "YYYY-MM-DD HH:MM:SS"
	parts := strings.Split(date, " ")
	if len(parts) != 2 {
		return date // Retourne la date originale si le format n'est pas correct
	}
	dateParts := strings.Split(parts[0], "/")
	if len(dateParts) != 3 {
		return date // Retourne la date originale si le format n'est pas correct
	}
	return fmt.Sprintf("%s-%s-%s %s:00", dateParts[2], dateParts[1], dateParts[0], parts[1])
}
