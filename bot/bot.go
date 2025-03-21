package bot

import (
	"FFG-Bot/bot/commands"
	"FFG-Bot/bot/monitoring"

	//"FFG-Bot/config"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Start(token string) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erreur lors de la création du bot", err)
	}

	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot en ligne !")
		commands.UnregisterGlobalCommands(s)
		//commands.UnregisterGuildCommands(s)
		commands.RegisterAllCommands(s)
	})

	err = s.Open()
	if err != nil {
		log.Fatal("Impossible de se connecter à Discord", err)
	}

	// Définition du statut du bot
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Status: "online", // Statut : "online", "idle", "dnd" (ne pas déranger) ou "invisible"
		Activities: []*discordgo.Activity{
			{
				Type: discordgo.ActivityTypeGame,
				Name: "la version 0.2.0",
			},
		},
	})

	//Lancement du chek cooldown
	monitoring.StartCooldownMonitor(s)

	select {}
}
