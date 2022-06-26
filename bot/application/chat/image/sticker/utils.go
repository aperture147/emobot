package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseCollectionName = "sticker"

type commandCollection struct {
	collection *mongo.Collection
}

func NewCommandCollection(database *mongo.Database) application.CommandCollection {
	return commandCollection{collection: database.Collection(databaseCollectionName)}
}

func (c commandCollection) GetAllCommands() []application.Command {
	getCommand := NewGetStickerCommand(c.collection)
	createCommand := NewCreateStickerCommand(c.collection)
	deleteCommand := NewDeleteStickerCommand(c.collection)
	return []application.Command{getCommand, createCommand, deleteCommand}
}

func GetStickerAutocompleteChoice(collection *mongo.Collection, findAttr string) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	stickers, err := db.GetStickerAutocompleteList(collection, findAttr)
	if err != nil {
		return nil, err
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, sticker := range stickers {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  sticker.Name,
			Value: sticker.ObjectId,
		})
	}

	return choices, nil
}
