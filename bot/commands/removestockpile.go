package commands

import (
	"FFG-Bot/json"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func RegisterRemoveStockpileCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "removestockpile",
		Description: "Supprime un stockpile du serveur",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "nom",
				Description:  "Nom du stockpile à supprimer",
				Required:     true,
				Autocomplete: true, // Active l'autocomplétion
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("❌ Impossible de créer la commande %s: %v", cmd.Name, err)
	} else {
		log.Printf("✅ Commande %s enregistrée avec succès", cmd.Name)
	}

	// Ajout des handlers
	s.AddHandler(removeStockpileHandler)
	s.AddHandler(removeStockpileAutocomplete)
}

// Autocomplétion pour afficher les stockpiles du serveur
func removeStockpileAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	if i.ApplicationCommandData().Name != "removestockpile" {
		return
	}

	stockpiles, err := json.GetStockpiles(i.GuildID)
	if err != nil || len(stockpiles) == 0 {
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, sp := range stockpiles {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  sp.Nom,
			Value: sp.Nom,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// Handler pour la commande /removestockpile
func removeStockpileHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "removestockpile" {
		return
	}

	nom := i.ApplicationCommandData().Options[0].StringValue()
	stockpiles, err := json.GetStockpiles(i.GuildID)
	if err != nil || len(stockpiles) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Aucun stockpile trouvé pour ce serveur.",
			},
		})
		return
	}

	var newStockpiles []json.Stockpiles
	found := false

	for _, sp := range stockpiles {
		if sp.Nom != nom {
			newStockpiles = append(newStockpiles, sp)
		} else {
			found = true
		}
	}

	if !found {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Le stockpile **%s** n'existe pas.", nom),
			},
		})
		return
	}

	// Sauvegarde des stockpiles mis à jour
	err = json.SaveStockpiles(i.GuildID, newStockpiles)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la suppression du stockpile.",
			},
		})
		log.Printf("Erreur lors de la suppression du stockpile : %v", err)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🗑️ Le stockpile **%s** a été supprimé avec succès.", nom),
		},
	})
}
