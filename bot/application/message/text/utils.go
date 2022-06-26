package text

import (
	"emobot/bot/application"
	"emobot/bot/application/message/text/quote"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandCollection struct {
	database *mongo.Database
}

func NewCommandCollection(database *mongo.Database) application.CommandCollection {
	return commandCollection{database: database}
}

func (c commandCollection) GetAllCommands() (commands []application.Command) {
	quoteCollection := quote.NewCommandCollection(c.database)

	commands = append(commands, quoteCollection.GetAllCommands()...)
	return commands
}
