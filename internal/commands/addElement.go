package commands

import (
	"FFG-Bot/internal/global"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "addelement",
		Description: "Ajouter un élément possible a la commande",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "element",
				Description:  "Nom de l'element à ajouter",
				Required:     true,
				Autocomplete: false,
			},
			{
				Type:         discordgo.ApplicationCommandOptionInteger,
				Name:         "prix",
				Description:  "Prix de l'element à ajouter",
				Required:     true,
				Autocomplete: false,
			},
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "money",
				Description:  "Type de monnaie (Souffre, )",
				Required:     true,
				Autocomplete: false,
			},
		},
	}, addElementHandler)
}

func addElementHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "addelement" {
		return
	}

	options := i.ApplicationCommandData().Options
	var element, money string
	var prix int64

	for _, option := range options {
		switch option.Name {
		case "element":
			element = option.StringValue()
		case "prix":
			prix = option.IntValue()
		case "money":
			money = option.StringValue()
		}
	}

	log.Printf("Ajout de l'élément : %s, Prix : %d, Monnaie : %s", element, prix, money)

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

	// Insérer l'élément dans la base de données
	_, err = db.Exec("INSERT INTO elements (name, price, money) VALUES (?, ?, ?)", element, prix, money)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de l'ajout de l'élément.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ L'élément **" + element + "** avec le prix **" + strconv.FormatInt(prix, 10) + " " + money + "** a été ajouté avec succès.",
		},
	})
	db.Close()
}
