package commands

import (
	//"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func RegisterHelpCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Affiche la liste des commandes disponibles et leurs descriptions",
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("❌ Impossible de créer la commande %s: %v", cmd.Name, err)
	} else {
		log.Printf("✅ Commande %s enregistrée avec succès", cmd.Name)
	}

	s.AddHandler(helpHandler)
}

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "help" {
		return
	}

	helpMessage := "**📜 Liste des commandes disponibles :**\n\n" +
		"📦 `/addstockpile <nom> <hexa> <code>` - Ajoute un stockpile avec un nom, un hexagone et un code d'accès.\n" +
		"📋 `/liststockpiles` - Affiche la liste des stockpiles du serveur avec leurs cooldowns.\n" +
		"⏳ `/resetstockpile <nom>` - Réinitialise le cooldown d'un stockpile à 49 heures.\n" +
		"🗑️ `/removestockpile <nom>` - Supprime un stockpile du serveur.\n" +
		"ℹ️ `/help` - Affiche ce message d'aide avec toutes les commandes disponibles.\n" +
		"‼️ `/setcooldownchannel` - Définit le salon pour les alertes de cooldown.\n"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMessage,
		},
	})
}
