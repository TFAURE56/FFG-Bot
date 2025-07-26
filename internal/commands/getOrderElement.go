package commands

import (
	"FFG-Bot/internal/global"
	"strconv"

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

// Handler d'autocomplétion pour les deux options
func getOrderElementAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

		elements := getElementsForOrder(orderID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: elements,
			},
		})
	}
}

// Récupère la liste des éléments pour un ordre donné
func getElementsForOrder(orderID int64) []*discordgo.ApplicationCommandOptionChoice {
	db, err := global.ConnectToDatabase()
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT ressource FROM order_elements WHERE order_id = ?", orderID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var choices []*discordgo.ApplicationCommandOptionChoice
	for rows.Next() {
		var element string
		if err := rows.Scan(&element); err == nil {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  element,
				Value: element,
			})
		}
	}
	return choices
}
