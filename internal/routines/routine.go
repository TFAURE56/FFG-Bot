package routines

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Start(ctx context.Context, dg *discordgo.Session, db *sql.DB) {

	log.Println("Démarrage des routines...")

	var wg sync.WaitGroup

	// Routine pour le cooldown des stockpiles
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Arrêt de la routine cooldown...")
				return
			case <-ticker.C:
				StartCooldown(dg, db)
			}
		}
	}()

	wg.Wait()
	log.Println("Routines terminées.")
}
