package monitoring

import (
	"FFG-Bot/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Délai entre chaque vérification (ex : toutes les 10 minutes)
const checkInterval = 10 * time.Minute

// Seuil du cooldown pour notification (ex: 1 heure avant expiration)
const cooldownThreshold = 1 * time.Hour

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
			cooldownEnd, err := strconv.ParseInt(sp.Cooldown, 10, 64)
			if err != nil {
				log.Printf("❌ Erreur conversion cooldown pour %s: %v", sp.Nom, err)
				continue
			}

			timeRemaining := cooldownEnd - now
			if timeRemaining > 0 && timeRemaining < int64(cooldownThreshold.Seconds()) {
				// 🚨 Cooldown bientôt terminé ! Envoyer une notification.
				notifyCooldown(s, guild.ID, sp, timeRemaining)
			}
		}
	}
}

func notifyCooldown(s *discordgo.Session, guildID string, sp json.Stockpiles, timeRemaining int64) {
	channelID := "ID_DU_CHANNEL" // Remplace par l'ID du channel où envoyer l'alerte
	message := fmt.Sprintf("⚠️ **Stockpile %s** dans **%s** sera bientôt prêt ! Temps restant : **%d minutes**.",
		sp.Nom, sp.City, timeRemaining/60)

	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("❌ Erreur envoi notification cooldown: %v", err)
	}
}
