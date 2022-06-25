package quote

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

const DatabaseCollectionName = "quote"

type CommandCollection struct {
	collection *mongo.Collection
}

func NewCommandCollection(database *mongo.Database) CommandCollection {
	return CommandCollection{collection: database.Collection(DatabaseCollectionName)}
}

func (c CommandCollection) GetAllCommands() []application.Command {
	getCommand := NewGetQuoteCommand(c.collection)
	return []application.Command{getCommand}
}

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
