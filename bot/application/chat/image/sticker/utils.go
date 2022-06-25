package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

const DatabaseCollectionName = "sticker"

type CommandCollection struct {
	collection *mongo.Collection
}

func NewCommandCollection(database *mongo.Database) CommandCollection {
	return CommandCollection{collection: database.Collection(DatabaseCollectionName)}
}

func (c CommandCollection) GetAllCommands() []application.Command {
	getCommand := NewGetStickerSlashCommand(c.collection)
	createCommand := NewCreateStickerSlashCommand(c.collection)
	deleteStickerCommand := NewDeleteStickerSlashCommand(c.collection)
	return []application.Command{getCommand, createCommand, deleteStickerCommand}
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
