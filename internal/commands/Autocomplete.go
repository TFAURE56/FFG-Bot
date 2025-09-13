package commands

import (
	"FFG-Bot/internal/global"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func ViewOrderAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {

	choices := global.GetOrderIDsFromDB()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// Handler d'autocomplétion pour les deux options
func GetOrderElementAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}
	if i.ApplicationCommandData().Name != "getorderelement" {
		return
	}

	// Détecter quelle option est en cours d'autocomplétion
	focused := ""
	var orderID int64
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Focused {
			focused = opt.Name
		}
		if opt.Name == "order_id" {
			// Correction : gérer le type selon le contexte (float64 pour autocomplete, int64 pour commande)
			switch v := opt.Value.(type) {
			case float64:
				orderID = int64(v)
			case int:
				orderID = int64(v)
			case int64:
				orderID = v
			case string:
				// Si jamais Discord renvoie un string (peu probable ici), essayer de parser
				if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
					orderID = parsed
				}
			}
		}
	}

	switch focused {
	case "order_id":
		choices := global.GetOrderIDsFromDB()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	case "element":

		elements := global.GetElementsForOrder(orderID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: elements,
			},
		})
	}
}
