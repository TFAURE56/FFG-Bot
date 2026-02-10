package global

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Récupère la liste des IDs d'ordres depuis la BDD
func GetOrderIDsFromDB() []*discordgo.ApplicationCommandOptionChoice {
	db, err := ConnectToDatabase()
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM orders")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var choices []*discordgo.ApplicationCommandOptionChoice
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  "Commande n°" + strconv.FormatInt(id, 10),
				Value: id,
			})
		}
	}
	return choices
}

// Récupère la liste des éléments pour un ordre donné
func GetElementsForOrder(orderID int64) []*discordgo.ApplicationCommandOptionChoice {
	db, err := ConnectToDatabase()
	if err != nil {
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT ressource FROM order_elements WHERE order_id = ?", orderID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var choices []*discordgo.ApplicationCommandOptionChoice
	for rows.Next() {
		var element string
		if err := rows.Scan(&element); err == nil {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  element,
				Value: element,
			})
		}
	}
	return choices
}
