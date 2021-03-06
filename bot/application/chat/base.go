package chat

import (
	"emobot/bot/application"
	"emobot/bot/application/chat/image"
	"emobot/bot/application/chat/server"
	"emobot/bot/application/chat/text"
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
	serverCollection := server.NewCommandCollection(c.client)

	commands = append(commands, imageCollection.GetAllCommands()...)
	commands = append(commands, textCollection.GetAllCommands()...)
	commands = append(commands, serverCollection.GetAllCommands()...)

	return commands
}
