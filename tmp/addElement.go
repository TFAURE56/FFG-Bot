package commands

import (
	"FFG-Bot/internal/global"
	"log"
	"strconv"
	"time"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Register(&discordgo.ApplicationCommand{
		Name:        "addelement",
		Description: "Ajouter un élément possible a la commande",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "Element",
				Description: "Nom de l'element à ajouter",
				Required:    true,
				Autocomplete: false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "Prix",
				Description: "Prix de l'element à ajouter",
				Required:    true,
				Autocomplete: false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "Money",
				Description: "Type de monnaie (Souffre, )",
				Required:    true,
				Autocomplete: false,
			},
		},
	}, addElementHandler)
}

func addElementHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "addelement" {
		return
	}

	options := i.ApplicationCommandData().Options
	var element, money string
	var prix int64	