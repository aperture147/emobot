package quote

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseCollectionName = "quote"

type CommandCollection struct {
	collection *mongo.Collection
}

func NewCommandCollection(database *mongo.Database) application.CommandCollection {
	return CommandCollection{collection: database.Collection(databaseCollectionName)}
}

func (c CommandCollection) GetAllCommands() []application.Command {
	getCommand := NewGetQuoteCommand(c.collection)
	deleteCommand := NewDeleteGetQuoteCommand(c.collection)
	return []application.Command{getCommand, deleteCommand}
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
