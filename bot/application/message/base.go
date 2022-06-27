package message

import (
	"emobot/bot/application"
	"emobot/bot/application/message/poll"
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
	pollCollection := poll.NewCommandCollection(c.database)

	commands = append(commands, pollCollection.GetAllCommands()...)

	return commands
}
