package routines

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Start(dg *discordgo.Session, db *sql.DB) {

	log.Println("Démarrage des routines...")

	var wg sync.WaitGroup

	// Routine pour la gestion des expériences
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			StartCooldown(dg, db)
			time.Sleep(15 * time.Second) // Ajustez la fréquence selon vos besoins
		}
	}()

	wg.Wait()
	log.Println("Routines terminées.")
}
