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

	helpMessage := "# **📜 Liste des commandes disponibles :**\n\n" +
		" ## **🛠️ Commandes de gestion des stockpiles :**\n" +
		"📦 `/addstockpile <nom> <hexa> <code>` - Ajoute un stockpile avec un nom, un hexagone et un code d'accès.\n" +
		"📋 `/liststockpiles` - Affiche la liste des stockpiles du serveur avec leurs cooldowns.\n" +
		"⏳ `/resetstockpile <nom>` - Réinitialise le cooldown d'un stockpile à 48 heures.\n" +
		"🗑️ `/removestockpile <stockpile_name>` - Retire un stockpile par son nom.\n\n" +
		" ## **🗳️ Commandes de gestions des orders :**\n" +
		"🆕 `/neworder <description> <date de livraison(JJ/MM/AAA)>` - Ajoute un order une description et une date de livraison.\n" +
		"📑 `/vieworderson` - Affiche la liste des orders actifs sur le serveur.\n" +
		"📄 `/vieworder <order id>` - Affiche les informations d'un order spécifique.\n" +
		"➕ `/addorderelement <order id> <élément> <quantitée>` - Ajoute un élément à un order existant.\n" +
		"⤵️ `/getorderelement <order id> <élément>` - Vous assigne a un élément d'un order existant.\n" +
		"❔ `/viewlements` - Affiche les elements disponible a la commande.\n" +

		"\nℹ️ `/help` - Affiche ce message d'aide avec toutes les commandes disponibles.\n"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMessage,
		},
	})
}
