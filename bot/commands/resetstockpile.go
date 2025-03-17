package commands

import (
	"FFG-Bot/json"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RegisterResetStockpileCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "resetstockpile",
		Description: "Remet le cooldown d'un stockpile à 49h",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "nom",
				Description:  "Nom du stockpile à reset",
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
	s.AddHandler(resetStockpileHandler)
	s.AddHandler(resetStockpileAutocomplete)
}

// Autocomplétion pour afficher les stockpiles du serveur
func resetStockpileAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	if i.ApplicationCommandData().Name != "resetstockpile" {
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
		Type: discordgo.InteractionApplicationCommandAutocompleteResult, // ✅ Correction ici
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

// Handler pour la commande /resetstockpile
func resetStockpileHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "resetstockpile" {
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

	found := false
	for idx, sp := range stockpiles {
		if sp.Nom == nom {
			stockpiles[idx].Cooldown = time.Now().Unix() + (49 * 3600)
			found = true
			break
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

	// ✅ Correction : Utilisation de SaveStockpiles
	err = json.SaveStockpiles(i.GuildID, stockpiles)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la mise à jour du stockpile.",
			},
		})
		log.Printf("Erreur lors du reset du stockpile : %v", err)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🔄 Le cooldown du stockpile **%s** a été remis à **49 heures**.", nom),
		},
	})
}
