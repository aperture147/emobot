package poll

import (
	"emobot/bot/application"
	"go.mongodb.org/mongo-driver/mongo"
)

const databaseCollectionName = "poll"

type commandCollection struct {
	collection *mongo.Collection
}

func NewCommandCollection(database *mongo.Database) application.CommandCollection {
	return commandCollection{collection: database.Collection(databaseCollectionName)}
}

func (c commandCollection) GetAllCommands() []application.Command {
	createCommand := NewCreatePollCommand(c.collection)
	return []application.Command{createCommand}
}
