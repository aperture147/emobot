package chat

import (
	"emobot/bot/application"
	"emobot/bot/application/chat/image"
	"emobot/bot/application/chat/server"
	"emobot/bot/application/chat/text"
	"emobot/bot/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandCollection struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewCommandCollection(guildId string, client *mongo.Client) CommandCollection {
	return CommandCollection{
		client:   client,
		database: db.GetGuildDatabase(guildId, client),
	}
}

func (c CommandCollection) GetAllCommands() (commands []application.Command) {
	imageCollection := image.NewCommandCollection(c.database)
	textCollection := text.NewCommandCollection(c.database)
	serverCollection := server.NewCommandCollection(c.client)

	commands = append(commands, imageCollection.GetAllCommands()...)
	commands = append(commands, textCollection.GetAllCommands()...)
	commands = append(commands, serverCollection.GetAllCommands()...)

	return commands
}
