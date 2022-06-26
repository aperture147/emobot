package message

import (
	"emobot/bot/application"
	"emobot/bot/application/message/image"
	"emobot/bot/application/message/text"
	"emobot/bot/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandCollection struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewCommandCollection(guildId string, client *mongo.Client) application.CommandCollection {
	return commandCollection{
		client:   client,
		database: db.GetGuildDatabase(guildId, client),
	}
}

func (c commandCollection) GetAllCommands() (commands []application.Command) {
	imageCollection := image.NewCommandCollection(c.database)
	textCollection := text.NewCommandCollection(c.database)

	commands = append(commands, imageCollection.GetAllCommands()...)
	commands = append(commands, textCollection.GetAllCommands()...)

	return commands
}
