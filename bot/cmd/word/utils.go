package word

import (
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetQuoteAutocompleteChoice(collection *mongo.Collection, findAttr string) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	stickers, err := db.GetQuoteAutocompleteList(collection, findAttr)
	if err != nil {
		return nil, err
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, sticker := range stickers {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  sticker.Title,
			Value: sticker.ObjectId,
		})
	}

	return choices, nil
}
