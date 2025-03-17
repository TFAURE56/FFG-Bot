package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var registeredCommands []*discordgo.ApplicationCommand

// Structure pour stocker dynamiquement les commandes
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

// Fonction pour enregistrer une commande
func RegisterCommand(s *discordgo.Session, guildID string, cmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("Impossible de créer la commande %s: %v", cmd.Name, err)
	} else {
		log.Printf("✅ Commande %s enregistrée avec succès", cmd.Name)
		registeredCommands = append(registeredCommands, cmd)
	}

	// Associer la commande à son handler
	commandHandlers[cmd.Name] = handler
}

func RegisterAllCommands(s *discordgo.Session) {
	if s.State.User == nil {
		log.Println("❌ Session utilisateur non disponible, attente de connexion...")
		return
	}

	// Récupérer toutes les guildes où le bot est présent
	guilds, err := s.UserGuilds(100, "", "", false)
	if err != nil {
		log.Printf("❌ Erreur lors de la récupération des guildes : %v", err)
		return
	}

	// Parcours de toutes les guildes et enregistrement des commandes
	for _, guild := range guilds {
		log.Printf("📌 Enregistrement des commandes pour la guilde %s (%s)", guild.Name, guild.ID)
		RegisterPingCommand(s, guild.ID)
		RegisterAddstockpileCommand(s, guild.ID)
		RegisterListstockpilesCommand(s, guild.ID)
		RegisterResetStockpileCommand(s, guild.ID)
		RegisterRemoveStockpileCommand(s, guild.ID)
		RegisterHelpCommand(s, guild.ID)
	}

	// Ajouter un seul handler global pour gérer toutes les commandes
	s.AddHandler(commandHandler)
}

// Handler unique pour toutes les commandes
func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, exists := commandHandlers[i.ApplicationCommandData().Name]; exists {
		handler(s, i)
	}
}

func UnregisterGlobalCommands(s *discordgo.Session) {
	// Supprimer les commandes globales
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Printf("Erreur lors de la récupération des commandes globales : %v", err)
		return
	}

	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Printf("Impossible de supprimer la commande globale %s: %v", cmd.Name, err)
		} else {
			log.Printf("Commande globale %s supprimée avec succès", cmd.Name)
		}
	}
}

func UnregisterGuildCommands(s *discordgo.Session) {
	// Récupérer toutes les guildes où le bot est présent
	guilds, err := s.UserGuilds(100, "", "", false) // Ajout de `false` pour correspondre à la nouvelle signature
	if err != nil {
		log.Printf("Erreur lors de la récupération des guildes : %v", err)
		return
	}

	for _, guild := range guilds {
		commands, err := s.ApplicationCommands(s.State.User.ID, guild.ID)
		if err != nil {
			log.Printf("Erreur lors de la récupération des commandes pour la guilde %s: %v", guild.ID, err)
			continue
		}

		for _, cmd := range commands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, cmd.ID)
			if err != nil {
				log.Printf("Impossible de supprimer la commande %s dans la guilde %s: %v", cmd.Name, guild.ID, err)
			} else {
				log.Printf("Commande %s supprimée avec succès dans la guilde %s", cmd.Name, guild.ID)
			}
		}
	}
}
