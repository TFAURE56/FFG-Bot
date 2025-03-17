package json

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Stockpile struct {
	Nom      string `json:"nom"`
	Hexa     string `json:"hexa"`
	Code     string `json:"code"`
	Cooldown int64  `json:"cooldown"` // Cooldown en timestamp UNIX (int64)
}

type ServerStockpiles struct {
	GuildID    string      `json:"guild_id"`
	Stockpiles []Stockpile `json:"stockpiles"`
}

const filename = "stockpiles.json"

// Récupérer les stockpiles d'une guilde
func GetStockpiles(guildID string) ([]Stockpile, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, errors.New("fichier stockpiles.json introuvable")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var serverStockpiles []ServerStockpiles
	err = json.Unmarshal(data, &serverStockpiles)
	if err != nil {
		return nil, err
	}

	for _, server := range serverStockpiles {
		if server.GuildID == guildID {
			return server.Stockpiles, nil
		}
	}

	return nil, nil
}

// Ajouter un stockpile
func AddStockpile(guildID, nom, hexa, code string, cooldown int64) error {
	var serverStockpiles []ServerStockpiles

	// Lire les données actuelles
	if _, err := os.Stat(filename); err == nil {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		json.Unmarshal(data, &serverStockpiles)
	}

	// Chercher si la guilde existe déjà
	var found bool
	for i, server := range serverStockpiles {
		if server.GuildID == guildID {
			serverStockpiles[i].Stockpiles = append(serverStockpiles[i].Stockpiles, Stockpile{
				Nom:      nom,
				Hexa:     hexa,
				Code:     code,
				Cooldown: cooldown,
			})
			found = true
			break
		}
	}

	// Si la guilde n'existe pas encore
	if !found {
		serverStockpiles = append(serverStockpiles, ServerStockpiles{
			GuildID: guildID,
			Stockpiles: []Stockpile{
				{
					Nom:      nom,
					Hexa:     hexa,
					Code:     code,
					Cooldown: cooldown,
				},
			},
		})
	}

	// Sauvegarde dans le fichier
	jsonData, err := json.MarshalIndent(serverStockpiles, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonData, 0644)
}

// Sauvegarde les stockpiles mis à jour
func SaveStockpiles(guildID string, stockpiles []Stockpile) error {
	// Charger les données existantes
	var allStockpiles []ServerStockpiles
	data, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(data, &allStockpiles)
	}

	// Mettre à jour les stockpiles du serveur
	updated := false
	for i, s := range allStockpiles {
		if s.GuildID == guildID {
			allStockpiles[i].Stockpiles = stockpiles
			updated = true
			break
		}
	}

	// Si le serveur n'existe pas encore dans le JSON
	if !updated {
		allStockpiles = append(allStockpiles, ServerStockpiles{
			GuildID:    guildID,
			Stockpiles: stockpiles,
		})
	}

	// Sauvegarde dans le fichier
	jsonData, err := json.MarshalIndent(allStockpiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}
