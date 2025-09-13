package commands

import (
	"FFG-Bot/internal/global"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "getorderelement",
		Description: "S'attribuer un élément d'une commande",
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
				Name:         "element",
				Description:  "Élément de la commande à s'attribuer",
				Required:     true,
				Autocomplete: true,
			},
		},
	}, getOrderElementHandler)
}

// Handler principal de la commande
func getOrderElementHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slave := i.Member.User.GlobalName

	if i.ApplicationCommandData().Name != "getorderelement" {
		return
	}

	orderID := i.ApplicationCommandData().Options[0].IntValue()
	element := i.ApplicationCommandData().Options[1].StringValue()

	// Ajouter le slave à l'élément de la commande dans la base de données
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
	_, err = db.Exec("UPDATE order_elements SET slave = ? WHERE order_id = ? AND ressource = ?", slave, orderID, element)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de l'attribution de l'élément de la commande.",
			},
		})

		return
	}

	// Répondre à l'utilisateur
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ L'élément de la commande a été attribué avec succès.",
		},
	})
}
