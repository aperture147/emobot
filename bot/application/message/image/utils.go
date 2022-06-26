package image

import (
	"emobot/bot/application"
	"emobot/bot/application/message/image/sticker"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandCollection struct {
	database *mongo.Database
}

func NewCommandCollection(database *mongo.Database) application.CommandCollection {
	return commandCollection{database: database}
}

func (c commandCollection) GetAllCommands() (commands []application.Command) {
	stickerCollection := sticker.NewCommandCollection(c.database)

	commands = append(commands, stickerCollection.GetAllCommands()...)
	return commands
}
