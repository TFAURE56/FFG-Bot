package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(
		&discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Affiche la latence du bot",
		},
		pingHandler,
	)
}

// Handler pour la commande ping
func pingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Milliseconds()
	response := fmt.Sprintf("Pong! Latence: %dms\n Je suis a ma version 0.4.0", latency)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
