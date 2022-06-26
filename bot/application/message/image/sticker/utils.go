package sticker

import (
	"emobot/bot/application"
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
	createCommand := NewCreateStickerCommand(c.collection)
	return []application.Command{createCommand}
}
