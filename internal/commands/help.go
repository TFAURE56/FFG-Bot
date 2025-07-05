package commands

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(
		&discordgo.ApplicationCommand{
			Name:        "help",
			Description: "Affiche la liste des commandes disponibles",
		},
		helpHandler,
	)
}

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "help" {
		return
	}

	helpMessage := "**📜 Liste des commandes disponibles :**\n\n" +
		"📦 `/addstockpile <nom> <hexa> <code>` - Ajoute un stockpile avec un nom, un hexagone et un code d'accès.\n" +
		"📋 `/liststockpiles` - Affiche la liste des stockpiles du serveur avec leurs cooldowns.\n" +
		"⏳ `/resetstockpile <nom>` - Réinitialise le cooldown d'un stockpile à 48 heures.\n" +
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
