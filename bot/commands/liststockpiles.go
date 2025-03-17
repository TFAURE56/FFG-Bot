package commands

import (
	"FFG-Bot/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RegisterListstockpilesCommand(s *discordgo.Session, guildID string) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "liststockpiles",
		Description: "Affiche la liste des stockpiles enregistrés pour ce serveur",
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
	if err != nil {
		log.Printf("❌ Impossible de créer la commande %s: %v", cmd.Name, err)
	} else {
		log.Printf("✅ Commande %s enregistrée avec succès", cmd.Name)
	}

	s.AddHandler(listStockpilesHandler)
}

func listStockpilesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "liststockpiles" {
		return
	}

	// Récupération de l'utilisateur et du serveur
	user := i.Member.User.Username
	guildID := i.GuildID

	// Changer le fuseau horaire
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Printf("�� Erreur lors du chargement du fuseau horaire : %v", err)
		return
	}
	time.Local = loc

	timestamp := time.Now().Format("02-01-2006 15:04:05")

	// Enregistrement dans le fichier de log
	logEntry := fmt.Sprintf("[%s] %s a exécuté /liststockpiles sur le serveur %s\n", timestamp, user, guildID)
	appendLog("logs/liststockpiles.log", logEntry)

	stockpiles, err := json.GetStockpiles(i.GuildID)
	if err != nil || len(stockpiles) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Aucun stockpile enregistré pour ce serveur.",
			},
		})
		return
	}

	var response string
	for _, sp := range stockpiles {
		// Calcul du temps restant
		cooldownStr := getTimeRemaining(sp.Cooldown)

		response += fmt.Sprintf("# 📦 **%s**\n🗺️ Hexagone: %s\n🔑 Code: ||`%s`||\n⏳ Temps restant: %s\n\n",
			sp.Nom, sp.Hexa, sp.Code, cooldownStr)
	}

	response += ("\n\n**Pour réinitialiser un stockpile, utilisez la commande `/resetstockpile <nom>`**")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Retourne le temps restant avant expiration du cooldown
func getTimeRemaining(expiration int64) string {
	now := time.Now().Unix()
	remaining := expiration - now

	if remaining <= 0 {
		return "✅ Disponible !"
	}

	hours := remaining / 3600
	minutes := (remaining % 3600) / 60
	seconds := remaining % 60

	return fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)
}

// Fonction pour ajouter une entrée au fichier de log
func appendLog(filename, entry string) {
	// Création du dossier logs s'il n'existe pas
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Ouverture du fichier en mode ajout
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("❌ Erreur d'écriture dans le fichier de log : %v", err)
		return
	}
	defer f.Close()

	// Écriture de l'entrée dans le fichier
	if _, err := f.WriteString(entry); err != nil {
		log.Printf("❌ Impossible d'écrire dans le fichier de log : %v", err)
	}
}
