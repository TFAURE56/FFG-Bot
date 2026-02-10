package commands

import (
	"FFG-Bot/internal/global"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Autocomplétion pour afficher les commandes de production
func ViewOrderAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {

	choices := global.GetOrderIDsFromDB()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// Autocomplétion pour les deux options
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

// Autocomplétion pour afficher les stockpiles du serveur
func NameStockpileAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	// Vérifier que la commande est bien /resetstockpile ou /removestockpile
	if i.ApplicationCommandData().Name != "resetstockpile" && i.ApplicationCommandData().Name != "removestockpile" {
		return
	}

	db, err := global.ConnectToDatabase()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur de connexion à la base de données.",
			},
		})
		log.Println("Erreur de connexion à la base de données pour l'autocomplétion des stockpiles:", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM stockpiles")
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la récupération des stockpiles.",
			},
		})
		log.Println("Erreur lors de la récupération des stockpiles:", err)
		return
	}
	defer rows.Close()

	var options []*discordgo.ApplicationCommandOptionChoice
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: name,
			})
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: options,
		},
	})
}

// Autocomplétion pour le champ "type" de la commande /addstockpile
func AddStockpileTypeAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	if i.ApplicationCommandData().Name != "addstockpile" {
		return
	}

	// Détecter l'option en cours d'autocomplétion
	focused := ""
	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Focused {
			focused = opt.Name
			break
		}
	}

	if focused != "type" {
		return
	}

	choices := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Dépot", Value: "normal"},
		{Name: "Dépot Naval", Value: "bateau"},
		{Name: "Dépot Aérien", Value: "avion"},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}
