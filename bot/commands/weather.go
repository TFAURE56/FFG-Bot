package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

const weatherAPIURL = "https://api.openweathermap.org/data/2.5/weather"

// Structure pour parser la réponse de l'API météo
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func RegisterWeatherCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "weather",
		Description: "Affiche la météo pour une ville donnée",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "city",
				Description: "Nom de la ville",
				Required:    true,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("❌ Impossible de créer la commande %s: %v", cmd.Name, err)
	}

	// Enregistrer la commande avec son handler
	RegisterCommand(s, guildID, cmd, weatherHandler)
}

func weatherHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "weather" {
		return
	}

	city := i.ApplicationCommandData().Options[0].StringValue()
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Println("❌ La clé API pour OpenWeatherMap n'est pas définie.")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ La clé API pour OpenWeatherMap n'est pas configurée.",
			},
		})
		return
	}

	// Construire l'URL de la requête
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=fr", weatherAPIURL, city, apiKey)

	// Faire la requête HTTP
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("❌ Erreur lors de la requête météo : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Impossible de récupérer les informations météo.",
			},
		})
		return
	}
	defer resp.Body.Close()

	// Décoder la réponse JSON
	var weatherData WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		log.Printf("❌ Erreur lors du décodage de la réponse météo : %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors du traitement des données météo.",
			},
		})
		return
	}

	// Construire l'embed avec les informations météo
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Météo pour %s", weatherData.Name),
		Description: fmt.Sprintf("**%s**", weatherData.Weather[0].Description),
		Color:       0x1E90FF, // Bleu
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Température",
				Value:  fmt.Sprintf("%.1f°C", weatherData.Main.Temp),
				Inline: true,
			},
			{
				Name:   "Humidité",
				Value:  fmt.Sprintf("%d%%", weatherData.Main.Humidity),
				Inline: true,
			},
		},
	}

	// Répondre avec l'embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
