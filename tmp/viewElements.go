package commands

import (
	"FFG-Bot/internal/global"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "viewelements",
		Description: "Voir les éléments possibles pour les orders",
	}, viewElementsHandler)
}

func viewElementsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "viewelements" {
		return
	}

	db, err := global.ConnectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données :", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur de connexion à la base de données.",
			},
		})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT name, price, money, enable FROM elements")
	if err != nil {
		log.Println("Erreur lors de la récupération des éléments :", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Une erreur est survenue lors de la récupération des éléments.",
			},
		})
		return
	}

	defer rows.Close()

	var elements []string
	for rows.Next() {
		var name string
		var money string
		var price int
		var enable bool
		if err := rows.Scan(&name, &price, &money, &enable); err != nil {
			log.Println("Erreur lors du scan des éléments :", err)
			continue
		}

		// Si l'élément n'est pas activé, ajouter (Désactivé) à son nom
		if !enable {
			elements = append(elements, "⭕ "+name+" : "+strconv.Itoa(price)+" "+money)
		} else {
			elements = append(elements, "✔️ "+name+" : "+strconv.Itoa(price)+" "+money)
		}
	}

	if len(elements) == 0 {

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ℹ️ Aucun élément disponible pour le moment.",
			},
		})
		return
	}

	responseMessage := "# **📜 Liste des éléments disponibles pour les orders :**\n\n" +
		"\n " +
		"```\n" +
		strings.Join(elements, "\n") + " ``` \n" +
		"-# ⭕ : Élement non disponible pour le moment"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
		},
	})
}
