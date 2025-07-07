package routines

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var cooldownRoutineStarted = false

// Start initializes and starts the cooldown management routine.
func StartCooldown(dg *discordgo.Session, db *sql.DB) {
	if cooldownRoutineStarted {
		return
	}
	cooldownRoutineStarted = true

	log.Println("Démarrage de la routine de gestion des cooldowns...")

	go func() {
		for {
			getCooldowns(dg, db)

			time.Sleep(15 * time.Second)
		}
	}()

	log.Println("Routine de gestion des cooldowns démarrée.")
}

// function pour récupérer les cooldowns
func getCooldowns(dg *discordgo.Session, db *sql.DB) {

	rows, err := db.Query("SELECT name, hexa, cooldown, alerted FROM stockpiles WHERE alerted = 0 AND cooldown > 0")
	if err != nil {
		log.Printf("Erreur lors de la récupération des stockpiles: %v", err)
		return
	}
	defer rows.Close()

	var name, hexa string
	var cooldownTime time.Time
	var alerted int
	for rows.Next() {
		err := rows.Scan(&name, &hexa, &cooldownTime, &alerted)
		if err != nil {
			log.Printf("Erreur lors de la lecture des données: %v", err)
			continue
		}
		cooldown := cooldownTime.Unix()
		remaining := cooldown - time.Now().Unix()

		// Si le cooldown est inférieur à 2h (7200 secondes) et supérieur à 0
		if remaining > 0 && remaining < 2*3600 {
			sendCooldownAlertTimestamp(name, hexa, cooldown, dg)
			// Mettre à jour la base de données pour marquer comme alerté
			_, err := db.Exec("UPDATE stockpiles SET alerted = 1 WHERE name = ? AND hexa = ?", name, hexa)
			if err != nil {
				log.Printf("Erreur lors de la mise à jour du stockpile %s (hexa: %s): %v", name, hexa, err)
			} else {
				log.Printf("Stockpile %s (hexa: %s) marqué comme alerté.", name, hexa)
			}
		}
	}
}

// Envoie un message d'alerte avec un timestamp Discord si cooldown < 2h
func sendCooldownAlertTimestamp(name, hexa string, cooldownUnix int64, dg *discordgo.Session) {
	const channelID = "1307259612755525696"

	message := fmt.Sprintf(
		"⏰ Le stockpile **%s** (hexa: %s) expire <t:%d:R> (<t:%d:F>) !",
		name, hexa, cooldownUnix, cooldownUnix,
	)

	// Ajout d'un bouton "Reset"
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Reset",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("reset_stockpile_%s_%s", name, hexa),
				},
			},
		},
	}

	if dg != nil {
		_, err := dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content:    message,
			Components: components,
		})
		if err != nil {
			log.Printf("Erreur lors de l'envoi de l'alerte cooldown: %v", err)
			log.Printf("⚠️ Vérifiez que le bot a accès au salon %s et la permission d'envoyer des messages.", channelID)
		} else {
			log.Printf("Alerte cooldown envoyée pour le stockpile %s (hexa: %s) dans le salon %s", name, hexa, channelID)
		}
	} else {
		log.Printf("Session Discord non initialisée, impossible d'envoyer l'alerte.")
	}
}
