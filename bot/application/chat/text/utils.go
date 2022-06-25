package text

import (
	"emobot/bot/application"
	"emobot/bot/application/chat/text/quote"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandCollection struct {
	database *mongo.Database
}

func NewCommandCollection(database *mongo.Database) CommandCollection {
	return CommandCollection{database: database}
}

func (c CommandCollection) GetAllCommands() (commands []application.Command) {
	stickerCollection := quote.NewCommandCollection(c.database)

	commands = append(commands, stickerCollection.GetAllCommands()...)
	return commands
}
