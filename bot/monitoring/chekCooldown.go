package monitoring

import (
	"FFG-Bot/json"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Délai entre chaque vérification (ex : toutes les 10 minutes)
const checkInterval = 1 * time.Hour

// Seuil du cooldown pour notification (ex: 1 heure avant expiration)
const cooldownThreshold = 12 * time.Hour

func StartCooldownMonitor(s *discordgo.Session) {
	go func() {
		for {
			checkStockpileCooldowns(s)
			time.Sleep(checkInterval)
		}
	}()
}

func checkStockpileCooldowns(s *discordgo.Session) {
	guilds := s.State.Guilds
	now := time.Now().Unix()

	for _, guild := range guilds {
		stockpiles, err := json.GetStockpiles(guild.ID)
		if err != nil {
			log.Printf("❌ Erreur lors de la récupération des stockpiles pour %s: %v", guild.ID, err)
			continue
		}

		for _, sp := range stockpiles {
			cooldownEnd := sp.Cooldown
			//			if err != nil {
			//				log.Printf("❌ Erreur conversion cooldown pour %s: %v", sp.Nom, err)
			//				continue
			//			}

			timeRemaining := cooldownEnd - now
			if timeRemaining > 0 && timeRemaining < int64(cooldownThreshold.Seconds()) {
				// 🚨 Cooldown bientôt terminé ! Envoyer une notification.
				alertCooldown(s, guild.ID, sp, timeRemaining)
			}
		}
	}
}

func alertCooldown(s *discordgo.Session, guildID string, sp json.Stockpiles, timeRemaining int64) {
	channelID, exists := json.GetCooldownChannel(guildID)
	if !exists {
		log.Printf("⚠️ Aucun salon défini pour les alertes cooldown sur le serveur %s\n", guildID)
		return
	}

	hours := timeRemaining / 3600
	minutes := (timeRemaining % 3600) / 60

	message := fmt.Sprintf("# ⚠️ **Alerte Cooldown** ⚠️\n\nLe stockpile **%s** situé à **%s** sera bientôt perdu. \nTemps restant : %d heures et %d minutes", sp.Nom, sp.Hexa, hours, minutes)

	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("❌ Erreur lors de l'envoi de l'alerte cooldown pour %s: %v\n", sp.Nom, err)
	}
}
