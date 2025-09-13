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
				// Récupérer la liste des ressources depuis la base de données
				Choices: func() []*discordgo.ApplicationCommandOptionChoice {
					ressources, err := getRessourcesList()
					if err != nil {
						log.Printf("Erreur lors de la récupération des ressources : %v", err)
						return []*discordgo.ApplicationCommandOptionChoice{}
					}
					choices := make([]*discordgo.ApplicationCommandOptionChoice, len(ressources))
					for i, ressource := range ressources {
						choices[i] = &discordgo.ApplicationCommandOptionChoice{
							Name:  ressource,
							Value: ressource,
						}
					}
					return choices
				}(),
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

	// Si l'user a choisi une ressource non disponible, on bloque l'ajout
	value := i.ApplicationCommandData().Options[1].StringValue()
	if value != "" && len(value) >= 17 && value[len(value)-17:] == " (non disponible)" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ La ressource choisie n'est pas disponible à la commande.",
			},
		})
		return
	}

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

// Fonction pour récupérer la liste des ressources dans la base de données
func getRessourcesList() ([]string, error) {
	db, err := global.ConnectToDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, enable FROM elements")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ressources []string
	for rows.Next() {
		var name string
		var enable int
		err := rows.Scan(&name, &enable)
		if err != nil {
			return nil, err
		}
		// Si la ressource est activée, on l'ajoute à la liste
		if enable == 1 {
			ressources = append(ressources, name)
		} else {
			ressources = append(ressources, name+" (non disponible)")
		}
	}
	return ressources, nil
}
