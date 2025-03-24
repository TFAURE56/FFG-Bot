package settings

import (
	"FFG-Bot/json"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Enregistre la commande
func RegisterSetCooldownChannelCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "setcooldownchannel",
		Description: "Définit le salon où les alertes de cooldown seront envoyées",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Salon où envoyer les alertes",
				Required:    true,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("❌ Impossible de créer la commande %s: %v", cmd.Name, err)
	} else {
		log.Printf("✅ Commande %s enregistrée avec succès", cmd.Name)
	}

	s.AddHandler(setCooldownChannelHandler)
}

// Gère l'exécution de la commande
func setCooldownChannelHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "setcooldownchannel" {
		return
	}

	// Vérifier si l'utilisateur est admin
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("❌ Erreur lors de la récupération du membre :", err)
		return
	}

	// Vérifier les permissions ADMINISTRATOR
	isAdmin := false
	for _, roleID := range member.Roles {
		role, err := s.State.Role(i.GuildID, roleID)
		if err == nil && (role.Permissions&discordgo.PermissionAdministrator) != 0 {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "🚫 Vous devez être administrateur pour exécuter cette commande.",
				Flags:   discordgo.MessageFlagsEphemeral, // Message visible uniquement par l'auteur
			},
		})
		return
	}

	// Récupérer l'option channel
	channelID := i.ApplicationCommandData().Options[0].ChannelValue(s).ID

	// Sauvegarder le channel ID dans le fichier JSON
	err = json.SaveCooldownChannel(i.GuildID, channelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Erreur lors de la sauvegarde du salon.",
			},
		})
		log.Printf("Erreur d'enregistrement du salon cooldown : %v", err)
		return
	}

	// Réponse de confirmation
	response := fmt.Sprintf("✅ Le salon des alertes de cooldown est maintenant défini sur <#%s>.", channelID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
