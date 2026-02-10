package commands

import (
	"FFG-Bot/internal/global"
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "removestockpile",
		Description: "Retirer un stockpile",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "stockpile_name",
				Description:  "Nom du stockpile à retirer",
				Required:     true,
				Autocomplete: true,
			},
		},
	}, removeStockpileHandler)
}

// Handler pour la commande /removestockpile
func removeStockpileHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "removestockpile" {
		return
	}

	stockpileName := i.ApplicationCommandData().Options[0].StringValue()
	log.Printf("Tentative de suppression du stockpile : %s", stockpileName)

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

	// Supprimer le stockpile de la base de données
	_, err = db.Exec("DELETE FROM stockpiles WHERE name = ?", stockpileName)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la suppression du stockpile.",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ Le stockpile avec le nom **" + stockpileName + "** a été supprimé avec succès.",
		},
	})
	db.Close()
}
