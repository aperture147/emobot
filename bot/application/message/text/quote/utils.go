package quote

import (
	"emobot/bot/application"
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
	createCommand := NewCreateQuoteCommand(c.collection)
	return []application.Command{createCommand}
}
