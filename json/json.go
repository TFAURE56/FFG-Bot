package json

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

type Stockpiles struct {
	Nom      string `json:"nom"`
	Hexa     string `json:"hexa"`
	Code     string `json:"code"`
	Cooldown int64  `json:"cooldown"` // Cooldown en timestamp UNIX (int64)
}

type ServerStockpiles struct {
	GuildID    string       `json:"guild_id"`
	Stockpiles []Stockpiles `json:"stockpiles"`
}

const filename = "stockpiles.json"

// Récupérer les stockpiles d'une guilde
func GetStockpiles(guildID string) ([]Stockpiles, error) {
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
			serverStockpiles[i].Stockpiles = append(serverStockpiles[i].Stockpiles, Stockpiles{
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
			Stockpiles: []Stockpiles{
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
func SaveStockpiles(guildID string, stockpiles []Stockpiles) error {
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

// Structure pour stocker le salon des cooldowns par serveur
type CooldownSettings struct {
	sync.Mutex
	GuildChannels map[string]string `json:"guild_channels"`
}

var cooldownData = CooldownSettings{GuildChannels: make(map[string]string)}

const settingsFile = "settings.json"

// Charge le fichier JSON au démarrage
func LoadSettings() error {
	file, err := os.Open(settingsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&cooldownData)
}

// Sauvegarde un salon pour les cooldowns
func SaveCooldownChannel(guildID, channelID string) error {
	cooldownData.Lock()
	defer cooldownData.Unlock()

	cooldownData.GuildChannels[guildID] = channelID

	file, err := os.Create(settingsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cooldownData)
}

// Récupère le salon des cooldowns d'un serveur
func GetCooldownChannel(guildID string) (string, bool) {
	cooldownData.Lock()
	defer cooldownData.Unlock()

	channel, exists := cooldownData.GuildChannels[guildID]
	return channel, exists
}
